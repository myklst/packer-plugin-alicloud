package eds

import (
	"context"
	"fmt"
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
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

function Log([string]$Message) {
    Write-Output "[packer-ssh-setup] $Message"
}

# 1. Install OpenSSH Server
if ($null -eq (Get-Service -Name sshd -ErrorAction SilentlyContinue)) {
    Log 'Installing OpenSSH Server from GitHub...'

    $zipPath = Join-Path $env:TEMP 'OpenSSH-Win64.zip'

    $downloadArgs = @{
        Uri             = 'https://github.com/PowerShell/Win32-OpenSSH/releases/download/v9.5.0.0p1-Beta/OpenSSH-Win64.zip'
        OutFile         = $zipPath
        UseBasicParsing = $true
        TimeoutSec      = 360
    }
    Invoke-WebRequest @downloadArgs

    Expand-Archive -Path $zipPath -DestinationPath 'C:\Program Files\OpenSSH' -Force
    Remove-Item -Path $zipPath -Force -ErrorAction SilentlyContinue

    $installScript = 'C:\Program Files\OpenSSH\OpenSSH-Win64\install-sshd.ps1'

    if (-not (Test-Path $installScript)) {
        throw 'install-sshd.ps1 not found.'
    }

    & powershell.exe -NoProfile -ExecutionPolicy Bypass -File $installScript

    if ($LASTEXITCODE -ne 0) {
        throw "OpenSSH installation failed with exit code $LASTEXITCODE."
    }

    Log 'OpenSSH Server installed successfully.'
}

# 2. Start SSH service
Set-Service -Name sshd -StartupType Automatic
Start-Service -Name sshd
Log 'sshd is running.'

# 3. Configure Windows Firewall
if ($null -eq (Get-NetFirewallRule -Name 'sshd-Server-In-TCP' -ErrorAction SilentlyContinue)) {
    $firewallArgs = @{
        Name        = 'sshd-Server-In-TCP'
        DisplayName = 'Inbound rule for OpenSSH Server on TCP port 22'
        Action      = 'Allow'
        Direction   = 'Inbound'
        Enabled     = 'True'
        Profile     = 'Any'
        Protocol    = 'TCP'
        LocalPort   = 22
    }
    New-NetFirewallRule @firewallArgs | Out-Null
}

Log 'Firewall rule for SSH is present.'

# 4. Set PowerShell as default SSH shell
$openSSHKey = 'HKLM:\SOFTWARE\OpenSSH'

if (-not (Test-Path $openSSHKey)) {
    New-Item -Path $openSSHKey -Force | Out-Null
}

$defaultShellArgs = @{
    Path         = $openSSHKey
    Name         = 'DefaultShell'
    Value        = 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe'
    PropertyType = 'String'
    Force        = $true
}
New-ItemProperty @defaultShellArgs | Out-Null

Log 'Default SSH shell set to PowerShell.'

Log 'SSH service setup completed.'
exit 0`
	case "linux":
		runcmdType = "RunShellScript"
		cmd = `
# Allow SSH traffic through the firewall, when ufw is in use.
if command -v ufw >/dev/null 2>&1; then
    sudo ufw allow 22/tcp
fi

exit 0`
	default:
		state.Put("error", fmt.Errorf("unsupported OS type: %s", osType))
		return multistep.ActionHalt
	}

	return common.RunCommand(ctx, state, &alieds.RunCommandRequest{
		RegionId:        alitea.String(s.RegionId),
		Type:            alitea.String(runcmdType),
		CommandContent:  alitea.String(cmd),
		ContentEncoding: alitea.String("PlainText"),
		DesktopId:       common.NilOrStringSlice(computerId),
		Timeout:         alitea.Int64(1800),
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
$ErrorActionPreference = 'Continue'

if ($null -ne (Get-Service -Name sshd -ErrorAction SilentlyContinue)) {
    Set-Service -Name sshd -Status Stopped -StartupType Disabled
}
if ($null -ne (Get-NetFirewallRule -Name 'sshd-Server-In-TCP' -ErrorAction SilentlyContinue)) {
    Remove-NetFirewallRule -Name 'sshd-Server-In-TCP'
}

exit 0`
	case "linux":
		runcmdType = "RunShellScript"
		cmd = `
# Deny SSH traffic through the firewall, when ufw is in use.
if command -v ufw >/dev/null 2>&1; then
    sudo ufw deny 22/tcp
fi

exit 0`
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
		Timeout:         alitea.Int64(900),
	})
}
