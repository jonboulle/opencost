package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type S3Configuration struct {
	Bucket        string
	Region        string
	Account       string
	AWSConfigurer AWSConfigurer
}

func (s3c *S3Configuration) Validate() error {
	// Validate Configurer
	if s3c.AWSConfigurer == nil {
		return fmt.Errorf("S3Configuration: missing configurer")
	}

	err := s3c.AWSConfigurer.Validate()
	if err != nil {
		return err
	}

	// Validate base properties
	if s3c.Bucket == "" {
		return fmt.Errorf("S3Configuration: missing bucket")
	}

	if s3c.Region == "" {
		return fmt.Errorf("S3Configuration: missing region")
	}

	if s3c.Account == "" {
		return fmt.Errorf("S3Configuration: missing account")
	}

	return nil
}

func (s3c *S3Configuration) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*S3Configuration)
	if !ok {
		return false
	}

	if !s3c.AWSConfigurer.Equals(thatConfig.AWSConfigurer) {
		return false
	}

	if s3c.Bucket != thatConfig.Bucket {
		return false
	}

	if s3c.Region != thatConfig.Region {
		return false
	}

	if s3c.Account != thatConfig.Account {
		return false
	}

	return true
}

func (s3c *S3Configuration) Key() string {
	return fmt.Sprintf("%s/%s", s3c.Account, s3c.Bucket)
}

func (s3c *S3Configuration) CreateAWSConfig() (aws.Config, error) {
	return s3c.AWSConfigurer.CreateAWSConfig(s3c.Region)
}
