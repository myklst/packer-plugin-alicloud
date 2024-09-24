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

## Outputs
|    Name     | Type   |
|-------------|--------|
|image_id     | string |
|image_name   | string |
|image_family | string |
|os_type      | string |
|architecture | string |

## Example
```
packer {
  required_plugins {
     alicloud = {
      source  = "github.com/myklst/alicloud"
      version = "~> 0.1"
    }
  }
}

data "alicloud-image" "test_image" {
  access_key = "v1-gastisthisisnotmyaccesskey"
  secret_key = "v9-adftthisfathisisnotmysecretkey"
  region  = "cn-hongkong"
  image_name = "aliyun_3_x64_20G_alibase_*.vhd"
}
```
