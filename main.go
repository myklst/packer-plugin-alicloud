package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"

	edsBuilder "github.com/myklst/packer-plugin-alicloud/builder/eds"
	ecsimageDatasource "github.com/myklst/packer-plugin-alicloud/datasource/ecsimage"
	"github.com/myklst/packer-plugin-alicloud/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("eds", new(edsBuilder.Builder))
	pps.RegisterDatasource("ecsimage", new(ecsimageDatasource.Datasource))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
