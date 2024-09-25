variable "access_key" {
  type    = string
  default = "${env("ALICLOUD_ACCESS_KEY")}"
}

variable "secret_key" {
  type    = string
  default =  "${env("ALICLOUD_SECRET_KEY")}"
}

variable "region_id" {
  type    = string
  default = "cn-hongkong"
}

variable "image_name" {
  type    = string
  default = "aliyun_3_x64_20G_alibase_*.vhd"
}

data "alicloud-image" "test_image" {
  access_key = var.access_key
  secret_key = var.secret_key
  region  = var.region_id
  image_name = var.image_name
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = ["source.null.basic-example"]

  provisioner "shell-local" {
    inline = [
      "echo image_id: ${data.alicloud-image.test_image.images[0].image_id}",
    ]

  }
}
