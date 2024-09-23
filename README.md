# packer-plugin-alicloud

## Input
### Required:

- `access_key` **(string)**
- `secret_key` **(string)**
- `region_id` **(string)**

### Optional:

- `image_name` **(string)**
- `image_family` **(string)**
- `os_type` **(string)**
	- valid values: "windows", "linux"
- `architecture` **(string)**
	- valid values: "i386", "x86_64", "arm64"
- `usage` **(string)**
	- valid values: "instance", "none".
		- Instance: The image is already in use and running on an ECS instance.
		- None: The image is idle.

## Output
- `image_id` **(string)**
- `image_name` **(string)**
- `image_family` **(string)**
- `os_type` **(string)**
- `architecture` **(string)**

## Example
```
data "alicloud-image" "test_image" {
  access_key = "v1-gastisthisisnotmyaccesskey"
  secret_key = "v9-adftthisfathisisnotmysecretkey"
  region  = "cn-hongkong"
  image_name = "aliyun_3_x64_20G_alibase_*.vhd"
}
```
