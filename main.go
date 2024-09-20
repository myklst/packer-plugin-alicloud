// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import "github.com/myklst/packer-plugin-alicloud/alicloud"

func main() {
	// pps := plugin.NewSet()
	// pps.RegisterDatasource("datasource", new(alicloud.Datasource))
	// pps.SetVersion(version.PluginVersion)
	// err := pps.Run()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	os.Exit(1)
	// }
	err := alicloud.GetTEE()
	if err != nil {
		panic(err)
	}
}
