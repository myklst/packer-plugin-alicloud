packer {
  required_plugins {
    st-alicloud = {
      source  = "github.com/myklst/alicloud"
      version = "0.0.1-dev"
    }
  }
}

variable "region" {
  type    = string
  default = "ap-southeast-1"
}

variable "access_key" {
  type    = string
  default = "${env("ALICLOUD_ACCESS_KEY")}"
}

variable "secret_key" {
  type      = string
  default   = "${env("ALICLOUD_SECRET_KEY")}"
  sensitive = true
}

source "alicloud-eds" "test" {
  region     = var.region
  access_key = var.access_key
  secret_key = var.secret_key

  end_user {
    name     = "packer-user-01"
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

//   user_commands {
//     type     = "RunPowerShellScript"
//     content  = <<EOL
// Install-PackageProvider -Name NuGet -MinimumVersion 2.8.5.201 -Force
// EOL
//     encoding = "PlainText"
//     timeout  = 180
//   }

//   user_commands {
//     type     = "RunPowerShellScript"
//     content  = <<EOL
// Install-Module -Name Microsoft.WinGet.Client -Force -Scope AllUsers
// Repair-WinGetPackageManager -AllUsers
// EOL
//     encoding = "PlainText"
//     timeout  = 180
//   }

  user_commands {
    type     = "RunPowerShellScript"
    content  = <<EOL
$filePath = "C:\packer.winget"
$multilineString = @"
# yaml-language-server: $schema=https://aka.ms/configuration-dsc-schema/0.2

###################################################################################
# This configuration will install the tools necessary for Backend Team on Windows #
###################################################################################

properties:
  configurationVersion: 0.2.0
  resources:
    - id: Atlassian Sourcetree
      directives:
        description: Install Atlassian Sourcetree
      resource: Microsoft.WinGet.DSC/WinGetPackage
      settings:
        id: Atlassian.Sourcetree
        source: winget

    - id: Notepad++
      directives:
        description: Install Notepad++
      resource: Microsoft.WinGet.DSC/WinGetPackage
      settings:
        id: Notepad++.Notepad++
        source: winget

    - id: Postman
      directives:
        description: Install Postman
      resource: Microsoft.WinGet.DSC/WinGetPackage
      settings:
        id: Postman.Postman
        source: winget

    - id: RocketChat
      directives:
        description: Install RocketChat
        securityContext: elevated
      resource: Microsoft.WinGet.DSC/WinGetPackage
      settings:
        id: RocketChat.RocketChat
        source: winget

    - id: Visual Studio
      directives:
        description: Install Microsoft VisualStudio 2022 Professional
        securityContext: elevated
      resource: Microsoft.WinGet.DSC/WinGetPackage
      settings:
        id: Microsoft.VisualStudio.2022.Professional
        source: winget

    - id: Visual Studio Code
      directives:
        description: Install Visual Studio Code
      resource: Microsoft.WinGet.DSC/WinGetPackage
      settings:
        id: Microsoft.VisualStudioCode
        source: winget
"@
Set-Content -Path $filePath -Value $multilineString
EOL
    encoding = "PlainText"
    timeout  = 30
  }

  user_commands {
    type     = "RunPowerShellScript"
    content  = <<EOL
winget.exe configure --accept-configuration-agreements C:/packer.winget
EOL
    encoding = "PlainText"
    timeout  = 900
  }
}

build {
  sources = ["source.alicloud-eds.test"]

# provisioner "breakpoint" {
#   disable = false
#   note    = "this is a breakpoint"
# }
}
