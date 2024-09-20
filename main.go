// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/myklst/packer-plugin-alicloud/alicloud/datasource"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource("datasource", new(datasource.Datasource))
	// pps.SetVersion("v0.0.2")
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
