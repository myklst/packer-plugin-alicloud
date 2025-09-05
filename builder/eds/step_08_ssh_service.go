package eds

import (
	"context"
	"strings"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"
	alitea "github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

type StepSshService struct {
	RegionId  string
	EndUserId string
}

func (s *StepSshService) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)
	osType := state.Get("os_type").(string)
	computerId := state.Get("instance_id").(string)

	ui.Say("Setting up SSH service...")

	var (
		cmd        string
		runcmdType string
	)
	switch strings.ToLower(osType) {
	case "windows":
		runcmdType = "RunPowerShellScript"
		cmd = `
Get-WindowsCapability -Name OpenSSH.Server* -Online | Add-WindowsCapability -Online
Set-Service -Name sshd -Status Running

# Allow SSH traffic through firewall
$firewallParams = @{
    Name        = 'sshd-Server-In-TCP'
    DisplayName = 'Inbound rule for OpenSSH Server (sshd) on TCP port 22'
    Action      = 'Allow'
    Direction   = 'Inbound'
    Enabled     = 'True'  # This is not a boolean but an enum
    Profile     = 'Any'
    Protocol    = 'TCP'
    LocalPort   = 22
}
New-NetFirewallRule @firewallParams

# Change default shell to PowerShell for Ansible SSH comminucation.
$shellParams = @{
    Path         = 'HKLM:\SOFTWARE\OpenSSH'
    Name         = 'DefaultShell'
    Value        = 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe'
    PropertyType = 'String'
    Force        = $true
}
New-ItemProperty @shellParams`
	case "linux":
		runcmdType = "RunShellScript"
		cmd = `
# Allow SSH traffic through firewall
sudo ufw allow 22/tcp`
	default:
		ui.Errorf("Unsupported OS type: %s", osType)
		return multistep.ActionHalt
	}

	return common.RunCommand(ctx, state, &alieds.RunCommandRequest{
		RegionId:        alitea.String(s.RegionId),
		Type:            alitea.String(runcmdType),
		CommandContent:  alitea.String(cmd),
		ContentEncoding: alitea.String("PlainText"),
		DesktopId:       common.NilOrStringSlice(computerId),
		EndUserId:       common.NilOrString(s.EndUserId),
		Timeout:         alitea.Int64(900),
	})
}

func (s *StepSshService) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packersdk.Ui)
	osType := state.Get("os_type").(string)
	computerId := state.Get("instance_id").(string)

	ui.Say("Stopping SSH service...")

	var (
		cmd        string
		runcmdType string
	)
	switch strings.ToLower(osType) {
	case "windows":
		runcmdType = "RunPowerShellScript"
		cmd = `
Set-Service -Name sshd -Status Stopped -StartupType Disabled

# Allow SSH traffic through firewall
$firewallParams = @{
    Name = 'sshd-Server-In-TCP'
}
Remove-NetFirewallRule @firewallParams`
	case "linux":
		runcmdType = "RunShellScript"
		cmd = `
# Allow SSH traffic through firewall
sudo ufw allow 22/tcp`
	default:
		ui.Errorf("Unsupported OS type: %s", osType)
		return
	}

	common.RunCommand(context.TODO(), state, &alieds.RunCommandRequest{
		RegionId:        alitea.String(s.RegionId),
		Type:            alitea.String(runcmdType),
		CommandContent:  alitea.String(cmd),
		ContentEncoding: alitea.String("PlainText"),
		DesktopId:       common.NilOrStringSlice(computerId),
		EndUserId:       common.NilOrString(s.EndUserId),
		Timeout:         alitea.Int64(900),
	})
}
