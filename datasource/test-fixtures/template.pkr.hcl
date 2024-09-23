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
  region     = "cn-hongkong"
  image_name   = "aliyun_3_x64_20G_alibase_20240819.vhd"
}

// locals{
//   image_id = data.alicloud-image.test_image.image_id
// }

# Null builder to fulfill the requirement of a build block
source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = ["source.null.basic-example"]

  provisioner "shell-local" {
    inline = [
      "echo image_id: ${data.alicloud-image.test_image.image_id}",
    ]
  }
}
