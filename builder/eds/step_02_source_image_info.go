package eds

import (
	"context"
	"time"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"
	"github.com/myklst/packer-plugin-alicloud/builder/common"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"
)

type StepSourceImageInfo struct {
	RegionId          string
	SourceImageFilter *EdsImageFilter
}

func (s *StepSourceImageInfo) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("alieds20200930").(*alieds.Client)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Querying source image info...")

	var (
		resp *alieds.DescribeImagesResponse
		err  error
	)
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, err2 := common.IsRetryableError(err)
			if !retryable {
				ui.Errorf("Failed to query alicloud image on source image info step: %s", err2)
			}
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		resp, err = client.DescribeImages(&alieds.DescribeImagesRequest{
			RegionId:            common.NilOrString(s.RegionId),
			ImageId:             common.NilOrStringSlice(s.SourceImageFilter.ImageId),
			ImageName:           common.NilOrString(s.SourceImageFilter.ImageName),
			FotaVersion:         common.NilOrString(s.SourceImageFilter.ImageVersion),
			ImageType:           common.NilOrString(s.SourceImageFilter.ImageType),
			ImageStatus:         common.NilOrString(s.SourceImageFilter.ImageStatus),
			ProtocolType:        common.NilOrString(s.SourceImageFilter.ProtocolType),
			DesktopInstanceType: common.NilOrString(s.SourceImageFilter.InstanceType),
			SessionType:         common.NilOrString(s.SourceImageFilter.SessionType),
			GpuCategory:         common.NilOrBool(s.SourceImageFilter.GpuCategory),
			GpuDriverVersion:    common.NilOrString(s.SourceImageFilter.GpuDriverVersion),
			OsType:              common.NilOrString(s.SourceImageFilter.OsType),
			LanguageType:        common.NilOrString(s.SourceImageFilter.OsLanguageType),
		})
		return err
	})
	if err != nil {
		return multistep.ActionHalt
	}

	images := resp.Body.Images
	if len(images) <= 0 {
		ui.Errorf("Image not found.")
		return multistep.ActionHalt
	}
	if len(resp.Body.Images) > 1 {
		ui.Errorf("Multiple images found: %d", len(resp.Body.Images))
		return multistep.ActionHalt
	}

	state.Put("image_id", *resp.Body.Images[0].ImageId)
	state.Put("os_type", *resp.Body.Images[0].OsType)

	return multistep.ActionContinue
}

func (s *StepSourceImageInfo) Cleanup(multistep.StateBag) {}
