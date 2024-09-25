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

  // tags ={
  //   env = "basic"
  // }
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
    // inline = [
    //   "echo rysn_image_id: ${local.rsynImageID[0]}, webserver_image_id: ${local.webserverImageID[0]}",
    // ]
  }
}

// locals {
//   rsynImageID = compact(flatten([for v in data.alicloud-image.test_image.images : [
//     for tag in v.tags : [
//      tag.key == "usage" && tag.value == "rsyn"? v.image_id : null
//     ]
//   ]]))
//   webserverImageID = compact(flatten([for v in data.alicloud-image.test_image.images : [
//     for tag in v.tags : [
//      tag.key == "usage" && tag.value == "webserver"? v.image_id : null
//     ]
//   ]]))
// }
