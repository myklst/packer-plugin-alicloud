// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

//go:embed test-fixtures/template.pkr.hcl
var testDatasourceHCL2Basic string

// Run with: PACKER_ACC=1 go test -count 1 -v ./datasource/data_acc_test.go  -timeout=120m
func TestAccAliCloudDatasource(t *testing.T) {
	// Define the test case
	testCase := &acctest.PluginTestCase{
		Name: "alicloud_datasource_basic_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testDatasourceHCL2Basic,
		Type:     "alicloud-image",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil && buildCommand.ProcessState.ExitCode() != 0 {
				return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable to open logfile: %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := os.ReadFile(logfile)
			if err != nil {
				return fmt.Errorf("Unable to read logfile: %s", logfile)
			}

			logsString := string(logsBytes)

			expectedLog := "null.basic-example: image_id: aliyun_3_x64_20G_alibase_.*.vhd"

			if matched, _ := regexp.MatchString(expectedLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected value %q", logsString)
			}

			return nil
		},
	}

	// Run the test case
	acctest.TestPlugin(t, testCase)
}
