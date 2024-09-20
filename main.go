// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/myklst/packer-plugin-alicloud/datasource"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource("image", new(datasource.Datasource))
	// pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// if err := alicloud.GetTEE(); err != nil {
	// 	logrus.Error(err)
	// }
}
