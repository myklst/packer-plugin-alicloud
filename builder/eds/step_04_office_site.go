package eds

import (
	"context"
	"fmt"
	"time"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"
	alitea "github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

type StepOfficeSite struct {
	RegionId             string
	OfficeSiteId         string
	CidrBlock            string
	EnableInternetAccess bool
	InternetBandwidth    int32
	CenId                string
	CenOwnerId           int64
	CenVerifyCode        string

	autoCreated bool
}

func (s *StepOfficeSite) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("alieds20200930").(*alieds.Client)
	ui := state.Get("ui").(packersdk.Ui)

	var err error
	if s.OfficeSiteId != "" {
		ui.Say("Querying office network...")

		var resp *alieds.DescribeOfficeSitesResponse
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to query office network: %s", err2.Error())
				}
				return retryable
			},
			RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
		}.Run(ctx, func(ctx context.Context) error {
			resp, err = client.DescribeOfficeSites(&alieds.DescribeOfficeSitesRequest{
				RegionId:     common.NilOrString(s.RegionId),
				OfficeSiteId: common.NilOrStringSlice(s.OfficeSiteId),
			})
			return err
		})
		if err != nil {
			return multistep.ActionHalt
		}

		if len(resp.Body.OfficeSites) <= 0 {
			ui.Errorf("Office network %s is not found", s.OfficeSiteId)
			return multistep.ActionHalt
		}
	} else {
		ui.Say("Creating office network...")

		if s.CidrBlock == "" {
			s.CidrBlock = "192.168.0.0/24"
		}

		var (
			officeSiteName = fmt.Sprintf("packer-office-site-%s", common.RandomString(5))

			resp *alieds.CreateSimpleOfficeSiteResponse
		)
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to create office network: %s", err2.Error())
				}
				return retryable
			},
			RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
		}.Run(ctx, func(ctx context.Context) error {
			resp, err = client.CreateSimpleOfficeSite(&alieds.CreateSimpleOfficeSiteRequest{
				RegionId:             common.NilOrString(s.RegionId),
				OfficeSiteName:       common.NilOrString(officeSiteName),
				DesktopAccessType:    common.NilOrString("Internet"),
				CidrBlock:            common.NilOrString(s.CidrBlock),
				EnableInternetAccess: common.NilOrBool(s.EnableInternetAccess),
				Bandwidth:            alitea.Int32(s.InternetBandwidth),
				CenId:                common.NilOrString(s.CenId),
				CenOwnerId:           common.NilOrInt64(s.CenOwnerId),
				VerifyCode:           common.NilOrString(s.CenVerifyCode),
			})
			return err
		})
		if err != nil {
			return multistep.ActionHalt
		}

		s.OfficeSiteId = *resp.Body.OfficeSiteId
		s.autoCreated = true

		if err := s.waitUntil(ctx, "REGISTERED", client); err != nil {
			ui.Say("Waiting for office network to be registered...")
			return multistep.ActionHalt
		}

		state.Put("office_site_id", s.OfficeSiteId)
	}

	return multistep.ActionContinue
}

func (s *StepOfficeSite) Cleanup(state multistep.StateBag) {
	if s.autoCreated && s.OfficeSiteId != "" {
		client := state.Get("alieds20200930").(*alieds.Client)
		ui := state.Get("ui").(packersdk.Ui)

		ui.Say("Deleting office network...")

		var (
			ctx = context.TODO()
			err error
		)
		// Trigger the deletion of the office network
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to delete office network: %s", err2.Error())
				}
				return retryable
			},
			RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
		}.Run(ctx, func(ctx context.Context) error {
			_, err = client.DeleteOfficeSites(&alieds.DeleteOfficeSitesRequest{
				RegionId:     common.NilOrString(s.RegionId),
				OfficeSiteId: common.NilOrStringSlice(s.OfficeSiteId),
			})
			return err
		})
		if err != nil {
			return
		}

		// Wait for the office network to be deleted
		s.waitUntil(ctx, "DELETED", client)
	}
}

func (s *StepOfficeSite) waitUntil(ctx context.Context, targetStatus string, client *alieds.Client) error {
	err := retry.Config{
		StartTimeout: 10 * time.Minute,
		ShouldRetry: func(err error) bool {
			retryable, _ := common.IsRetryableError(err)
			return retryable
		},
		RetryDelay: func() time.Duration {
			return 5 * time.Second
		},
	}.Run(ctx, func(ctx context.Context) error {
		resp, err := client.DescribeOfficeSites(&alieds.DescribeOfficeSitesRequest{
			RegionId:     common.NilOrString(s.RegionId),
			OfficeSiteId: common.NilOrStringSlice(s.OfficeSiteId),
		})
		if err != nil {
			return err
		}
		if (targetStatus == "DELETED") && len(resp.Body.OfficeSites) != 0 {
			return fmt.Errorf("office network %s still exists", s.OfficeSiteId)
		}
		if len(resp.Body.OfficeSites) > 0 {
			current := *resp.Body.OfficeSites[0].Status
			if current == targetStatus {
				return nil
			}
			return fmt.Errorf("unexpected office network (%s) status: %s, expected: %s", s.OfficeSiteId, current, targetStatus)
		}
		return nil
	})
	return err
}
