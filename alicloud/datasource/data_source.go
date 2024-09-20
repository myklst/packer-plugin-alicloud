// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type DatasourceOutput,Config
package datasource

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

type Datasource struct {
	config Config
}

type Config struct {
	AccessId     string `mapstructure:"access_key" required:"true"`
	AccessKey    string `mapstructure:"secret_key" required:"true"`
	Region       string `mapstructure:"region" required:"true"`
	ImageId      string `mapstructure:"image_id"`
	ImageName    string `mapstructure:"image_name"`
	ImageFamily  string `mapstructure:"image_family"`
	OSType       string `mapstructure:"os_type"`
	Architecture string `mapstructure:"architecture"`
	Usage        string `mapstructure:"usage"`
}

type DatasourceOutput struct {
	Region       string `mapstructure:"regionId"`
	ImageId      string `mapstructure:"imageId"`
	ImageFamily  string `mapstructure:"imageFamily"`
	OSType       string `mapstructure:"osType"`
	Architecture string `mapstructure:"architecture"`
}

type DescribeImagesOutput struct {
	ImageList ImageList `mapstructure:"Images"`
	Region    string    `mapstructure:"RegionId"`
}

type ImageList struct {
	Image []Image `mapstructure:"Image"`
}

type Image struct {
	ImageId      string `mapstructure:"ImageId"`
	ImageFamily  string `mapstructure:"ImageFamily"`
	OSType       string `mapstructure:"OsType"`
	Architecture string `mapstructure:"Architecture"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	if err := config.Decode(&d.config, nil, raws...); err != nil {
		return fmt.Errorf("error parsing configuration: %v", err)
	}

	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	client, err := openapi.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(d.config.AccessId),
		AccessKeySecret: tea.String(d.config.AccessKey),
		Endpoint:        tea.String(fmt.Sprintf("ecs.%s.aliyuncs.com", d.config.Region)),
	})
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	params := &openapi.Params{
		Action:      tea.String("DescribeImages"),
		Version:     tea.String("2014-05-26"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		Pathname:    tea.String("/"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}

	queries := map[string]interface{}{
		"ImageId":     tea.String(d.config.ImageId),
		"ImageName":   tea.String(d.config.ImageId),
		"RegionId":    tea.String(d.config.Region),
		"ImageFamily": tea.String(d.config.ImageFamily),
		"Usage":       tea.String(d.config.ImageFamily),
	}

	if d.config.OSType != "" {
		queries["OSType"] = tea.String(d.config.OSType)
	}

	if d.config.Architecture != "" {
		queries["Architecture"] = tea.String(d.config.Architecture)
	}

	// Make API request
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	resp, err := client.CallApi(params, request, runtime)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), fmt.Errorf("error querying AliCloud API: %v", err)
	}

	output, err := getFilteredImage(resp)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}

func getFilteredImage(resp map[string]interface{}) (DatasourceOutput, error) {
	var out DescribeImagesOutput
	var dataSourceOut DatasourceOutput

	if body, ok := resp["body"].(map[string]interface{}); ok {
		mapstructure.Decode(body, &out)
	}

	if len(out.ImageList.Image) == 0 {
		return dataSourceOut, fmt.Errorf("no found matching filters")
	}

	if len(out.ImageList.Image) > 1 {
		return dataSourceOut, fmt.Errorf("query return more then one result, please specific search")
	}
	output := DatasourceOutput{
		Region:       out.Region,
		ImageId:      out.ImageList.Image[0].ImageId,
		ImageFamily:  out.ImageList.Image[0].ImageFamily,
		OSType:       out.ImageList.Image[0].OSType,
		Architecture: out.ImageList.Image[0].Architecture,
	}
	return output, nil
}
