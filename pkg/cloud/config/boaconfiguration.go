package config

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

// BoaConfiguration is the BSS open API configuration for Alibaba's Billing information
type BoaConfiguration struct {
	Account    string `json:"account"`
	Region     string `json:"region"`
	Configurer AlibabaConfigurer
}

func (bc *BoaConfiguration) CreateAlibabaConfig() (*credentials.AccessKeyCredential, error) {
	return bc.Configurer.CreateAlibabaConfig()
}

func (bc *BoaConfiguration) Validate() error {
	// Validate Configurer
	if bc.Configurer == nil {
		return fmt.Errorf("BoaConfiguration: missing configurer")
	}

	err := bc.Configurer.Validate()
	if err != nil {
		return err
	}

	// Validate base properties
	if bc.Region == "" {
		return fmt.Errorf("BoaConfiguration: missing region")
	}

	if bc.Account == "" {
		return fmt.Errorf("BoaConfiguration: missing account")
	}
	return nil
}

func (bc *BoaConfiguration) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*BoaConfiguration)
	if !ok {
		return false
	}

	if bc.Configurer != nil {
		if !bc.Configurer.Equals(thatConfig.Configurer) {
			return false
		}
	} else {
		if thatConfig.Configurer != nil {
			return false
		}
	}

	if bc.Region != thatConfig.Region {
		return false
	}
	return true
}

func (bc *BoaConfiguration) Key() string {
	return fmt.Sprintf("%s/%s", bc.Account, bc.Region)
}
