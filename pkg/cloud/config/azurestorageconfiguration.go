package config

import (
	"fmt"
)

type AzureStorageConfiguration struct {
	SubscriptionID string
	Account        string
	Container      string
	Path           string
	Cloud          string
	Configurer     AzureConfigurer
}

// Check ensures that all required fields are set, and throws an error if they are not
func (asc *AzureStorageConfiguration) Validate() error {

	if asc.Configurer == nil {
		return fmt.Errorf("AzureStorageConfiguration: missing configurer")
	}

	err := asc.Configurer.Validate()
	if err != nil {
		return err
	}

	if asc.SubscriptionID == "" {
		return fmt.Errorf("AzureStorageConfiguration: missing Subcription ID")
	}

	if asc.Account == "" {
		return fmt.Errorf("AzureStorageConfiguration: missing Account")
	}

	if asc.Container == "" {
		return fmt.Errorf("AzureStorageConfiguration: missing Container")
	}

	return nil
}

func (asc *AzureStorageConfiguration) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*AzureStorageConfiguration)
	if !ok {
		return false
	}

	if asc.Configurer != nil {
		if !asc.Configurer.Equals(thatConfig.Configurer) {
			return false
		}
	} else {
		if thatConfig.Configurer != nil {
			return false
		}
	}

	if asc.SubscriptionID != thatConfig.SubscriptionID {
		return false
	}

	if asc.Account != thatConfig.Account {
		return false
	}

	if asc.Container != thatConfig.Container {
		return false
	}

	if asc.Path != thatConfig.Path {
		return false
	}

	if asc.Cloud != thatConfig.Cloud {
		return false
	}

	return true
}

func (asc *AzureStorageConfiguration) Key() string {
	key := fmt.Sprintf("%s/%s", asc.SubscriptionID, asc.Container)
	// append container path to key if it exists
	if asc.Path != "" {
		key = fmt.Sprintf("%s/%s", key, asc.Path)
	}
	return key
}
