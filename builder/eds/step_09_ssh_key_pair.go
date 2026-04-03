package eds

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh"

	alieds "github.com/alibabacloud-go/ecd-20200930/v5/client"
	alitea "github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

type sshKeyPair struct {
	PubKey  []byte
	PrivKey []byte
}

type StepSshKeyPair struct {
	Comm *communicator.Config

	RegionId  string
	EndUserId string

	sshKeyPair   *sshKeyPair
	debugKeyPath string
}

func (s *StepSshKeyPair) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)

	if s.Comm.SSHPrivateKeyFile != "" {
		ui.Say("Using existing SSH private key")
		privateKeyBytes, err := s.Comm.ReadSSHPrivateKeyFile()
		if err != nil {
			ui.Errorf("Failed to read existing SSH private key file: %s", err)
			return multistep.ActionHalt
		}

		s.Comm.SSHPrivateKey = privateKeyBytes

		return multistep.ActionContinue
	}

	if s.Comm.SSHAgentAuth && s.Comm.SSHKeyPairName == "" {
		ui.Say("Using SSH Agent with key pair in Source Image")
		return multistep.ActionContinue
	}

	if s.Comm.SSHAgentAuth && s.Comm.SSHKeyPairName != "" {
		ui.Say(fmt.Sprintf("Using SSH Agent for existing key pair %s", s.Comm.SSHKeyPairName))
		return multistep.ActionContinue
	}

	if s.Comm.SSHTemporaryKeyPairName == "" {
		ui.Say("Not using temporary keypair")
		s.Comm.SSHKeyPairName = ""
		return multistep.ActionContinue
	}

	ui.Say("Creating ephemeral SSH key pair...")

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		ui.Errorf("Failed to generate Ed25519 key: %v", err)
		return multistep.ActionHalt
	}

	// Marshal the private key to OpenSSH format (PEM encoded)
	pemBlock, err := ssh.MarshalPrivateKey(privateKey, "OPENSSH PRIVATE KEY")
	if err != nil {
		ui.Errorf("Failed to marshal SSH private key: %v", err)
		return multistep.ActionHalt
	}
	privatePem := pem.EncodeToMemory(pemBlock)

	// Create the public key in OpenSSH authorized_keys format
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		ui.Errorf("Failed to create SSH public key: %v", err)
		return multistep.ActionHalt
	}
	publicAuthorizedKey := ssh.MarshalAuthorizedKey(sshPublicKey)

	s.Comm.SSHKeyPairName = s.Comm.SSHTemporaryKeyPairName
	s.Comm.SSHPrivateKey = privatePem

	s.sshKeyPair = &sshKeyPair{
		PubKey:  publicAuthorizedKey,
		PrivKey: privatePem,
	}
	state.Put("ssh_key_pair", s.sshKeyPair)

	// Output the private key to the working directory.
	instanceIp := state.Get("instance_ip").(string)
	s.debugKeyPath = fmt.Sprintf("packer-%s", strings.ReplaceAll(instanceIp, ".", "-"))

	ui.Sayf("Saving ephemeral SSH key for local development purposes: %s", s.debugKeyPath)

	f, err := os.Create(s.debugKeyPath)
	if err != nil {
		ui.Errorf("Error saving debug key: %s", err)
		return multistep.ActionHalt
	}
	defer f.Close()

	// Write the key out
	if _, err := f.Write(s.sshKeyPair.PrivKey); err != nil {
		ui.Errorf("Error saving debug key: %s", err)
		return multistep.ActionHalt
	}

	// Chmod it so that it is SSH ready
	if runtime.GOOS != "windows" {
		if err := f.Chmod(0600); err != nil {
			ui.Errorf("Error setting permissions of debug key: %s", err)
			return multistep.ActionHalt
		}
	}

	return s.attachSshKeyPair(ctx, state)
}

func (s *StepSshKeyPair) Cleanup(state multistep.StateBag) {
	if s.Comm.SSHTemporaryKeyPairName == "" {
		return
	}

	var (
		cmd        string
		runcmdType string
		computerId = state.Get("instance_id").(string)
		osType     = state.Get("os_type").(string)
		ui         = state.Get("ui").(packersdk.Ui)
	)

	ui.Say("Trying to remove ephemeral keys from authorized_keys file...")

	switch strings.ToLower(osType) {
	case "windows":
		runcmdType = "RunPowerShellScript"
		cmd = concatWinCmd(fmt.Sprintf(`
Get-Content "$sshFolder\$sshFile" `+"`"+
			`  | Where-Object { -not $_.EndsWith(' %s') } `+"`"+
			`  | Out-File -FilePath "$sshFolder\$sshFile"`, s.Comm.SSHTemporaryKeyPairName))
	case "linux":
		// Do nothing because will use commonsteps.StepCleanupTempKeys in packer-plugin-sdk .
	default:
		ui.Errorf("Unsupported OS type: %s", osType)
	}

	if cmd != "" {
		ctx := context.TODO()
		common.RunCommand(ctx, state, &alieds.RunCommandRequest{
			RegionId:        alitea.String(s.RegionId),
			Type:            alitea.String(runcmdType),
			CommandContent:  alitea.String(cmd),
			ContentEncoding: alitea.String("PlainText"),
			DesktopId:       common.NilOrStringSlice(computerId),
			Timeout:         alitea.Int64(900),
		})
	}

	// Also remove the physical SSH private key
	ui.Say("Deleting ephemeral SSH key pair...")
	if err := os.Remove(s.debugKeyPath); err != nil {
		ui.Errorf("Error on delete ephemeral SSH key '%s': %s", s.debugKeyPath, err)
	}
}

func (s *StepSshKeyPair) attachSshKeyPair(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	computerId := state.Get("instance_id").(string)
	osType := state.Get("os_type").(string)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Attaching ephemeral SSH key pair to computer...")

	var (
		cmd        string
		pubKey     = strings.TrimSpace(string(s.sshKeyPair.PubKey))
		runcmdType string
	)
	switch strings.ToLower(osType) {
	case "windows":
		runcmdType = "RunPowerShellScript"
		cmd = concatWinCmd(fmt.Sprintf(`
# Write SSH authorized_keys
New-Item -Force -ItemType Directory -Path "$sshFolder"
Add-Content -Force -Path "$sshFolder\$sshFile" -Value "%s %s"`, pubKey, s.Comm.SSHTemporaryKeyPairName))
	case "linux":
		runcmdType = "RunShellScript"
		cmd = fmt.Sprintf(`
# Write SSH authorized_keys
mkdir -p ~/.ssh
chmod 700 ~/.ssh
touch ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
echo "%s %s" >> ~/.ssh/authorized_keys`, pubKey, s.Comm.SSHTemporaryKeyPairName)
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
		Timeout:         alitea.Int64(900),
	})
}

func concatWinCmd(cmd string) string {
	return fmt.Sprintf(`
# Check if admin account and setup variables for writing SSH authorized_keys file
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
If ($currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator))
{
  $sshFolder = "C:\ProgramData\ssh"
  $sshFile = "administrators_authorized_keys"
}
Else
{
  $sshFolder = "$env:USERPROFILE\.ssh"
  $sshFile = "authorized_keys"
}
%s`, cmd)
}
