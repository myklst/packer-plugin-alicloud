package eds

import (
	"context"
	"fmt"
	"strings"
	"time"

	alieds "github.com/alibabacloud-go/eds-user-20210308/client"
	alitea "github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

type StepCloudComputerUser struct {
	Comm *communicator.Config
	User *EdsUser
}

func (s *StepCloudComputerUser) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("alieds20210308").(*alieds.Client)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Creating cloud computer users...")

	var (
		users    = make([]*alieds.CreateUsersRequestUsers, 0, 1)
		password = common.RandomString(8)
		resp     *alieds.CreateUsersResponse
		err      error
	)
	if s.User.Name != "" {
		users = append(users, &alieds.CreateUsersRequestUsers{
			EndUserId:    alitea.String(s.User.Name),
			Password:     alitea.String(password),
			Email:        common.NilOrString(s.User.Email),
			Phone:        common.NilOrString(s.User.Phone),
			OwnerType:    common.NilOrString(s.User.OwnerType),
			OrgId:        common.NilOrString(s.User.OrgId),
			Remark:       common.NilOrString(s.User.Remark),
			RealNickName: common.NilOrString(s.User.RealNickName),
		})
	}

	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, err2 := common.IsRetryableError(err)
			if !retryable {
				ui.Errorf("Failed to create cloud computer users: %s", err2)
			}
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		resp, err = client.CreateUsers(&alieds.CreateUsersRequest{
			IsLocalAdmin: common.NilOrBool(true),
			Users:        users,
		})
		if err == nil {
			if len(resp.Body.CreateResult.FailedUsers) > 0 {
				return fmt.Errorf("failed users: %v", resp.Body.CreateResult.FailedUsers)
			}
		}
		return err
	})
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Say("====================================")
	ui.Sayf("Alicloud EDS Username: %s", s.User.Name)
	ui.Sayf("Alicloud EDS Password: %s", password)
	ui.Say("====================================")

	osType := state.Get("os_type").(string)
	if !strings.EqualFold(osType, "linux") {
		if s.Comm.SSHUsername == "" || strings.EqualFold(s.Comm.SSHUsername, "root") {
			s.Comm.SSHUsername = s.User.Name
		}
	}
	state.Put("cloud_computer_user", s.User.Name)

	return multistep.ActionContinue
}

func (s *StepCloudComputerUser) Cleanup(state multistep.StateBag) {
	client := state.Get("alieds20210308").(*alieds.Client)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Deleting cloud computer users...")

	var (
		ctx = context.TODO()
		err error
	)
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, err2 := common.IsRetryableError(err)
			if !retryable {
				ui.Errorf("Failed to delete cloud computer users: %s", err2)
			}
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(ctx, func(ctx context.Context) error {
		_, err = client.RemoveUsers(&alieds.RemoveUsersRequest{
			Users: common.NilOrStringSlice(s.User.Name),
		})
		return err
	})
	if err != nil {
		ui.Errorf("Failed to delete cloud computer users: %s", err)
	}
}
