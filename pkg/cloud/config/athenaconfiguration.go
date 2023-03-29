package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// AthenaConfiguration
type AthenaConfiguration struct {
	Bucket     string `json:"bucket"`
	Region     string `json:"region"`
	Database   string `json:"database"`
	Table      string `json:"table"`
	Workgroup  string `json:"workgroup"`
	Account    string `json:"account"`
	Configurer AWSConfigurer
}

func (ac *AthenaConfiguration) CreateAWSConfig() (aws.Config, error) {
	return ac.Configurer.CreateAWSConfig(ac.Region)
}

func (ac *AthenaConfiguration) Validate() error {

	// Validate Configurer
	if ac.Configurer == nil {
		return fmt.Errorf("AthenaConfiguration: missing configurer")
	}

	err := ac.Configurer.Validate()
	if err != nil {
		return fmt.Errorf("AthenaConfiguration: %s", err)
	}

	// Validate base properties
	if ac.Bucket == "" {
		return fmt.Errorf("AthenaConfiguration: missing bucket")
	}

	if ac.Region == "" {
		return fmt.Errorf("AthenaConfiguration: missing region")
	}

	if ac.Database == "" {
		return fmt.Errorf("AthenaConfiguration: missing database")
	}

	if ac.Table == "" {
		return fmt.Errorf("AthenaConfiguration: missing table")
	}

	if ac.Account == "" {
		return fmt.Errorf("AthenaConfiguration: missing account")
	}

	return nil
}

func (ac *AthenaConfiguration) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*AthenaConfiguration)
	if !ok {
		return false
	}

	if ac.Configurer != nil {
		if !ac.Configurer.Equals(thatConfig.Configurer) {
			return false
		}
	} else {
		if thatConfig.Configurer != nil {
			return false
		}
	}

	if ac.Bucket != thatConfig.Bucket {
		return false
	}

	if ac.Region != thatConfig.Region {
		return false
	}

	if ac.Database != thatConfig.Database {
		return false
	}

	if ac.Table != thatConfig.Table {
		return false
	}

	if ac.Workgroup != thatConfig.Workgroup {
		return false
	}

	if ac.Account != thatConfig.Account {
		return false
	}

	return true
}

func (ac *AthenaConfiguration) Key() string {
	return fmt.Sprintf("%s/%s", ac.Account, ac.Bucket)
}
