package eds

import (
	"context"
	"fmt"
	"strings"
	"time"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"
	alitea "github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

type StepUserCommand struct {
	RegionId        string
	CommandType     string
	CommandContent  string
	ContentEncoding string
	EndUserId       string
	CommandRole     string
	Timeout         uint64
}

func (s *StepUserCommand) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if s.CommandContent == "" {
		return multistep.ActionContinue
	}

	client := state.Get("client20200930").(*alieds.Client)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Running command:")
	ui.Say("```")
	ui.Say(strings.TrimSpace(s.CommandContent))
	ui.Say("```")

	var (
		cloudComputerId = state.Get("instance_id").(string)
		req             = alieds.RunCommandRequest{
			RegionId:        common.NilOrString(s.RegionId),
			Type:            common.NilOrString(s.CommandType),
			CommandContent:  common.NilOrString(s.CommandContent),
			ContentEncoding: common.NilOrString(s.ContentEncoding),
			DesktopId:       common.NilOrStringSlice(cloudComputerId),
			EndUserId:       common.NilOrString(s.EndUserId),
			CommandRole:     common.NilOrString(s.CommandRole),
		}
		resp *alieds.RunCommandResponse
		err  error
	)
	if s.Timeout > 0 {
		req.Timeout = alitea.Int64(int64(s.Timeout))
	}

	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, err2 := common.IsRetryableError(err)
			if !retryable {
				ui.Errorf("Failed to run command: %s", err2)
			}
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		resp, err = client.RunCommand(&req)
		return err
	})
	if err != nil {
		return multistep.ActionHalt
	}

	invokeId := *resp.Body.InvokeId
	ui.Sayf("Command triggered, invoke_id: %s", invokeId)

	return s.waitUntilCommandFinished(ctx, state, client, invokeId)
}

func (s *StepUserCommand) Cleanup(multistep.StateBag) {}

func (s *StepUserCommand) waitUntilCommandFinished(ctx context.Context, state multistep.StateBag, client *alieds.Client, invokeId string) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)

	var (
		resp *alieds.DescribeInvocationsResponse
		err  error
	)
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, err2 := common.IsRetryableError(err)
			if !retryable {
				ui.Errorf("Failed to describe invocation result: %s", err2)
			}
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		resp, err = client.DescribeInvocations(&alieds.DescribeInvocationsRequest{
			RegionId: common.NilOrString(s.RegionId),
			InvokeId: common.NilOrString(invokeId),
		})
		if err == nil {
			if len(resp.Body.Invocations) == 0 {
				return fmt.Errorf("invocation (%s) not found", invokeId)
			}
			if resp.Body.Invocations[0].InvocationStatus == alitea.String("Finished") {
				return nil
			}
		}
		return err
	})
	if err != nil {
		return multistep.ActionHalt
	}

	ui.Sayf("Command (%s) run successfully!", invokeId)
	return multistep.ActionContinue
}
