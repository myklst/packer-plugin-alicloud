package eds

import (
	"context"
	"fmt"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

func AccessHost(client *alieds.Client, regionId string) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		ui := state.Get("ui").(packersdk.Ui)
		instanceId := state.Get("instance_id").(string)

		var (
			resp *alieds.DescribeDesktopsResponse
			err  error
		)
		err = retry.Config{
			Tries: 8,
			ShouldRetry: func(err error) bool {
				retryable, err2 := common.IsRetryableError(err)
				if !retryable {
					ui.Errorf("Failed to describe cloud computer: %s", err2)
				}
				return retryable
			},
		}.Run(context.Background(), func(ctx context.Context) error {
			resp, err = client.DescribeDesktops(&alieds.DescribeDesktopsRequest{
				RegionId:  common.NilOrString(regionId),
				DesktopId: common.NilOrStringSlice(instanceId),
			})
			if err == nil {
				if len(resp.Body.Desktops) == 0 {
					return fmt.Errorf("cloud computer (%s) not found", instanceId)
				}
				ui.Sayf("Cloud computer (%s) found, IP: %s", instanceId, *resp.Body.Desktops[0].NetworkInterfaceIp)
			}
			return err
		})
		if err != nil {
			return "", err
		}

		return *resp.Body.Desktops[0].NetworkInterfaceIp, nil
	}
}
