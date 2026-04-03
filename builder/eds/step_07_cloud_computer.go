package eds

import (
	"context"
	"fmt"
	"time"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

type StepCloudComputer struct {
	RegionId                string
	ResourceGroupId         string
	ComputerPoolId          string
	ComputerTemplateId      string
	OfficeSiteId            string
	DesktopName             string
	DesktopNameSuffix       bool
	PolicyGroupId           string
	PromotionId             string
	Hostname                string
	DesktopMemberIp         string
	VolumeEncryptionEnabled bool
	VolumeEncryptionKey     string

	instanceId string
}

func (s *StepCloudComputer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if s.OfficeSiteId == "" {
		s.OfficeSiteId = state.Get("office_site_id").(string)
	}

	if s.ComputerTemplateId == "" {
		s.ComputerTemplateId = state.Get("computer_template_id").(string)
	}

	if s.PolicyGroupId == "" {
		s.PolicyGroupId = state.Get("policy_group_id").(string)
	}

	client := state.Get("alieds20200930").(*alieds.Client)
	ui := state.Get("ui").(packersdk.Ui)
	cloudComputerUser := state.Get("cloud_computer_user").(string)

	ui.Say("Creating cloud computer...")

	computerName := s.DesktopName
	if s.DesktopName == "" {
		computerName = fmt.Sprintf("packer-computer-%s", common.RandomString(5))
	}

	var (
		resp *alieds.CreateDesktopsResponse
		err  error
	)
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, err2 := common.IsRetryableError(err)
			if !retryable {
				ui.Errorf("Error creating cloud computer: %s", err2)
			}
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		resp, err = client.CreateDesktops(&alieds.CreateDesktopsRequest{
			ResourceGroupId:         common.NilOrString(s.ResourceGroupId),
			RegionId:                common.NilOrString(s.RegionId),
			GroupId:                 common.NilOrString(s.ComputerPoolId),
			BundleId:                common.NilOrString(s.ComputerTemplateId),
			OfficeSiteId:            common.NilOrString(s.OfficeSiteId),
			DesktopName:             common.NilOrString(computerName),
			DesktopNameSuffix:       common.NilOrBool(s.DesktopNameSuffix),
			PolicyGroupId:           common.NilOrString(s.PolicyGroupId),
			Hostname:                common.NilOrString(s.Hostname),
			EndUserId:               common.NilOrStringSlice(cloudComputerUser),
			DesktopMemberIp:         common.NilOrString(s.DesktopMemberIp),
			VolumeEncryptionEnabled: common.NilOrBool(s.VolumeEncryptionEnabled),
			VolumeEncryptionKey:     common.NilOrString(s.VolumeEncryptionKey),
			PromotionId:             common.NilOrString(s.PromotionId),
		})
		return err
	})
	if err != nil {
		return multistep.ActionHalt
	}

	computer := resp.Body.DesktopId
	s.instanceId = *computer[0]
	state.Put("instance_id", *computer[0])

	s.waitUntil(ctx, state, "Running", client)

	return multistep.ActionContinue
}

func (s *StepCloudComputer) Cleanup(state multistep.StateBag) {
	if s.instanceId == "" {
		return
	}

	if s.instanceId != "" {
		client := state.Get("alieds20200930").(*alieds.Client)
		ui := state.Get("ui").(packersdk.Ui)

		ui.Say("Terminating the cloud computer...")
		if _, err := client.DeleteDesktops(&alieds.DeleteDesktopsRequest{
			RegionId:  common.NilOrString(s.RegionId),
			DesktopId: common.NilOrStringSlice(s.instanceId),
		}); err != nil {
			ui.Errorf("Error terminating cloud computer, may still be around: %s", err)
			return
		}

		ctx := context.TODO()
		err := retry.Config{
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to get computer info: %s", err2)
				}
				return retryable
			},
			RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
		}.Run(ctx, func(ctx context.Context) error {
			_, err := client.DescribeDesktops(&alieds.DescribeDesktopsRequest{
				RegionId:  common.NilOrString(s.RegionId),
				DesktopId: common.NilOrStringSlice(s.instanceId),
			})
			return err
		})
		if err != nil {
			ui.Errorf("Error waiting desktop termination, may still be around: %s", err)
			return
		}

		s.waitUntil(ctx, state, "Deleted", client)
	}
}

func (s *StepCloudComputer) waitUntil(ctx context.Context, state multistep.StateBag, targetStatus string, client *alieds.Client) error {
	var (
		resp *alieds.DescribeDesktopsResponse
		err  error
	)
	err = retry.Config{
		StartTimeout: 10 * time.Minute,
		ShouldRetry: func(err error) bool {
			retryable, _ := common.IsRetryableError(err)
			return retryable
		},
		RetryDelay: func() time.Duration {
			return 5 * time.Second
		},
	}.Run(ctx, func(ctx context.Context) error {
		resp, err = client.DescribeDesktops(&alieds.DescribeDesktopsRequest{
			RegionId:  common.NilOrString(s.RegionId),
			DesktopId: common.NilOrStringSlice(s.instanceId),
		})
		if err != nil {
			return err
		}

		if targetStatus == "Deleted" {
			if len(resp.Body.Desktops) > 0 {
				return fmt.Errorf("cloud computer %s still exists", s.instanceId)
			}
			return nil
		}
		if len(resp.Body.Desktops) > 0 {
			current := *resp.Body.Desktops[0].DesktopStatus
			if current == targetStatus {
				return nil
			}
			return fmt.Errorf("unexpected cloud computer (%s) status: %s, expected: %s", s.instanceId, current, targetStatus)
		}
		return nil
	})
	if err == nil && len(resp.Body.Desktops) > 0 {
		state.Put("instance_ip", *resp.Body.Desktops[0].NetworkInterfaceIp)
		time.Sleep(60 * time.Second) // Wait for the instance to be fully ready
	}

	return err
}
