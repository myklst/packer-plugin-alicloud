package eds

import (
	"context"
	"time"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"
	alitea "github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

type StepCreateImage struct {
	RegionId     string
	NewImageName string
	Description  string

	autoCreated bool
}

func (s *StepCreateImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("alieds20200930").(*alieds.Client)
	computerId := state.Get("instance_id").(string)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Sayf("Creating image (%s) from cloud computer (%s)...", s.NewImageName, computerId)

	var (
		resp *alieds.CreateImageResponse
		err  error
	)
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, err2 := common.IsRetryableError(err)
			if !retryable {
				ui.Errorf("Failed to create image: %s", err2)
			}
			return retryable
		},
	}.Run(ctx, func(ctx context.Context) error {
		resp, err = client.CreateImage(&alieds.CreateImageRequest{
			RegionId:          common.NilOrString(s.RegionId),
			DesktopId:         common.NilOrString(computerId),
			ImageName:         common.NilOrString(s.NewImageName),
			Description:       common.NilOrString(s.Description),
			AutoCleanUserdata: alitea.Bool(true),
			DiskType:          alitea.String("SYSTEM"),
		})
		return err
	})
	if err != nil {
		return multistep.ActionHalt
	}

	s.autoCreated = true
	imageId := *resp.Body.ImageId
	state.Put("new_image_id", imageId)

	// It will take too long to wait for the image to be available.
	// And Alicloud EDS Team said:
	//   只要打镜像那个操作在实例被删掉前成功触发就行，后面会异步来制作镜像。
	// So we comment it out for now.
	// return s.waitUntilImageAvailable(ctx, state, client, imageId)

	return multistep.ActionContinue
}

func (s *StepCreateImage) Cleanup(state multistep.StateBag) {
	if !s.autoCreated {
		return
	}

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if !cancelled && !halted {
		return
	}

	client := state.Get("alieds20200930").(*alieds.Client)
	imageId := state.Get("new_image_id").(string)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Deregistering the Image because of cancellation, or error...")

	var (
		ctx = context.TODO()
		err error
	)
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, err2 := common.IsRetryableError(err)
			if !retryable {
				ui.Errorf("Failed to deregister image: %s", err2)
			}
			return retryable
		},
	}.Run(ctx, func(ctx context.Context) error {
		_, err = client.DeleteImages(&alieds.DeleteImagesRequest{
			RegionId: common.NilOrString(s.RegionId),
			ImageId:  common.NilOrStringSlice(imageId),
		})
		return err
	})
}

func (s *StepCreateImage) waitUntilImageAvailable(client *alieds.Client, imageId string) multistep.StepAction {
	for {
		retry := false
		resp, err := client.DescribeImages(&alieds.DescribeImagesRequest{
			RegionId:  common.NilOrString(s.RegionId),
			ImageType: alitea.String("CUSTOM"),
			ImageId:   common.NilOrStringSlice(imageId),
		})
		if err != nil {
			retryable, _ := common.IsRetryableError(err)
			if !retryable {
				return multistep.ActionHalt
			}
			retry = true
		}

		if *resp.Body.Images[0].Status == "Creating" {
			retry = true
		}

		if retry {
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}

	return multistep.ActionContinue
}
