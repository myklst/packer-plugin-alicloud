# packer-plugin-alicloud

## Input
### Required:

- `access_key` (string)
- `secret_key` (string)
- `region_id` (string)

### Optional:

- `image_name` (string)
- `image_family` (string)
- `os_type` (string): windows, linux
- `architecture` (string): i386, x86_64, arm64
- `usage` (string): instance, none. Instance - The image is already in use and running on an ECS instance,; None - The image is idle.

```
data "alicloud-image" "test_image" {
  access_key = var.access_key
  secret_key = var.secret_key
  region_id  = var.region_id
  image_name = var.image_name
}
```

## Output
- 	Region:       out.Region,
	ImageName:    out.ImageList.Image[0].ImageName,
	ImageFamily:  out.ImageList.Image[0].ImageFamily,
	OSType:       out.ImageList.Image[0].OSType,
	Architecture: out.ImageList.Image[0].Architecture,
