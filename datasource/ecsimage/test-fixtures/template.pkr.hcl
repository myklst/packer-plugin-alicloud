packer {
  required_plugins {
    st-alicloud = {
      source  = "github.com/myklst/alicloud"
      version = "0.0.1-dev"
    }
  }
}

variable "region_id" {
  type    = string
  default = "cn-hongkong"
}

variable "access_key" {
  type    = string
  default = "${env("ALICLOUD_ACCESS_KEY")}"
}

variable "secret_key" {
  type    = string
  default =  "${env("ALICLOUD_SECRET_KEY")}"
}

variable "image_name" {
  type    = string
  default = "aliyun_3_x64_20G_alibase_*.vhd"
}

data "alicloud-ecsimage" "test_image" {
  region     = var.region_id
  access_key = var.access_key
  secret_key = var.secret_key

  image_name = var.image_name
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = ["source.null.basic-example"]

  provisioner "shell-local" {
    inline = [
      "echo image_id: ${data.alicloud-ecsimage.test_image.image.image_id}",
    ]
  }
}
