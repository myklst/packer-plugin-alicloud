//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config,alipacker.AlicloudAccessConfig

package eds

import (
	"context"

	aliclient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	alieds20200930 "github.com/alibabacloud-go/ecd-20200930/v5/client"
	alieds20210308 "github.com/alibabacloud-go/eds-user-20210308/client"
	alitea "github.com/alibabacloud-go/tea/tea"
	alipacker "github.com/hashicorp/packer-plugin-alicloud/builder/ecs"

	"github.com/hashicorp/hcl/v2/hcldec"
	packercommon "github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"

	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

// The unique ID for this builder
const BuilderId = "myklst.alicloudeds"

type Config struct {
	packercommon.PackerConfig      `mapstructure:",squash"`
	alipacker.AlicloudAccessConfig `mapstructure:",squash"`
	RunConfig                      `mapstructure:",squash"`

	ctx interpolate.Context
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		PluginType:         BuilderId,
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	// Accumulate any errors
	var errs *packersdk.MultiError
	var warns []string

	errs = packersdk.MultiErrorAppend(errs, b.config.AlicloudAccessConfig.Prepare(&b.config.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, b.config.RunConfig.Prepare(&b.config.ctx)...)

	if errs != nil && len(errs.Errors) > 0 {
		return nil, warns, errs
	}

	// Return the placeholder for the generated data that will become
	// available to provisioners and post-processors.
	return getGeneratedDataList(), nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packersdk.Ui, hook packersdk.Hook) (packersdk.Artifact, error) {
	// Used for create EDS resources (except EDS users)
	alieds20200930, err := alieds20200930.NewClient(&aliclient.Config{
		RegionId:        common.NilOrString(b.config.AlicloudRegion),
		AccessKeyId:     common.NilOrString(b.config.AlicloudAccessKey),
		AccessKeySecret: common.NilOrString(b.config.AlicloudSecretKey),
	})
	if err != nil {
		return nil, err
	}

	// Used for create EDS users.
	alieds20210308, err := alieds20210308.NewClient(&aliclient.Config{
		RegionId:        alitea.String("ap-southeast-1"), // Cannot use b.config.AlicloudRegion because it's not supported by AliCloud EDS
		AccessKeyId:     common.NilOrString(b.config.AlicloudAccessKey),
		AccessKeySecret: common.NilOrString(b.config.AlicloudSecretKey),
	})
	if err != nil {
		return nil, err
	}

	// Setup the state bag and initial state for the steps
	state := new(multistep.BasicStateBag)
	state.Put("hook", hook)
	state.Put("ui", ui)
	state.Put("alieds20200930", alieds20200930)
	state.Put("alieds20210308", alieds20210308)
	state.Put("config", &b.config)
	generatedData := &packerbuilderdata.GeneratedData{State: state}

	steps := []multistep.Step{
		&StepPreValidate{
			RegionId:     b.config.AlicloudRegion,
			NewImageName: "", // TODO
		},
		&StepSourceImageInfo{
			RegionId:          b.config.AlicloudRegion,
			SourceImageFilter: &b.config.ComputerTemplate.SourceImageFilter,
		},
		&StepCloudComputerUser{
			Comm: &b.config.RunConfig.Comm,
			User: &b.config.RunConfig.EndUser,
		},
		&StepOfficeSite{
			RegionId:             b.config.AlicloudRegion,
			OfficeSiteId:         b.config.RunConfig.OfficeSite.Id,
			CidrBlock:            b.config.RunConfig.OfficeSite.CidrBlock,
			EnableInternetAccess: b.config.RunConfig.OfficeSite.InternetAccess.Enabled,
			InternetBandwidth:    b.config.RunConfig.OfficeSite.InternetAccess.Bandwidth,
			CenId:                b.config.RunConfig.OfficeSite.Cen.Id,
			CenOwnerId:           b.config.RunConfig.OfficeSite.Cen.OwnerId,
			CenVerifyCode:        b.config.RunConfig.OfficeSite.Cen.VerifyCode,
		},
		&StepCloudComputerTemplate{
			RegionId:                 b.config.AlicloudRegion,
			ComputerTemplateId:       b.config.RunConfig.ComputerTemplate.Id,
			InstanceType:             b.config.ComputerTemplate.InstanceType,
			RootDiskSizeGib:          b.config.ComputerTemplate.RootDiskSizeGib,
			RootDiskPerformanceLevel: b.config.ComputerTemplate.RootDiskPerformanceLevel,
			UserDiskSizeGib:          b.config.ComputerTemplate.UserDiskSizeGib,
			UserDiskPerformanceLevel: b.config.ComputerTemplate.UserDiskPerformanceLevel,
			Language:                 b.config.ComputerTemplate.Language,
		},
		&StepPolicyGroup{
			RegionId:      b.config.AlicloudRegion,
			PolicyGroupId: b.config.RunConfig.PolicyGroup.Id,
		},
		&StepCloudComputer{
			RegionId:                b.config.AlicloudRegion,
			ResourceGroupId:         b.config.RunConfig.ResourceGroupId,
			ComputerPoolId:          b.config.RunConfig.ComputerPoolId,
			ComputerTemplateId:      b.config.RunConfig.ComputerTemplate.Id,
			OfficeSiteId:            b.config.RunConfig.OfficeSite.Id,
			DesktopNameSuffix:       b.config.RunConfig.DesktopNameSuffix,
			PolicyGroupId:           b.config.RunConfig.PolicyGroup.Id,
			Hostname:                b.config.RunConfig.Hostname,
			DesktopMemberIp:         b.config.RunConfig.DesktopIp,
			VolumeEncryptionEnabled: b.config.RunConfig.VolumeEncryptionEnabled,
			VolumeEncryptionKey:     b.config.RunConfig.VolumeEncryptionKey,
		},
		&StepSshService{
			RegionId:  b.config.AlicloudRegion,
			EndUserId: b.config.RunConfig.EndUser.Name,
		},
		&StepSshKeyPair{
			Comm:      &b.config.RunConfig.Comm,
			RegionId:  b.config.AlicloudRegion,
			EndUserId: b.config.RunConfig.EndUser.Name,
		},
	}

	if len(b.config.RunConfig.UserCommands) > 0 {
		for _, command := range b.config.RunConfig.UserCommands {
			steps = append(steps, &StepUserCommand{
				RegionId:        b.config.AlicloudRegion,
				CommandType:     command.Type,
				CommandContent:  command.Content,
				ContentEncoding: command.Encoding,
				EndUserId:       b.config.RunConfig.EndUser.Name,
				CommandRole:     command.Role,
				Timeout:         command.Timeout,
			})
		}
	}

	// SSH only works when connected with our office network
	if b.config.RunConfig.OfficeSite.Id != "" {
		steps = append(steps, &communicator.StepConnect{
			Config:    &b.config.RunConfig.Comm,
			Host:      SshHost(),
			SSHConfig: b.config.RunConfig.Comm.SSHConfigFunc(),
		})
	}
	steps = append(steps, &StepSetGeneratedData{
		GeneratedData: generatedData,
	})
	steps = append(steps, &commonsteps.StepProvision{})
	steps = append(steps, &StepCreateImage{
		RegionId:     b.config.AlicloudRegion,
		NewImageName: b.config.RunConfig.Artifact.ImageName,
		Description:  b.config.RunConfig.Artifact.ImageDescription,
	})
	steps = append(steps, &commonsteps.StepCleanupTempKeys{
		Comm: &b.config.RunConfig.Comm,
	})

	// Run!
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	// If there was an error, return that
	if err, ok := state.GetOk("error"); ok {
		return nil, err.(error)
	}

	artifactStateData := map[string]interface{}{"generated_data": state.Get("generated_data")}
	if imageId, ok := state.GetOk("new_image_id"); ok {
		artifactStateData["new_image_id"] = imageId
	}

	artifact := &Artifact{
		RegionId:       b.config.AlicloudRegion,
		BuilderIdValue: BuilderId,
		StateData:      artifactStateData,
		Alieds20200930: alieds20200930,
		Alieds20210308: alieds20210308,
	}
	return artifact, nil
}
