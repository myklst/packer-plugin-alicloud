// This file is auto-generated, don't edit it. Thanks.
package alicloud

import (
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// Description:
//
// # Initialize the Client with the AccessKey of the account
//
// @return Client
//
// @throws Exception
func CreateClient() (_result *openapi.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("ALICLOUD_ACCESS_KEY")),
		AccessKeySecret: tea.String(os.Getenv("ALICLOUD_SECRET_KEY")),
	}
	// See https://api.alibabacloud.com/product/Ecs.
	config.Endpoint = tea.String("ecs.cn-hongkong.aliyuncs.com")
	_result = &openapi.Client{}
	_result, _err = openapi.NewClient(config)
	return _result, _err
}

// Description:
//
// # API Info
//
// @param path - string Path parameters
//
// @return OpenApi.Params
func CreateApiInfo() (_result *openapi.Params) {
	params := &openapi.Params{
		// API Name
		Action: tea.String("DescribeImages"),
		// API Version
		Version: tea.String("2014-05-26"),
		// Protocol
		Protocol: tea.String("HTTPS"),
		// HTTP Method
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),
		// API PATH
		Pathname: tea.String("/"),
		// Request body content format
		ReqBodyType: tea.String("json"),
		// Response body content format
		BodyType: tea.String("json"),
	}
	_result = params
	return _result
}

func GetTEE() (_err error) {
	client, _err := CreateClient()
	if _err != nil {
		return _err
	}

	queries := map[string]interface{}{}
	queries["RegionId"] = tea.String("cn-hongkong")
	queries["ImageFamily"] = tea.String("acs:alibaba_cloud_linux_3_2104_lts_x64")
	// queries["IsSupportCloudinit"] = tea.Bool(true)
	queries["OSType"] = tea.String("linux")
	queries["Architecture"] = tea.String("x86_64")
	queries["Usage"] = tea.String("instance")

	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	// The return value is of Map type, and three types of data can be obtained from Map: response body, response headers, HTTP status code.
	params := CreateApiInfo()
	resp, _err := client.CallApi(params, request, runtime)
	if _err != nil {
		return _err
	}

	var out Output
	if body, ok := resp["body"].(map[string]interface{}); ok {
		mapstructure.Decode(body, &out)
	}

	if len(out.ImageList.Image) == 0 {
		return fmt.Errorf("no found matching filters")
	}

	if len(out.ImageList.Image) > 1 {
		return fmt.Errorf("query return more then one result, please specific search")
	}

	logrus.Info("Region :", out.Region)
	logrus.Info("ImageId :", out.ImageList.Image[0].ImageId)
	logrus.Info("ImageFamily :", out.ImageList.Image[0].ImageFamily)
	logrus.Info("OSType :", out.ImageList.Image[0].OSType)
	logrus.Info("Architecture :", out.ImageList.Image[0].Architecture)

	return nil
}

type Output struct {
	ImageList ImageList `mapstructure:"Images"`
	Region    string    `mapstructure:"RegionId"`
}

type ImageList struct {
	Image []Image `mapstructure:"Image"`
}

type Image struct {
	ImageId      string `mapstructure:"imageId"`
	ImageFamily  string `mapstructure:"imageFamily"`
	OSType       string `mapstructure:"osType"`
	Architecture string `mapstructure:"Architecture"`
}
