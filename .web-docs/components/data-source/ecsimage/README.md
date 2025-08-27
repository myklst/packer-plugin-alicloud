Type: `alicloud-ecsimage`

The AliCloud ECS Images data source will filter and fetch an ECS Image, and
output all the ECS Image information that will be then available to use in the
[AliCloud builders](/builder/eds).

-> **Note:** Data sources is a feature exclusively available to HCL2 templates.

## Basic Example

```hcl
data "alicloud-ecsimage" "example" {
  image_name   = ""
  image_family = ""
  is_public    = false
  source       = ""
  owner_id     = ""
  os_type      = ""
  architecture = ""
  usage        = ""

  tags {
    key = ""
    value = ""
  }

  tags {
    key = ""
    value = ""
  }
}
```

Note that the data source will fail unless *EXACTLY* ONE image is returned.

## Configuration Reference

**Optional:**

<!-- Code generated from the comments of the Image struct in datasource/ecsimage/data.go; DO NOT EDIT MANUALLY -->

- `image_id` (string) - Image Id

- `image_name` (string) - Image Name

- `image_family` (string) - Image Family

- `is_public` (bool) - Is Public

- `source` (string) - Source

- `owner_id` (int64) - Owner Id

- `os_type` (string) - OS Type

- `architecture` (string) - Architecture

- `tags` (ImageTags) - Tags

- `usage` (string) - Usage

<!-- End of code generated from the comments of the Image struct in datasource/ecsimage/data.go; -->


#### Tags Filter

<!-- Code generated from the comments of the ImageTag struct in datasource/ecsimage/data.go; DO NOT EDIT MANUALLY -->

- `key` (string) - Tag Key

- `value` (string) - Tag Value

<!-- End of code generated from the comments of the ImageTag struct in datasource/ecsimage/data.go; -->


## Output Data

<!-- Code generated from the comments of the DatasourceOutput struct in datasource/ecsimage/data.go; DO NOT EDIT MANUALLY -->

- `image` (\*Image) - Image

<!-- End of code generated from the comments of the DatasourceOutput struct in datasource/ecsimage/data.go; -->


## Authentication

The authentication for AliCloud Data Sources uses the same configuration options
as [AliCloud Builders](/builder/eds). To learn more about all of the available
authentication options please see
[AliCloud Builders authentication](/docs/builder/eds#authentication).

-> **Note:** The authentication session started by a data source is separate
from any authentication sessions started by an AliCloud builder. Users are
encouraged to use `variables` for defining and sharing configuration values
between datasources and builders.

Basic example of an AliCloud data source authentication using `assume_role`:

```hcl
data "amazon-secretsmanager" "basic-example" {
  name = "packer_test_secret"
  key  = "packer_test_key"

  assume_role {
      role_arn     = "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME"
      session_name = "SESSION_NAME"
      external_id  = "EXTERNAL_ID"
  }
}
```
