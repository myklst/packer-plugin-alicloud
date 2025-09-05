package eds

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

func SshHost() func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		instanceId := state.Get("instance_id").(string)
		instanceIp := state.Get("instance_ip").(string)
		if instanceIp == "" {
			return "", fmt.Errorf("cloud computer (%s) not found", instanceId)
		}
		return instanceIp, nil
	}
}
