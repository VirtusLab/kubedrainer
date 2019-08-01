package aws

import (
	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	// DefaultProfile is the default profile to be used when
	// loading configuration from the config files if another profile name
	// is not provided.
	DefaultProfile = session.DefaultSharedConfigProfile
)

// SessionConfig returns AWS API client session and config with given region and profile
func SessionConfig(region, profile string) (*session.Session, *aws.Config, error) {
	// AWS_DEFAULT_PROFILE environment variable can be also used to set profile
	config := aws.NewConfig().WithRegion(region)
	awsSession, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config:  *config,
	})

	return awsSession, config, errors.Wrap(err)
}
