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
	sshUsername := ""
	sshPrivKey := []byte(nil)

	if rawConfig, ok := state.GetOk("config"); ok {
		if cfg, ok := rawConfig.(*Config); ok && cfg != nil {
			sshUsername = cfg.RunConfig.Comm.SSHUsername
			sshPrivKey = cfg.RunConfig.Comm.SSHPrivateKey
		}
	}

	if sshUsername == "" {
		if cloudComputerUser, ok := state.Get("cloud_computer_user").(string); ok {
			sshUsername = cloudComputerUser
		}
	}

	if len(sshPrivKey) == 0 {
		if rawKeyPair, ok := state.GetOk("ssh_key_pair"); ok {
			if keyPair, ok := rawKeyPair.(*sshKeyPair); ok && keyPair != nil {
				sshPrivKey = keyPair.PrivKey
			}
		}
	}

	s.GeneratedData.Put("InstanceId", instanceId)
	s.GeneratedData.Put("SshHost", instanceIp)
	s.GeneratedData.Put("SshUsername", sshUsername)
	s.GeneratedData.Put("SshPrivKey", sshPrivKey)

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
