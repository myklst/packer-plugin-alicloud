package eds

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	// aliacc "github.com/myklst/packer-plugin-alicloud/builder/eds/acceptance"
)

//go:embed test-fixtures/basic.pkr.hcl
var testBuilderAcc_Basic string

//go:embed test-fixtures/with-exist-office-site.pkr.hcl
var testBuilderAcc_WithExistOfficeSite string

// Run with: PACKER_ACC=1 go test -count 1 -v ./builder/eds/builder_acc_test.go -timeout=120m
func TestAccBuilder_Eds_Basic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pkr  string
	}{
		{
			name: "basic",
			pkr: testBuilderAcc_Basic,
		},
		{
			name: "with-exist-office-site",
			pkr:  testBuilderAcc_WithExistOfficeSite,
		},
	}
	// alihelper := &aliacc.AlicloudHelper{}

	for _, test := range tests {
		testCase := &acctest.PluginTestCase{
			Name:     test.name,
			Template: test.pkr,
			Type:     "st-alicloud-eds",
			Setup: func() error {
				return nil
			},
			Teardown: func() error {
				return nil
			},
			Check: func(buildCommand *exec.Cmd, logfile string) error {
				if buildCommand.ProcessState != nil {
					if buildCommand.ProcessState.ExitCode() != 0 {
						return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
					}
				}

				logs, err := os.ReadFile(logfile)
				if err != nil {
					return fmt.Errorf("couldn't read logs from logfile %s: %s", logfile, err)
				}
				if strings.Contains(string(logs), "Uploading SSH public key") {
					return fmt.Errorf("SSH key was uploaded, but shouldn't have been")
				}

				return nil
			},
		}
		acctest.TestPlugin(t, testCase)
	}
}
