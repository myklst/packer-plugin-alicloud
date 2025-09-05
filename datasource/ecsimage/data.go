//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type ImageTag,Image,Config,DatasourceOutput
package ecsimage

import (
	"context"
	"fmt"
	"time"

	aliopenapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliecs20140526 "github.com/alibabacloud-go/ecs-20140526/v7/client"
	alitea "github.com/alibabacloud-go/tea/tea"
	alipacker "github.com/hashicorp/packer-plugin-alicloud/builder/ecs"

	"github.com/hashicorp/hcl/v2/hcldec"
	packercommon "github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"
	"github.com/hashicorp/packer-plugin-sdk/template/config"

	"github.com/zclconf/go-cty/cty"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

type ImageTag struct {
	TagKey   string `mapstructure:"key"`
	TagValue string `mapstructure:"value"`
}
type ImageTags []ImageTag

type Image struct {
	ImageId      string    `mapstructure:"image_id"`
	ImageName    string    `mapstructure:"image_name"`
	ImageFamily  string    `mapstructure:"image_family"`
	IsPublic     bool      `mapstructure:"is_public"`
	Source       string    `mapstructure:"source"`
	OwnerId      int64     `mapstructure:"owner_id"`
	OSType       string    `mapstructure:"os_type"`
	Architecture string    `mapstructure:"architecture"`
	Tags         ImageTags `mapstructure:"tags"`
	Usage        string    `mapstructure:"usage"`
}

type Config struct {
	packercommon.PackerConfig      `mapstructure:",squash"`
	alipacker.AlicloudAccessConfig `mapstructure:",squash"`
	Image                          `mapstructure:",squash"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	Image *Image `mapstructure:"image"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec { return d.config.FlatMapstructure().HCL2Spec() }

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	if err := config.Decode(&d.config, nil, raws...); err != nil {
		return fmt.Errorf("error parsing configuration: %v", err)
	}

	var errs *packer.MultiError
	if d.config.AlicloudAccessConfig.AlicloudRegion == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("region is missing"))
	}

	if d.config.AlicloudAccessConfig.AlicloudAccessKey == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("access_key is missing"))
	}

	if d.config.AlicloudAccessConfig.AlicloudSecretKey == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("secret_key is missing"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (d *Datasource) Execute() (cty.Value, error) {
	config := &aliopenapi.Config{
		AccessKeyId:     alitea.String(d.config.AlicloudAccessConfig.AlicloudAccessKey),
		AccessKeySecret: alitea.String(d.config.AlicloudAccessConfig.AlicloudSecretKey),
	}
	config.Endpoint = alitea.String(fmt.Sprintf("ecs.%s.aliyuncs.com", d.config.AlicloudAccessConfig.AlicloudRegion))
	client, err := aliecs20140526.NewClient(config)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	var (
		describeImagesRequest = aliecs20140526.DescribeImagesRequest{
			RegionId:        common.NilOrString(d.config.AlicloudAccessConfig.AlicloudRegion),
			ImageId:         common.NilOrString(d.config.ImageId),
			ImageName:       common.NilOrString(d.config.ImageName),
			ImageFamily:     common.NilOrString(d.config.ImageFamily),
			IsPublic:        common.NilOrBool(d.config.IsPublic),
			ImageOwnerAlias: common.NilOrString(d.config.Source),
			ImageOwnerId:    alitea.Int64(d.config.OwnerId),
			OSType:          common.NilOrString(d.config.OSType),
			Architecture:    common.NilOrString(d.config.Architecture),
			Usage:           common.NilOrString(d.config.Usage),
		}
		tags []*aliecs20140526.DescribeImagesRequestTag
	)
	for _, tag := range d.config.Tags {
		tag := &aliecs20140526.DescribeImagesRequestTag{
			Key:   alitea.String(tag.TagKey),
			Value: alitea.String(tag.TagValue),
		}
		tags = append(tags, tag)
	}
	if len(tags) > 0 {
		describeImagesRequest.Tag = tags
	}

	var resp *aliecs20140526.DescribeImagesResponse
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, _ := common.IsRetryableError(err)
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(context.TODO(), func(ctx context.Context) error {
		resp, err = client.DescribeImages(&describeImagesRequest)
		if err == nil {
			if *resp.Body.TotalCount == 0 {
				return fmt.Errorf("no image found matching the filters")
			}

			if *resp.Body.TotalCount > 1 {
				return fmt.Errorf("multiple images found matching the filters")
			}
		}
		return err
	})
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	output := &DatasourceOutput{}
	for _, img := range resp.Body.Images.Image {
		var tags []ImageTag
		for _, imgtag := range img.Tags.Tag {
			tag := ImageTag{
				TagKey:   alitea.StringValue(imgtag.TagKey),
				TagValue: alitea.StringValue(imgtag.TagValue),
			}
			tags = append(tags, tag)
		}

		output.Image = &Image{
			ImageId:      alitea.StringValue(img.ImageId),
			ImageName:    alitea.StringValue(img.ImageName),
			ImageFamily:  alitea.StringValue(img.ImageFamily),
			IsPublic:     alitea.BoolValue(img.IsPublic),
			Source:       alitea.StringValue(img.ImageOwnerAlias),
			OwnerId:      alitea.Int64Value(img.ImageOwnerId),
			OSType:       alitea.StringValue(img.OSType),
			Architecture: alitea.StringValue(img.Architecture),
			Tags:         tags,
		}
	}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
