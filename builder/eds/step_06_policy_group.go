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

type StepPolicyGroup struct {
	RegionId      string
	PolicyGroupId string

	autoCreated bool
}

func (s *StepPolicyGroup) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("alieds20200930").(*alieds.Client)
	ui := state.Get("ui").(packersdk.Ui)

	var err error
	if s.PolicyGroupId != "" {
		ui.Say("Querying policy group...")

		var resp *alieds.DescribePolicyGroupsResponse
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to query policy group: %s", err2.Error())
				}
				return retryable
			},
			RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
		}.Run(ctx, func(ctx context.Context) error {
			resp, err = client.DescribePolicyGroups(&alieds.DescribePolicyGroupsRequest{
				RegionId:      common.NilOrString(s.RegionId),
				PolicyGroupId: common.NilOrStringSlice(s.PolicyGroupId),
			})
			return err
		})
		if err != nil {
			return multistep.ActionHalt
		}

		if len(resp.Body.DescribePolicyGroups) <= 0 {
			ui.Errorf("Policy group %s not found", s.PolicyGroupId)
			return multistep.ActionHalt
		}
	} else {
		ui.Say("Creating policy group...")

		var (
			policyGroupName = fmt.Sprintf("packer-policy-group-%s", common.RandomString(5))
			resp            *alieds.CreatePolicyGroupResponse
		)
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to create policy group: %s", err2.Error())
				}
				return retryable
			},
		}.Run(ctx, func(ctx context.Context) error {
			resp, err = client.CreatePolicyGroup(&alieds.CreatePolicyGroupRequest{
				RegionId:  common.NilOrString(s.RegionId),
				Name:      common.NilOrString(policyGroupName),
				Clipboard: common.NilOrString("readwrite"),

				AuthorizeSecurityPolicyRule: []*alieds.CreatePolicyGroupRequestAuthorizeSecurityPolicyRule{
					{
						Description: alitea.String("Allow income all traffic"),
						Priority:    alitea.String("1"),
						Type:        alitea.String("inflow"),
						IpProtocol:  alitea.String("ALL"),
						CidrIp:      alitea.String("0.0.0.0/0"),
						PortRange:   alitea.String("-1/-1"),
						Policy:      alitea.String("accept"),
					},
				},
			})
			return err
		})
		if err != nil {
			return multistep.ActionHalt
		}

		s.PolicyGroupId = *resp.Body.PolicyGroupId
		s.autoCreated = true

		state.Put("policy_group_id", s.PolicyGroupId)
	}

	return multistep.ActionContinue
}

func (s *StepPolicyGroup) Cleanup(state multistep.StateBag) {
	if s.autoCreated && s.PolicyGroupId != "" {
		client := state.Get("alieds20200930").(*alieds.Client)
		ui := state.Get("ui").(packersdk.Ui)

		ui.Say("Deleting policy group...")

		var (
			ctx = context.TODO()
			err error
		)
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to delete policy group: %s", err2.Error())
				}
				return retryable
			},
		}.Run(ctx, func(ctx context.Context) error {
			_, err = client.DeletePolicyGroups(&alieds.DeletePolicyGroupsRequest{
				RegionId:      common.NilOrString(s.RegionId),
				PolicyGroupId: common.NilOrStringSlice(s.PolicyGroupId),
			})
			return err
		})
		if err != nil {
			return
		}
	}
}
