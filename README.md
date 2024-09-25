# packer-plugin-alicloud

## Inputs

### Required:
|    Name   |  Type  |
|-----------|--------|
|access_key | string |
|secret_key | string |
|region_id  | string |

### Optional:

|    Name     | Type   |     Valid Value     | Description                                                                                            |
|-------------|--------|---------------------|--------------------------------------------------------------------------------------------------------|
|image_id     | string | any string          | ID of the image.                                                                                       |
|image_name   | string | any string          | Name of the image.                                                                                     |
|image_family | string | any string          | Famaily of the image.                                                                                  |
|os_type      | string | windows, linux      | OS type of the image.                                                                                  |
|architecture | string | i386, x86_64, arm64 | Architectre of the images.                                                                             |
|usage        | string | instance, none      |- Instance: The image is already in use and running on an ECS instance. <br> - None: The image is idle. |
|tags         | map    | map of string       | Tags of the image.                                                                                     |

## Outputs
|    Name     | Type           |
|-------------|----------------|
|images       | <pre>list of object([{<br>  image_id     = string<br>  image_name   = string<br>  image_family = string<br>  os_type      = string<br>  architecture = string<br>  tags         = list of object([{<br>    key   = string<br>    value = string<br>  }])<br>}])</pre> |


## Example
```
packer {
  required_plugins {
     st-alicloud = {
      source  = "github.com/myklst/alicloud"
      version = "~> 0.1"
    }
  }
}

data "st-alicloud-image" "ecs_images" {
  access_key = "v1-gastisthisisnotmyaccesskey"
  secret_key = "v9-adftthisfathisisnotmysecretkey"
  region  = "cn-hongkong"
  image_name = "ThisIsMyImageName_*.vhd"

  tag = {
    status = "activated"
  }
}

locals {
  prodImageID = compact(flatten([for v in data.alicloud-image.ecs_images.images : [
    for tag in v.tags : [
     tag.key == "env" && tag.value == "prod"? v.image_id : null
    ]
  ]]))

  devImageID = compact(flatten([for v in data.alicloud-image.ecs_images.images : [
    for tag in v.tags : [
     tag.key == "env" && tag.value == "dev"? v.image_id : null
    ]
  ]]))
}
```
