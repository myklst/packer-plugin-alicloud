//go:generate packer-sdc mapstructure-to-hcl2 -type DatasourceOutput,Image,Config
package datasource

import (
	"encoding/json"
	"fmt"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

type Datasource struct {
	config Config
}

type Config struct {
	AccessKey    string            `mapstructure:"access_key" required:"true"`
	SecretKey    string            `mapstructure:"secret_key" required:"true"`
	Region       string            `mapstructure:"region" required:"true"`
	ImageId      string            `mapstructure:"image_id"`
	ImageName    string            `mapstructure:"image_name"`
	ImageFamily  string            `mapstructure:"image_family"`
	OsType       string            `mapstructure:"os_type"`
	Architecture string            `mapstructure:"architecture"`
	Usage        string            `mapstructure:"usage"`
	Tag          map[string]string `mapstructure:"tags"`
}

type DatasourceOutput struct {
	Images []Image `mapstructure:"images"`
}

type Image struct {
	ImageId      string `mapstructure:"image_id"`
	ImageName    string `mapstructure:"image_name"`
	ImageFamily  string `mapstructure:"image_family"`
	OSType       string `mapstructure:"os_type"`
	Architecture string `mapstructure:"architecture"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	if err := config.Decode(&d.config, nil, raws...); err != nil {
		return fmt.Errorf("error parsing configuration: %v", err)
	}

	var errs *packer.MultiError
	if d.config.AccessKey == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("access_key is missing"))
	}

	if d.config.SecretKey == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("secret_key is missing"))
	}

	if d.config.Region == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("region is missing"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func CreateClient(d *Datasource) (_result *ecs20140526.Client, _err error) {
	// The project code leakage may result in the leakage of AccessKey, posing a threat to the security of all resources under the account. The following code examples are for reference only.
	// It is recommended to use the more secure STS credential. For more credentials, please refer to: https://www.alibabacloud.com/help/en/alibaba-cloud-sdk-262060/latest/configure-credentials-378661.
	config := &openapi.Config{
		// Required, please ensure that the environment variables ALIBABA_CLOUD_ACCESS_KEY_ID is set.
		AccessKeyId: tea.String(d.config.AccessKey),
		// Required, please ensure that the environment variables ALIBABA_CLOUD_ACCESS_KEY_SECRET is set.
		AccessKeySecret: tea.String(d.config.SecretKey),
	}
	// See https://api.alibabacloud.com/product/Ecs.
	config.Endpoint = tea.String(fmt.Sprintf("ecs.%s.aliyuncs.com", d.config.Region))
	_result = &ecs20140526.Client{}
	_result, _err = ecs20140526.NewClient(config)
	return _result, _err
}

func (d *Datasource) Execute() (cty.Value, error) {
	client, err := CreateClient(d)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	var tags []*ecs20140526.DescribeImagesRequestTag
	for key := range d.config.Tag {
		tag := &ecs20140526.DescribeImagesRequestTag{
			Key:   tea.String(key),
			Value: tea.String(d.config.Tag[key]),
		}
		tags = append(tags, tag)
	}

	describeImagesRequest := &ecs20140526.DescribeImagesRequest{
		RegionId:    tea.String(d.config.Region),
		ImageName:   tea.String(d.config.ImageName),
		ImageId:     tea.String(d.config.ImageId),
		ImageFamily: tea.String(d.config.ImageFamily),
		Tag:         tags,
	}

	if d.config.OsType != "" {
		describeImagesRequest.OSType = tea.String(d.config.OsType)
	}

	if d.config.Architecture != "" {
		describeImagesRequest.Architecture = tea.String(d.config.Architecture)
	}

	if d.config.Usage != "" {
		describeImagesRequest.Usage = tea.String(d.config.Usage)
	}

	// Make API request
	var imgOutput DatasourceOutput
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// Copy the code to run, please print the return value of the API by yourself.
		resp, err := client.DescribeImagesWithOptions(describeImagesRequest, runtime)
		if err != nil {
			return err
		}

		imgOutput, err = getFilteredImage(resp)
		if err != nil {
			return err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// Only a printing example. Please be careful about exception handling and do not ignore exceptions directly in engineering projects.
		// print error message
		// Please click on the link below for diagnosis.
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend := m["Recommend"]
			return cty.NullVal(cty.EmptyObject), fmt.Errorf(tryErr.Error(), recommend)
		}

		return cty.NullVal(cty.EmptyObject), fmt.Errorf("%s", tryErr.Error())
	}

	return hcl2helper.HCL2ValueFromConfig(imgOutput, d.OutputSpec()), nil
}

func getFilteredImage(resp *ecs20140526.DescribeImagesResponse) (DatasourceOutput, error) {
	var dataSourceOut DatasourceOutput

	if *resp.Body.TotalCount == 0 {
		return dataSourceOut, fmt.Errorf("no image found matching the filters")
	}

	for _, img := range resp.Body.Images.Image {
		image := Image{
			ImageId:      *img.ImageId,
			ImageName:    *img.ImageName,
			ImageFamily:  *img.ImageFamily,
			OSType:       *img.OSType,
			Architecture: *img.Architecture,
		}
		dataSourceOut.Images = append(dataSourceOut.Images, image)
	}

	return dataSourceOut, nil

}
