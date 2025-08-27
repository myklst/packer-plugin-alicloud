package eds

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type StepCreateImage struct {
}

func (s *StepCreateImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	// TODO: Create image
	return multistep.ActionContinue
}

func (s *StepCreateImage) Cleanup(state multistep.StateBag) {}
