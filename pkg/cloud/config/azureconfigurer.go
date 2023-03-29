package config

import (
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
)

type AzureConfigurer interface {
	Config
	GetBlobCredentials() (azblob.Credential, error)
	FromInterface(any) error
	ToDTO() AzureConfigurerDTO
}

type AzureAccessKey struct {
	AccessKey string
	Account   string
}

func (aak *AzureAccessKey) Validate() error {
	if aak.AccessKey == "" {
		return fmt.Errorf("AzureAccessKey: missing access key")
	}
	if aak.Account == "" {
		return fmt.Errorf("AzureAccessKey: missing account")
	}
	return nil
}

func (aak *AzureAccessKey) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*AzureAccessKey)
	if !ok {
		return false
	}

	if aak.AccessKey != thatConfig.AccessKey {
		return false
	}
	if aak.Account != thatConfig.Account {
		return false
	}

	return true
}

func (aak *AzureAccessKey) GetBlobCredentials() (azblob.Credential, error) {
	// Create a default request pipeline using your storage account name and account key.
	return azblob.NewSharedKeyCredential(aak.Account, aak.AccessKey)
}
