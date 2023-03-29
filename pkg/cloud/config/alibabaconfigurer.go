package config

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

// AlibabaConfigurer provide *bssopenapi.Client for Alibaba cloud BOS for Billing related SDK calls
type AlibabaConfigurer interface {
	Config
	CreateAlibabaConfig() (*credentials.AccessKeyCredential, error)
	FromInterface(any) error
	ToDTO() AlibabaConfigurerDTO
}

// AlibabaAccessKey holds Alibaba credentials parsing from the service-key.json file.
type AlibabaAccessKey struct {
	AccessKeyId     string
	AccessKeySecret string
}

func (aak *AlibabaAccessKey) Validate() error {
	if aak.AccessKeyId == "" {
		return fmt.Errorf("AlibabaAccessKey: missing Access key ID")
	}
	if aak.AccessKeySecret == "" {
		return fmt.Errorf("AlibabaAccessKey: missing Access Key secret")
	}
	return nil
}

func (aak *AlibabaAccessKey) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*AlibabaAccessKey)
	if !ok {
		return false
	}

	if aak.AccessKeyId != thatConfig.AccessKeyId {
		return false
	}
	if aak.AccessKeySecret != thatConfig.AccessKeySecret {
		return false
	}
	return true
}

// CreateAlibabaConfig creates an
func (aak *AlibabaAccessKey) CreateAlibabaConfig() (*credentials.AccessKeyCredential, error) {
	err := aak.Validate()
	if err != nil {
		return nil, err
	}
	return &credentials.AccessKeyCredential{AccessKeyId: aak.AccessKeyId, AccessKeySecret: aak.AccessKeySecret}, nil
}
