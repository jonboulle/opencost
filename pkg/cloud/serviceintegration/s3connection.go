package serviceintegration

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	config2 "github.com/opencost/opencost/pkg/cloud/config"
)

type S3Connection struct {
	config2.S3Configuration
}

func (s3c *S3Connection) Equals(config config2.Config) bool {
	thatConfig, ok := config.(*S3Connection)
	if !ok {
		return false
	}

	return s3c.S3Configuration.Equals(&thatConfig.S3Configuration)
}

func (s3c *S3Connection) GetS3Client() (*s3.Client, error) {
	cfg, err := s3c.CreateAWSConfig()
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}

func (s3c *S3Connection) ListObjects(cli *s3.Client) (*s3.ListObjectsOutput, error) {
	objs, err := cli.ListObjects(context.TODO(), &s3.ListObjectsInput{
		Bucket: aws.String(s3c.Bucket),
	})
	if err != nil {
		return nil, err
	}
	return objs, err
}
