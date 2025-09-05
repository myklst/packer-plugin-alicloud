package eds

import (
	"context"
	"time"

	alieds20200930 "github.com/alibabacloud-go/ecd-20200930/v5/client"
	alieds20210308 "github.com/alibabacloud-go/eds-user-20210308/client"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/retry"
	"github.com/myklst/packer-plugin-alicloud/builder/common"
)

// Artifact is an artifact implementation that contains built AMIs.
type Artifact struct {
	RegionId string

	// BuilderId is the unique ID for the builder that created this Image
	BuilderIdValue string

	// StateData should store data such as GeneratedData
	// to be shared with post-processors
	StateData map[string]interface{}

	// AliCloud EDS connection for performing API stuff.
	Alieds20200930 *alieds20200930.Client
	Alieds20210308 *alieds20210308.Client
}

func (a *Artifact) BuilderId() string {
	return a.BuilderIdValue
}

func (*Artifact) Files() []string {
	// We have no files
	return nil
}

func (a *Artifact) Id() string {
	return a.State("new_image_id").(string)
}

func (a *Artifact) String() string {
	return a.State("new_image_id").(string)
}

func (a *Artifact) State(name string) interface{} {
	if _, ok := a.StateData[name]; ok {
		return a.StateData[name]
	}

	switch name {
	default:
		return nil
	}
}

func (a *Artifact) Destroy() error {
	var (
		err     error
		errors  = make([]error, 0)
		imageId = a.State("new_image_id").(string)
	)
	err = retry.Config{
		Tries: 8,
		ShouldRetry: func(err error) bool {
			retryable, _ := common.IsRetryableError(err)
			return retryable
		},
		RetryDelay: (&retry.Backoff{InitialBackoff: 1 * time.Second, MaxBackoff: 30 * time.Second, Multiplier: 2}).Linear,
	}.Run(context.TODO(), func(ctx context.Context) error {
		_, err = a.Alieds20200930.DeleteImages(&alieds20200930.DeleteImagesRequest{
			RegionId: common.NilOrString(a.RegionId),
			ImageId:  common.NilOrStringSlice(imageId),
		})
		return err
	})
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		if len(errors) == 1 {
			return errors[0]
		} else {
			return &packersdk.MultiError{Errors: errors}
		}
	}

	return nil
}
