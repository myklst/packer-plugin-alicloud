package eds

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"
)

type StepSetGeneratedData struct {
	GeneratedData *packerbuilderdata.GeneratedData
}

func (s *StepSetGeneratedData) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	instanceId := state.Get("instance_id").(string)
	instanceIp := state.Get("instance_ip").(string)
	cloudComputerUser := state.Get("cloud_computer_user").(string)
	sshKeyPair := state.Get("ssh_key_pair").(*sshKeyPair)

	s.GeneratedData.Put("InstanceId", instanceId)
	s.GeneratedData.Put("SshHost", instanceIp)
	s.GeneratedData.Put("SshUsername", cloudComputerUser)
	s.GeneratedData.Put("SshPrivKey", sshKeyPair.PrivKey)

	return multistep.ActionContinue
}

func (s *StepSetGeneratedData) Cleanup(state multistep.StateBag) {
}

func getGeneratedDataList() []string {
	return []string{
		"InstanceId",
		"SshHost",
		"SshUsername",
		"SshPrivKey",
	}
}
