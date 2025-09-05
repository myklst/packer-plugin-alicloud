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

type StepCloudComputerTemplate struct {
	RegionId                 string
	ComputerTemplateId       string
	InstanceType             string
	RootDiskSizeGib          int
	RootDiskPerformanceLevel string
	UserDiskSizeGib          []int32
	UserDiskPerformanceLevel string
	Language                 string

	autoCreated bool
}

func (s *StepCloudComputerTemplate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("alieds20200930").(*alieds.Client)
	sourceImageId := state.Get("image_id").(string)
	ui := state.Get("ui").(packersdk.Ui)

	var err error
	if s.ComputerTemplateId != "" {
		ui.Say("Querying cloud computer template...")

		var (
			resp *alieds.DescribeBundlesResponse
		)
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, _ := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to describe cloud computer template: %s", err)
				}
				return retryable
			},
		}.Run(ctx, func(ctx context.Context) error {
			resp, err = client.DescribeBundles(&alieds.DescribeBundlesRequest{
				RegionId: common.NilOrString(s.RegionId),
				BundleId: common.NilOrStringSlice(s.ComputerTemplateId),
			})
			if err == nil {
				if len(resp.Body.Bundles) == 0 {
					return fmt.Errorf("cloud computer template (%s) not found", s.ComputerTemplateId)
				}
			}
			return err
		})
		if err != nil {
			return multistep.ActionHalt
		}
	} else {
		ui.Say("Creating cloud computer template...")

		var (
			computerTemplateName = fmt.Sprintf("packer-computer-template-%s", common.RandomString(5))
			resp                 *alieds.CreateBundleResponse
		)
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to create cloud computer template: %s", err2)
				}
				return retryable
			},
			RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
		}.Run(ctx, func(ctx context.Context) error {
			resp, err = client.CreateBundle(&alieds.CreateBundleRequest{
				RegionId:                 common.NilOrString(s.RegionId),
				BundleName:               common.NilOrString(computerTemplateName),
				DesktopType:              common.NilOrString(s.InstanceType),
				ImageId:                  common.NilOrString(sourceImageId),
				RootDiskSizeGib:          alitea.Int32(int32(s.RootDiskSizeGib)),
				RootDiskPerformanceLevel: common.NilOrString(s.RootDiskPerformanceLevel),
				UserDiskSizeGib:          alitea.Int32Slice(s.UserDiskSizeGib),
				UserDiskPerformanceLevel: common.NilOrString(s.UserDiskPerformanceLevel),
				Language:                 common.NilOrString(s.Language),
			})
			return err
		})
		if err != nil {
			return multistep.ActionHalt
		}

		s.ComputerTemplateId = *resp.Body.BundleId
		s.autoCreated = true
	}

	state.Put("computer_template_id", s.ComputerTemplateId)
	return multistep.ActionContinue
}

func (s *StepCloudComputerTemplate) Cleanup(state multistep.StateBag) {
	if s.autoCreated && s.ComputerTemplateId != "" {
		client := state.Get("alieds20200930").(*alieds.Client)
		ui := state.Get("ui").(packersdk.Ui)

		ui.Say("Deleting cloud computer template...")

		var (
			ctx = context.TODO()
			err error
		)
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to delete cloud computer template: %s", err2)
				}
				return retryable
			},
			RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
		}.Run(ctx, func(ctx context.Context) error {
			_, err = client.DeleteBundles(&alieds.DeleteBundlesRequest{
				RegionId: common.NilOrString(s.RegionId),
				BundleId: common.NilOrStringSlice(s.ComputerTemplateId),
			})
			return err
		})
		if err != nil {
			return
		}
	}
}
