package eds

import (
	"context"
	"fmt"
	"time"

	"github.com/myklst/packer-plugin-alicloud/builder/common"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"
)

type StepPreValidate struct {
	RegionId     string
	NewImageName string
}

func (s *StepPreValidate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)

	if config.AlicloudSkipValidation {
		ui.Say("Skip validation flag found, skipping prevalidating.")
		return multistep.ActionContinue
	}

	if err := s.validateRegions(state); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if err := s.validateImageName(state); err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepPreValidate) validateRegions(state multistep.StateBag) error {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Prevalidating regions...")

	var errs *packersdk.MultiError
	regions := []string{config.AlicloudRegion, s.RegionId}
	for _, region := range regions {
		if err := config.ValidateRegion(region); err != nil {
			errs = packersdk.MultiErrorAppend(errs, err)
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (s *StepPreValidate) validateImageName(state multistep.StateBag) error {
	if s.NewImageName == "" {
		return nil
	}

	client := state.Get("alieds20200930").(*alieds.Client)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Prevalidating image name...")

	ctx := context.TODO()
	var (
		resp *alieds.DescribeImagesResponse
		err  error
	)
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, err2 := common.IsRetryableError(err)
			if !retryable {
				ui.Errorf("Failed to query alicloud image on pre-validate step: %s", err2)
			}
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		resp, err = client.DescribeImages(&alieds.DescribeImagesRequest{
			RegionId:  common.NilOrString(s.RegionId),
			ImageName: common.NilOrString(s.NewImageName),
		})
		return err
	})
	if err != nil {
		return err
	}

	images := resp.Body.Images
	if len(images) > 0 {
		return fmt.Errorf("image name (%s) is used by an existing alicloud image: %s",
			*images[0].Name, *images[0].ImageId)
	}

	return nil
}

func (s *StepPreValidate) Cleanup(multistep.StateBag) {}
