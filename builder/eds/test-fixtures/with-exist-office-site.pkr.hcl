# Test with: make dev && packer build ./builder/eds/test-fixtures/with-exist-office-site.pkr.hcl

packer {
  required_plugins {
    st-alicloud = {
      source  = "github.com/myklst/alicloud"
      version = "0.0.1"
    }
  }
}

variable "region" {
  type    = string
  default = "ap-southeast-3"
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

source "st-alicloud-eds" "test" {
  region     = var.region
  access_key = var.access_key
  secret_key = var.secret_key

  end_user {
    name  = "packer-user-01"
    email = "packer-user-01@example.com"
  }

  office_site {
    id = "ap-southeast-3+dir-3604693611"
  }

  computer_template {
    instance_type      = "eds.general.8c16g"
    root_disk_size_gib = 80
    user_disk_size_gib = [40]

    source_image_filter {
      image_id = "desktopimage-windows-11-64-asp"
      # image_id = "desktopimage-ubuntu-2204-asp"
    }
  }

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
    - id: Notepad++
      directives:
        description: Install Notepad++
      resource: Microsoft.WinGet.DSC/WinGetPackage
      settings:
        id: Notepad++.Notepad++
        source: winget
"@
Set-Content -Path $filePath -Value $multilineString
EOL
    encoding = "PlainText"
    timeout  = 30
  }

  artifact {
    image_name = "packer-test-with-exist-office-site"
  }
}

build {
  sources = ["source.st-alicloud-eds.test"]

  # provisioner "breakpoint" {
  #   disable = false
  #   note    = "this is a breakpoint"
  # }
}
