---
Type: |
  `alicloud-eds`
Artifact BuilderId: |
  `myklst.alicloudeds`
---

The `alicloud-eds` Packer builder is able to create Images for use in
[AliCloud EDS](https://www.alibabacloud.com/en/product/cloud-desktop?_p_lc=1).

This builder builds an Image by launching an EDS computer from a source Image,
provisioning that running machine, and then creating an Image from that machine.
This is all done in your own AliCloud account. The builder will create temporary
cloud computer user, office site (if `office_site.id` does not specified), etc.
that provide it temporary access to the instance while the image is being created.
This simplifies configuration quite a bit.

The builder does _not_ manage Images. Once it creates an Image and stores it in
your account, it is up to you to use, delete, etc. the Image.

-> **Note:** Temporary resources are, by default, all created with the
prefix `packer-`. This can be useful if you want to restrict the security groups
and key pairs Packer is able to operate on.

## Configuration Reference

There are many configuration options available for the builder. In addition to
the items listed here, you will want to look at the general configuration
references for [Access](#access-configuration),
[Cloud Computer](#cloud-computer-configuration),
[Cloud Computer User](#cloud-computer-user-configuration),
[Office Site](#office-site-configuration),
[Computer Template](#computer-template-configuration),
[Policy Group](#policy-group-configuration),
[User Command](#user-command-configuration) and
[Communicator](#communicator-configuration)
configuration references, which are necessary for this build to succeed and can
be found further down the page.

### Access Configuration

**Required:**

[HERE](https://github.com/hashicorp/packer-plugin-alicloud/blob/main/docs-partials/builder/ecs/AlicloudAccessConfig-required.mdx)

**Optional:**

[HERE](https://github.com/hashicorp/packer-plugin-alicloud/blob/main/docs-partials/builder/ecs/AlicloudAccessConfig-not-required.mdx)

### Cloud Computer Configuration

**Required:**

*NONE*

**Optional:**

<!-- Code generated from the comments of the RunConfig struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `resource_group_id` (string) - Resource Group Id

- `computer_pool_id` (string) - Computer Pool Id

- `desktop_name` (string) - Desktop Name

- `desktop_name_suffix` (bool) - Desktop Name Suffix

- `hostname` (string) - Hostname

- `desktop_ip` (string) - Desktop Ip

- `volume_encryption_enabled` (bool) - Volume Encryption Enabled

- `volume_encryption_key` (string) - Volume Encryption Key

- `end_user` (EdsUser) - End User

- `office_site` (EdsOfficeSite) - Office Site

- `computer_template` (EdsComputerTemplate) - Computer Template

- `policy_group` (EdsPolicyGroup) - Policy Group

- `user_commands` (EdsUserCommands) - User Commands

<!-- End of code generated from the comments of the RunConfig struct in builder/eds/run_config.go; -->


### Cloud Computer User Configuration

**Required:**

<!-- Code generated from the comments of the EdsUser struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `name` (string) - Name

- `email` (string) - Email

<!-- End of code generated from the comments of the EdsUser struct in builder/eds/run_config.go; -->


**Optional:**

<!-- Code generated from the comments of the EdsUser struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `phone` (string) - Phone

- `owner_type` (string) - Owner Type

- `org_id` (string) - Org Id

- `remark` (string) - Remark

- `real_nick_name` (string) - Real Nick Name

<!-- End of code generated from the comments of the EdsUser struct in builder/eds/run_config.go; -->


### Office Site Configuration

**Optional:**

<!-- Code generated from the comments of the EdsOfficeSite struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `id` (string) - Id

- `name` (string) - Name

- `cidr_block` (string) - Cidr Block

- `internet_access` (EdsOfficeSiteInternetAccess) - Internet Access

- `cen` (EdsOfficeSiteCen) - Cen

<!-- End of code generated from the comments of the EdsOfficeSite struct in builder/eds/run_config.go; -->


##### CEN Configuration

<!-- Code generated from the comments of the EdsOfficeSiteCen struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `id` (string) - Id

- `owner_id` (int64) - Owner Id

- `verify_code` (string) - Verify Code

<!-- End of code generated from the comments of the EdsOfficeSiteCen struct in builder/eds/run_config.go; -->


##### Internet Access Configuration

<!-- Code generated from the comments of the EdsOfficeSiteInternetAccess struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `enabled` (bool) - Enabled

- `bandwidth` (int32) - Bandwidth

<!-- End of code generated from the comments of the EdsOfficeSiteInternetAccess struct in builder/eds/run_config.go; -->


### Computer Template Configuration

**Required:**

<!-- Code generated from the comments of the EdsComputerTemplate struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `instance_type` (string) - Instance Type

- `root_disk_size_gib` (int) - Root Disk Size Gib

<!-- End of code generated from the comments of the EdsComputerTemplate struct in builder/eds/run_config.go; -->


**Optional:**

<!-- Code generated from the comments of the EdsComputerTemplate struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `id` (string) - Id

- `name` (string) - Name

- `description` (string) - Description

- `source_image_filter` (EdsImageFilter) - Source Image Filter

- `root_disk_performance_level` (string) - Root Disk Performance Level

- `user_disk_size_gib` ([]int32) - User Disk Size Gib

- `user_disk_performance_level` (string) - User Disk Performance Level

- `language` (string) - Language

<!-- End of code generated from the comments of the EdsComputerTemplate struct in builder/eds/run_config.go; -->


##### Source Image Filter Configuration

<!-- Code generated from the comments of the EdsImageFilter struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `image_id` (string) - The ID of the image.

- `image_name` (string) - The image name.

- `image_version` (string) - The image version.

- `image_type` (string) - The type of the image. Avaiable options: [SYSTEM, CUSTOM]

- `image_status` (string) - The status of the image. Avaiable options: [Creating, Available,
  CreateFailed]

- `protocol_type` (string) - The protocol type. Available options: [HDX, ASP]

- `desktop_instance_type` (string) - The instance type of the cloud computer. You can call the
  DescribeDesktopTypes operation to obtain the parameter value.

- `session_type` (string) - The session type. Avaiable options: [SINGLE_SESSION, MULTIPLE_SESSION]

- `gpu_category` (bool) - Specifies whether the images are GPU-accelerated images. Available
  options: [true, false]

- `gpu_driver_version` (string) - The version of the GPU driver.

- `os_type` (string) - The type of the operating system of the images.Available options:
  [Linux, Windows]

- `language_type` (string) - The language of the OS. Available options: [en-US, zh-HK, zh-CN, ja-JP]

<!-- End of code generated from the comments of the EdsImageFilter struct in builder/eds/run_config.go; -->


### Policy Group Configuration

**Required:**

*NONE*

**Optional:**

<!-- Code generated from the comments of the EdsPolicyGroup struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `id` (string) - Id

- `name` (string) - Name

<!-- End of code generated from the comments of the EdsPolicyGroup struct in builder/eds/run_config.go; -->


### User Command Configuration

**Required:**

<!-- Code generated from the comments of the EdsUserCommand struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `type` (string) - Type

- `content` (string) - Content

- `encoding` (string) - Encoding

<!-- End of code generated from the comments of the EdsUserCommand struct in builder/eds/run_config.go; -->


**Optional:**

<!-- Code generated from the comments of the EdsUserCommand struct in builder/eds/run_config.go; DO NOT EDIT MANUALLY -->

- `end_user_id` (string) - End User Id

- `role` (string) - Role

- `timeout` (uint64) - Timeout

<!-- End of code generated from the comments of the EdsUserCommand struct in builder/eds/run_config.go; -->


### Communicator Configuration

**Optional:**

<!-- Code generated from the comments of the Config struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `communicator` (string) - Packer currently supports three kinds of communicators:
  
  -   `none` - No communicator will be used. If this is set, most
      provisioners also can't be used.
  
  -   `ssh` - An SSH connection will be established to the machine. This
      is usually the default.
  
  -   `winrm` - A WinRM connection will be established.
  
  In addition to the above, some builders have custom communicators they
  can use. For example, the Docker builder has a "docker" communicator
  that uses `docker exec` and `docker cp` to execute scripts and copy
  files.

- `pause_before_connecting` (duration string | ex: "1h5m2s") - We recommend that you enable SSH or WinRM as the very last step in your
  guest's bootstrap script, but sometimes you may have a race condition
  where you need Packer to wait before attempting to connect to your
  guest.
  
  If you end up in this situation, you can use the template option
  `pause_before_connecting`. By default, there is no pause. For example if
  you set `pause_before_connecting` to `10m` Packer will check whether it
  can connect, as normal. But once a connection attempt is successful, it
  will disconnect and then wait 10 minutes before connecting to the guest
  and beginning provisioning.

<!-- End of code generated from the comments of the Config struct in communicator/config.go; -->


<!-- Code generated from the comments of the SSH struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `ssh_host` (string) - The address to SSH to. This usually is automatically configured by the
  builder.

- `ssh_port` (int) - The port to connect to SSH. This defaults to `22`.

- `ssh_username` (string) - The username to connect to SSH with. Required if using SSH.

- `ssh_password` (string) - A plaintext password to use to authenticate with SSH.

- `ssh_ciphers` ([]string) - This overrides the value of ciphers supported by default by Golang.
  The default value is [
    "aes128-gcm@openssh.com",
    "chacha20-poly1305@openssh.com",
    "aes128-ctr", "aes192-ctr", "aes256-ctr",
  ]
  
  Valid options for ciphers include:
  "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com",
  "chacha20-poly1305@openssh.com",
  "arcfour256", "arcfour128", "arcfour", "aes128-cbc", "3des-cbc",

- `ssh_clear_authorized_keys` (bool) - If true, Packer will attempt to remove its temporary key from
  `~/.ssh/authorized_keys` and `/root/.ssh/authorized_keys`. This is a
  mostly cosmetic option, since Packer will delete the temporary private
  key from the host system regardless of whether this is set to true
  (unless the user has set the `-debug` flag). Defaults to "false";
  currently only works on guests with `sed` installed.

- `ssh_key_exchange_algorithms` ([]string) - If set, Packer will override the value of key exchange (kex) algorithms
  supported by default by Golang. Acceptable values include:
  "curve25519-sha256@libssh.org", "ecdh-sha2-nistp256",
  "ecdh-sha2-nistp384", "ecdh-sha2-nistp521",
  "diffie-hellman-group14-sha1", and "diffie-hellman-group1-sha1".

- `ssh_certificate_file` (string) - Path to user certificate used to authenticate with SSH.
  The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_pty` (bool) - If `true`, a PTY will be requested for the SSH connection. This defaults
  to `false`.

- `ssh_timeout` (duration string | ex: "1h5m2s") - The time to wait for SSH to become available. Packer uses this to
  determine when the machine has booted so this is usually quite long.
  Example value: `10m`.
  This defaults to `5m`, unless `ssh_handshake_attempts` is set.

- `ssh_disable_agent_forwarding` (bool) - If true, SSH agent forwarding will be disabled. Defaults to `false`.

- `ssh_handshake_attempts` (int) - The number of handshakes to attempt with SSH once it can connect.
  This defaults to `10`, unless a `ssh_timeout` is set.

- `ssh_bastion_host` (string) - A bastion host to use for the actual SSH connection.

- `ssh_bastion_port` (int) - The port of the bastion host. Defaults to `22`.

- `ssh_bastion_agent_auth` (bool) - If `true`, the local SSH agent will be used to authenticate with the
  bastion host. Defaults to `false`.

- `ssh_bastion_username` (string) - The username to connect to the bastion host.

- `ssh_bastion_password` (string) - The password to use to authenticate with the bastion host.

- `ssh_bastion_interactive` (bool) - If `true`, the keyboard-interactive used to authenticate with bastion host.

- `ssh_bastion_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with the
  bastion host. The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_bastion_certificate_file` (string) - Path to user certificate used to authenticate with bastion host.
  The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_file_transfer_method` (string) - `scp` or `sftp` - How to transfer files, Secure copy (default) or SSH
  File Transfer Protocol.
  
  **NOTE**: Guests using Windows with Win32-OpenSSH v9.1.0.0p1-Beta, scp
  (the default protocol for copying data) returns a a non-zero error code since the MOTW
  cannot be set, which cause any file transfer to fail. As a workaround you can override the transfer protocol
  with SFTP instead `ssh_file_transfer_method = "sftp"`.

- `ssh_proxy_host` (string) - A SOCKS proxy host to use for SSH connection

- `ssh_proxy_port` (int) - A port of the SOCKS proxy. Defaults to `1080`.

- `ssh_proxy_username` (string) - The optional username to authenticate with the proxy server.

- `ssh_proxy_password` (string) - The optional password to use to authenticate with the proxy server.

- `ssh_keep_alive_interval` (duration string | ex: "1h5m2s") - How often to send "keep alive" messages to the server. Set to a negative
  value (`-1s`) to disable. Example value: `10s`. Defaults to `5s`.

- `ssh_read_write_timeout` (duration string | ex: "1h5m2s") - The amount of time to wait for a remote command to end. This might be
  useful if, for example, packer hangs on a connection after a reboot.
  Example: `5m`. Disabled by default.

- `ssh_remote_tunnels` ([]string) - Remote tunnels forward a port from your local machine to the instance.
  Format: ["REMOTE_PORT:LOCAL_HOST:LOCAL_PORT"]
  Example: "9090:localhost:80" forwards localhost:9090 on your machine to port 80 on the instance.

- `ssh_local_tunnels` ([]string) - Local tunnels forward a port from the instance to your local machine.
  Format: ["LOCAL_PORT:REMOTE_HOST:REMOTE_PORT"]
  Example: "8080:localhost:3000" allows the instance to access your local machine’s port 3000 via localhost:8080.

<!-- End of code generated from the comments of the SSH struct in communicator/config.go; -->


<!-- Code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `temporary_key_pair_type` (string) - `dsa` | `ecdsa` | `ed25519` | `rsa` ( the default )
  
  Specifies the type of key to create. The possible values are 'dsa',
  'ecdsa', 'ed25519', or 'rsa'.
  
  NOTE: DSA is deprecated and no longer recognized as secure, please
  consider other alternatives like RSA or ED25519.

- `temporary_key_pair_bits` (int) - Specifies the number of bits in the key to create. For RSA keys, the
  minimum size is 1024 bits and the default is 4096 bits. Generally, 3072
  bits is considered sufficient. DSA keys must be exactly 1024 bits as
  specified by FIPS 186-2. For ECDSA keys, bits determines the key length
  by selecting from one of three elliptic curve sizes: 256, 384 or 521
  bits. Attempting to use bit lengths other than these three values for
  ECDSA keys will fail. Ed25519 keys have a fixed length and bits will be
  ignored.
  
  NOTE: DSA is deprecated and no longer recognized as secure as specified
  by FIPS 186-5, please consider other alternatives like RSA or ED25519.

<!-- End of code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; -->


- `ssh_keypair_name` (string) - If specified, this is the key that will be used for SSH with the
  machine. The key must match a key pair name loaded up into the remote.
  By default, this is blank, and Packer will generate a temporary keypair
  unless [`ssh_password`](#ssh_password) is used.
  [`ssh_private_key_file`](#ssh_private_key_file) or
  [`ssh_agent_auth`](#ssh_agent_auth) must be specified when
  [`ssh_keypair_name`](#ssh_keypair_name) is utilized.


- `ssh_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with SSH.
  The `~` can be used in path and will be expanded to the home directory
  of current user.


- `ssh_agent_auth` (bool) - If true, the local SSH agent will be used to authenticate connections to
  the source instance. No temporary keypair will be created, and the
  values of [`ssh_password`](#ssh_password) and
  [`ssh_private_key_file`](#ssh_private_key_file) will be ignored. The
  environment variable `SSH_AUTH_SOCK` must be set for this option to work
  properly.


## Basic Example

Here is a basic example. It is completely valid except for the access keys:

**HCL2**

```hcl
source "alicloud-eds" "example" {
  region        = "ap-southeast-11"

  end_user {
    name     = "packer-user-01"
    password = "Password1234"
    email    = "packer-user-01@example.com"
  }

  office_site {
    internet_access {
      enabled   = true
      bandwidth = 10
    }
  }

  computer_template {
    instance_type      = "eds.general.8c16g"
    root_disk_size_gib = 80
    user_disk_size_gib = [40]

    source_image_filter {
      image_id = "desktopimage-windows-11-64-asp"
    }
  }

  user_commands {
    type     = "RunPowerShellScript"
    content  = <<EOL
New-Item -Path "C:\packer-test-01" -ItemType Directory
EOL
    encoding = "PlainText"
    timeout  = 30
  }

  user_commands {
    type     = "RunPowerShellScript"
    content  = <<EOL
New-Item -Path "C:\packer-test-02" -ItemType Directory
EOL
    encoding = "PlainText"
    timeout  = 30
  }
}
```

-> **Note:** Packer can also read the access key and secret access key from
environmental variables. See the configuration reference in the section above
for more information on what environmental variables Packer will look for.

## Accessing the Instance to Debug

If you need to access the instance to debug for some reason, run this builder
with the `-debug` flag. In debug mode, the builder will output the Wuying username
and password. You can use this information to access the instance as it is running.

## Build template data

In configuration directives marked as a template engine above, the following
variables are available:

*NONE*

## Build Shared Information Variables

This builder generates data that are shared with provisioner and post-processor
via build function of
[template engine](/packer/docs/templates/legacy_json_templates/engine) for JSON
and [contextual variables](/packer/docs/templates/hcl_templates/contextual-variables)
for HCL2.

The generated variables available for this builder are:

*NONE*

## Connecting to Windows instances using WinRM

If you want to launch a Windows instance and connect using WinRM, you will need
to configure WinRM on that instance. The following is a basic powershell script
that can be supplied to AliCloud EDS using the "user_commands" option. It enables
WinRM via HTTPS on port 5986, and creates a self-signed certificate to use to
connect. If you are using a certificate from a CA, rather than creating a
self-signed certificate, you can omit the "winrm_insecure" option mentioned below.

autogenerated_password_https_bootstrap.txt

```powershell
<powershell>

# MAKE SURE IN YOUR PACKER CONFIG TO SET:
#
#
#    "winrm_username": "Administrator",
#    "winrm_insecure": true,
#    "winrm_use_ssl": true,
#
#

Write-Output "Running User Data Script"
Write-Host "(host) Running User Data Script"

Set-ExecutionPolicy Unrestricted -Scope LocalMachine -Force -ErrorAction Ignore

# Don't set this before Set-ExecutionPolicy as it throws an error
$ErrorActionPreference = "stop"

# Remove HTTP listener
Remove-Item -Path WSMan:\Localhost\listener\listener* -Recurse

# Create a self-signed certificate to let ssl work
$Cert = New-SelfSignedCertificate -CertstoreLocation Cert:\LocalMachine\My -DnsName "packer"
New-Item -Path WSMan:\LocalHost\Listener -Transport HTTPS -Address * -CertificateThumbPrint $Cert.Thumbprint -Force

# WinRM
Write-Output "Setting up WinRM"
Write-Host "(host) setting up WinRM"

cmd.exe /c winrm quickconfig -q
cmd.exe /c winrm set "winrm/config" '@{MaxTimeoutms="1800000"}'
cmd.exe /c winrm set "winrm/config/winrs" '@{MaxMemoryPerShellMB="1024"}'
cmd.exe /c winrm set "winrm/config/service" '@{AllowUnencrypted="true"}'
cmd.exe /c winrm set "winrm/config/client" '@{AllowUnencrypted="true"}'
cmd.exe /c winrm set "winrm/config/service/auth" '@{Basic="true"}'
cmd.exe /c winrm set "winrm/config/client/auth" '@{Basic="true"}'
cmd.exe /c winrm set "winrm/config/service/auth" '@{CredSSP="true"}'
cmd.exe /c winrm set "winrm/config/listener?Address=*+Transport=HTTPS" "@{Port=`"5986`";Hostname=`"packer`";CertificateThumbprint=`"$($Cert.Thumbprint)`"}"
cmd.exe /c netsh advfirewall firewall set rule group="remote administration" new enable=yes
cmd.exe /c netsh advfirewall firewall add rule name="Port 5986" dir=in action=allow protocol=TCP localport=5986 profile=any
cmd.exe /c net stop winrm
cmd.exe /c sc config winrm start= auto
cmd.exe /c net start winrm

</powershell>
```

You'll notice that this config does not define a user or password; instead,
Packer will ask AWS to provide a random password that it generates
automatically. The following config will work with the above template:
