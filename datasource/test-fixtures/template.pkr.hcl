packer {
  required_plugins {
     alicloud = {
      source  = "github.com/Kid-debug/alicloud"
      version = "v0.0.8"
    }
  }
}

variable "access_key" {
  type    = string
  default = "${env("ALICLOUD_ACCESS_KEY")}"
}

variable "secret_key" {
  type    = string
  default =  "${env("ALICLOUD_SECRET_KEY")}"
}

variable "region" {
  type    = string
  default = "cn-hongkong"
}

# Datasource block to retrieve AliCloud image details
data "alicloud-image" "test_image" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
  image_id   = "aliyun_3_x64_20G_alibase_*.vhd"
}

# Null builder to fulfill the requirement of a build block
source "null" "test" {
  communicator = "none"
}

build {
  sources = ["source.null.test"]

  provisioner "file" {
    content     = "This is a test provisioner"
    destination = "/tmp/test_file.txt"
  }
}
