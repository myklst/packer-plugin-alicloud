package common

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"
)

func RunCommand(ctx context.Context, state multistep.StateBag, req *alieds.RunCommandRequest) multistep.StepAction {
	edsClient := state.Get("alieds20200930").(*alieds.Client)

	log.Printf(`
====================================
[DEBUG] Running command:
%s
====================================`, strings.TrimSpace(*req.CommandContent))

	var (
		resp *alieds.RunCommandResponse
		err  error
	)
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, _ := IsRetryableError(err)
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		resp, err = edsClient.RunCommand(req)

		if err != nil {
			log.Printf("[ERROR] RunCommand failed: %v", err)
			return err
		}

		if resp == nil || resp.Body == nil || resp.Body.InvokeId == nil {
			return fmt.Errorf("invalid RunCommand response")
		}

		log.Printf("[DEBUG] Command triggered, invoke_id: %s...", *resp.Body.InvokeId)
		return err
	})
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return waitUntilCommandFinished(state, edsClient, req.RegionId, resp.Body.InvokeId)
}

func waitUntilCommandFinished(state multistep.StateBag, client *alieds.Client, regionId, invokeId *string) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)

	var (
		resp *alieds.DescribeInvocationsResponse
		err  error
	)
	for {
		retry := false
		resp, err = client.DescribeInvocations(&alieds.DescribeInvocationsRequest{
			RegionId: regionId,
			InvokeId: invokeId,
		})
		if err != nil {
			retryable, err2 := IsRetryableError(err)
			if !retryable {
				ui.Errorf("Failed to describe invocation result: %s", err2)
				state.Put("error", fmt.Errorf("Failed to describe invocation result: %s", err2))
				return multistep.ActionHalt
			}
			retry = true
		} else {
			if len(resp.Body.Invocations) == 0 {
				ui.Errorf("Invocation (%s) not found", *invokeId)
				state.Put("error", fmt.Errorf("Invocation (%s) not found", *invokeId))
				return multistep.ActionHalt
			}

			status := *resp.Body.Invocations[0].InvocationStatus
			retry = status == "Pending" || status == "Running"
			if !retry && status != "Success" {
				log.Printf("[ERROR] Invocation (%s) failed with status: %s", *invokeId, status)
				state.Put("error", fmt.Errorf("Command Execution failed. Status: %s, ErrorInfo: %s", status, *resp.Body.Invocations[0].InvokeDesktops[0].ErrorInfo))
				return multistep.ActionHalt
			}
		}

		if retry {
			time.Sleep(10 * time.Second)
			continue
		}

		log.Printf("[DEBUG] Command (%s) run successfully!", *invokeId)
		return multistep.ActionContinue
	}
}
