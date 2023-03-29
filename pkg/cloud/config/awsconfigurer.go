package config

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/opencost/opencost/pkg/util/json"
)

// AWSConfigurer provide aws.Config for AWS SDK calls
type AWSConfigurer interface {
	Config
	CreateAWSConfig(string) (aws.Config, error)
	FromInterface(any) error
	ToDTO() AWSConfigurerDTO
}

// AWSAccessKey holds AWS credentials and fulfils the awsV2.CredentialsProvider interface
type AWSAccessKey struct {
	ID     string
	Secret string
}

// Retrieve returns a set of awsV2 credentials using the AWSAccessKey's key and secret.
// This fulfils the awsV2.CredentialsProvider interface contract.
func (ak *AWSAccessKey) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     ak.ID,
		SecretAccessKey: ak.Secret,
	}, nil
}

func (ak *AWSAccessKey) Validate() error {
	if ak.ID == "" {
		return fmt.Errorf("AWSAccessKey: missing ID")
	}
	if ak.Secret == "" {
		return fmt.Errorf("AWSAccessKey: missing Secret")
	}
	return nil
}

func (ak *AWSAccessKey) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*AWSAccessKey)
	if !ok {
		return false
	}

	if ak.ID != thatConfig.ID {
		return false
	}
	if ak.Secret != thatConfig.Secret {
		return false
	}
	return true
}

// CreateAWSConfig creates an AWS SDK V2 Config for the credentials that it contains for the provided region
func (ak *AWSAccessKey) CreateAWSConfig(region string) (cfg aws.Config, err error) {
	err = ak.Validate()
	if err != nil {
		return cfg, err
	}
	// The AWS SDK v2 requires an object fulfilling the CredentialsProvider interface, which cloud.AWSAccessKey does
	cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(ak), config.WithRegion(region))
	if err != nil {
		return cfg, fmt.Errorf("failed to initialize AWS SDK config for region %s: %s", region, err)
	}
	return cfg, nil
}

// AWSServiceAccount uses pod annotations along with a service account to authenticate integrations
type AWSServiceAccount struct{}

// Check has nothing to check at this level, connection will fail if Pod Annotation and Service Account are not configured correctly
func (asa *AWSServiceAccount) Validate() error {
	return nil
}

func (asa *AWSServiceAccount) Equals(config Config) bool {
	if config == nil {
		return false
	}
	_, ok := config.(*AWSServiceAccount)
	if !ok {
		return false
	}

	return true
}

func (asa *AWSServiceAccount) CreateAWSConfig(region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return cfg, fmt.Errorf("failed to initialize AWS SDK config for region from annotation %s: %s", region, err)
	}
	return cfg, nil
}

// AWSAssumeRoleConfigurer is a wrapper for another configurer which adds an assumed role to the configuration
type AWSAssumeRoleConfigurer struct {
	configurer AWSConfigurer
	roleARN    string
}

func (arc *AWSAssumeRoleConfigurer) CreateAWSConfig(region string) (aws.Config, error) {
	cfg, _ := arc.configurer.CreateAWSConfig(region)
	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the roleARN.
	stsSvc := sts.NewFromConfig(cfg)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, arc.roleARN)
	cfg.Credentials = aws.NewCredentialsCache(creds)
	return cfg, nil
}

func (arc *AWSAssumeRoleConfigurer) Validate() error {
	if arc.configurer == nil {
		return fmt.Errorf("AWSAssumeRoleConfigurer: misisng base configurer")
	}
	err := arc.configurer.Validate()
	if err != nil {
		return err
	}

	if arc.roleARN == "" {
		return fmt.Errorf("AWSAssumeRoleConfigurer: misisng roleARN configuration")
	}

	return nil
}

func (arc *AWSAssumeRoleConfigurer) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*AWSAssumeRoleConfigurer)
	if !ok {
		return false
	}
	if arc.configurer != nil {
		if !arc.configurer.Equals(thatConfig.configurer) {
			return false
		}
	} else {
		if thatConfig.configurer != nil {
			return false
		}
	}

	if arc.roleARN != thatConfig.roleARN {
		return false
	}

	return true
}

func (arc *AWSAssumeRoleConfigurer) UnmarshalJSON(bytes []byte) error {
	var f interface{}

	err := json.Unmarshal(bytes, &f)
	if err != nil {
		return err
	}

	return arc.FromInterface(f)
}

func GetInterfaceValue[T any](fmap map[string]interface{}, key string) (T, error) {
	var value T
	interfaceValue, ok := fmap[key]
	if !ok {
		return value, fmt.Errorf("FromInterface: missing '%s' property", key)
	}
	typedValue, ok := interfaceValue.(T)
	if !ok {
		return value, fmt.Errorf("GetInterfaceValue: property '%s' had expected type '%T' but did not match", key, value)
	}
	return typedValue, nil
}
