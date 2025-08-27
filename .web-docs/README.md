The AliCloud plugin can be used with HashiCorp Packer to create custom images on AliCloud. To achieve this, the plugin comes with
a builder, data source to build the Cloud Image depending on the strategy you want to use.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    amazon = {
      source  = "github.com/myklst/alicloud"
      version = "~> 0.1"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/myklst/alicloud
```

**Note: Update to Packer Plugin Installation**

With the new Packer release starting from version 1.14.0, the `packer init` command will automatically install official
plugins from the [HashiCorp release site.](https://releases.hashicorp.com/)

Going forward, to use newer versions of official Packer plugins, you'll need to upgrade to Packer version 1.14.0 or later.
If you're using an older version, you can still install plugins, but as a workaround, you'll need to
[manually install them using the CLI.](https://developer.hashicorp.com/packer/docs/plugins/install#manually-install-plugins-using-the-cli)

There is no change to the syntax or commands for installing plugins.

### Components

#### Builders

- [alicloud-eds](/builder/eds) - Create EDS Image by launching a source Image and re-packaging it into a new Image
  after provisioning.

#### Data sources

- [alicloud-ecs](/datasource/image) - Filter and fetch an AliCloud ECS Image to output all the Image information.

#### Post-Processors

*NONE*

### Authentication

The AliCloud provider offers a flexible means of providing credentials for authentication. The following methods are
supported, in this order, and explained below:

- Static credentials
- Environment variables
- Shared credentials file
- RAM Role

#### Static Credentials

Static credentials can be provided in the form of an access key id and secret.
These look like:

```hcl
source "alicloud-eds" "example" {
  region     = "ap-southeast-1"
  access_key = ""
  secret_key = ""
}
```

If you would like, you may also assume a role using the `ram_role_name`
configuration option. You must still have one of the valid credential resources
explained above, and your user must have permission to assume the role in
question. This is a way of running Packer with a more restrictive set of
permissions than your user.

AssumeRoleConfig lets users set configuration options for assuming a special
role when executing Packer.

Usage example:

HCL config example:

```HCL
source "alicloud-eds" "example" {
  ram_role_name    = ""
  ram_role_arn     = ""
  ram_session_name = ""
}
```

- `role_arn` (string) - Amazon Resource Name (ARN) of the IAM Role to assume.

- `duration_seconds` (int) - Number of seconds to restrict the assume role session duration.

- `external_id` (string) - The external ID to use when assuming the role. If omitted, no external
  ID is passed to the AssumeRole call.

- `policy` (string) - IAM Policy JSON describing further restricting permissions for the IAM
  Role being assumed.

- `policy_arns` ([]string) - Set of Amazon Resource Names (ARNs) of IAM Policies describing further
  restricting permissions for the IAM Role being

- `session_name` (string) - Session name to use when assuming the role.

- `tags` (map[string]string) - Map of assume role session tags.

- `transitive_tag_keys` ([]string) - Set of assume role session tag keys to pass to any subsequent sessions.

#### Environment variables

You can provide your credentials via the `ALICLOUD_ACCESS_KEY` and `ALICLOUD_SECRET_KEY`, environment variables,
representing your AliCloud Access Key and Secret Key, respectively. Note that setting your AliCloud credentials
using either these environment variables will override the use of `ALICLOUD_SHARED_CREDENTIALS_FILE` and
`ALICLOUD_PROFILE`. The `ALICLOUD_REGION` and `SECURITY_TOKEN` environment variables are also used, if applicable:

Usage:

    $ export ALICLOUD_ACCESS_KEY="anaccesskey"
    $ export ALICLOUD_SECRET_KEY="asecretkey"
    $ export ALICLOUD_REGION="us-west-2"
    $ packer build template.pkr.hcl

#### Shared Credentials file

You can use an AliCloud credentials file to specify your credentials. The default location is `$HOME/.aliyun/config.json`
on Linux and OS X, or `%USERPROFILE%.aliyun\config.json` for Windows users. If we fail to detect credentials inline, or
in the environment, the AliCloud Plugin will check this location. You can optionally specify a different location in the configuration by setting the environment with the `ALICLOUD_SHARED_CREDENTIALS_FILE` variable.

The format for the credentials file is like so

    [default]
    alicloud_access_key=<your access key id>
    alicloud_secret_key=<your secret access key>

You may also configure the profile to use by setting the `profile`
configuration option, or setting the `ALICLOUD_PROFILE` environment variable:

```hcl
source "alicloud-eds" "example" {
  region  = "us-east-1"
  profile = "customprofile"
}
```

#### RAM Task or Instance Role

Finally, the plugin will use credentials provided by the task's or instance's RAM role, if it has one.

This is a preferred approach over any other when running in AliCloud as you can avoid hard coding credentials. Instead
these are leased on-the-fly by the plugin, which reduces the chance of leakage.

The following policy document provides the minimal set permissions necessary for the AliCloud plugin to work:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecd:*"
      ],
      "Resource": "*"
    }
  ]
}
```

### Troubleshooting

#### TODO
