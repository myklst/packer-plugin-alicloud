package eds

import (
	"context"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"
	alitea "github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/packer-plugin-sdk/multistep"

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

	cloudComputerId := state.Get("instance_id").(string)

	req := alieds.RunCommandRequest{
		RegionId:        common.NilOrString(s.RegionId),
		Type:            common.NilOrString(s.CommandType),
		CommandContent:  common.NilOrString(s.CommandContent),
		ContentEncoding: common.NilOrString(s.ContentEncoding),
		DesktopId:       common.NilOrStringSlice(cloudComputerId),
		CommandRole:     common.NilOrString(s.CommandRole),
	}
	if s.Timeout > 0 {
		req.Timeout = alitea.Int64(int64(s.Timeout))
	}

	return common.RunCommand(ctx, state, &req)
}

func (s *StepUserCommand) Cleanup(multistep.StateBag) {}
