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
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

type Datasource struct {
	config Config
}

type Config struct {
	AccessKey    string `mapstructure:"access_key" required:"true"`
	SecretKey    string `mapstructure:"secret_key" required:"true"`
	Region       string `mapstructure:"region" required:"true"`
	ImageId      string `mapstructure:"image_id"`
	ImageName    string `mapstructure:"image_name"`
	ImageFamily  string `mapstructure:"image_family"`
	OsType       string `mapstructure:"os_type"`
	Architecture string `mapstructure:"architecture"`
	Usage        string `mapstructure:"usage"`
}

type DatasourceOutput struct {
	ImageId      string `mapstructure:"image_id"`
	ImageName    string `mapstructure:"image_name"`
	ImageFamily  string `mapstructure:"image_family"`
	OsType       string `mapstructure:"os_type"`
	Architecture string `mapstructure:"architecture"`
}

type DescribeImagesOutput struct {
	ImageList  ImageList `mapstructure:"Images"`
	TotalCount int       `mapstructure:"TotalCount"`
}

type ImageList struct {
	Image []Image `mapstructure:"Image"`
}

type Image struct {
	ImageId      string `mapstructure:"ImageId"`
	ImageName    string `mapstructure:"ImageName"`
	ImageFamily  string `mapstructure:"ImageFamily"`
	OsType       string `mapstructure:"OsType"`
	Architecture string `mapstructure:"Architecture"`
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

func (d *Datasource) Execute() (cty.Value, error) {
	client, err := openapi.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(d.config.AccessKey),
		AccessKeySecret: tea.String(d.config.SecretKey),
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
		"ImageName":   tea.String(d.config.ImageName),
		"RegionId":    tea.String(d.config.Region),
		"ImageFamily": tea.String(d.config.ImageFamily),
	}

	if d.config.OsType != "" {
		queries["OSType"] = tea.String(d.config.OsType)
	}

	if d.config.Architecture != "" {
		queries["Architecture"] = tea.String(d.config.Architecture)
	}

	if d.config.Usage != "" {
		queries["Usage"] = tea.String(d.config.Usage)
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
	var result DescribeImagesOutput
	var dataSourceOut DatasourceOutput

	if body, ok := resp["body"].(map[string]interface{}); ok {
		mapstructure.Decode(body, &result)
	}

	if result.TotalCount == 0 {
		return dataSourceOut, fmt.Errorf("no image found matching the filters")
	}

	if result.TotalCount > 1 {
		return dataSourceOut, fmt.Errorf("query return more then one result, please refine your search")
	}
	output := DatasourceOutput{
		ImageId:      result.ImageList.Image[0].ImageId,
		ImageName:    result.ImageList.Image[0].ImageName,
		ImageFamily:  result.ImageList.Image[0].ImageFamily,
		OsType:       result.ImageList.Image[0].OsType,
		Architecture: result.ImageList.Image[0].Architecture,
	}
	return output, nil
}
