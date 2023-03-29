package config

import (
	"fmt"
	"github.com/opencost/opencost/pkg/util/json"
)

type AzureStorageConfigurationDTO struct {
	SubscriptionID string             `json:"subscriptionID"`
	Account        string             `json:"account"`
	Container      string             `json:"container"`
	Path           string             `json:"path"`
	Cloud          string             `json:"cloud"`
	Configurer     AzureConfigurerDTO `json:"configurer"`
}

// Sanitize redacts sensitive information in configuration
func (ascd *AzureStorageConfigurationDTO) Sanitize() {
	ascd.Configurer.Sanitize()
}

func (ascd *AzureStorageConfigurationDTO) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = ascd.FromInterface(f)
	if err != nil {
		return err
	}

	return nil
}

func (ascd *AzureStorageConfigurationDTO) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	subscriptionID, err := GetInterfaceValue[string](fmap, "subscriptionID")
	if err != nil {
		return fmt.Errorf("AzureStorageConfigurationDTO: FromInterface: %s", err.Error())
	}
	ascd.SubscriptionID = subscriptionID

	account, err := GetInterfaceValue[string](fmap, "account")
	if err != nil {
		return fmt.Errorf("AzureStorageConfigurationDTO: FromInterface: %s", err.Error())
	}
	ascd.Account = account

	container, err := GetInterfaceValue[string](fmap, "container")
	if err != nil {
		return fmt.Errorf("AzureStorageConfigurationDTO: FromInterface: %s", err.Error())
	}
	ascd.Container = container

	path, err := GetInterfaceValue[string](fmap, "path")
	if err != nil {
		return fmt.Errorf("AzureStorageConfigurationDTO: FromInterface: %s", err.Error())
	}
	ascd.Path = path

	cloud, err := GetInterfaceValue[string](fmap, "cloud")
	if err != nil {
		return fmt.Errorf("AzureStorageConfigurationDTO: FromInterface: %s", err.Error())
	}
	ascd.Cloud = cloud

	if confAny, ok := fmap["configurer"]; ok {
		var configurerDTO AzureConfigurerDTO
		err = configurerDTO.FromInterface(confAny)
		if err != nil {
			return fmt.Errorf("AthenaConfigurationDTO: FromInterface: %s", err.Error())
		}
		ascd.Configurer = configurerDTO
	}

	return nil
}

func (asc *AzureStorageConfiguration) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	subscriptionID, err := GetInterfaceValue[string](fmap, "subscriptionID")
	if err != nil {
		return fmt.Errorf("AzureStorageConfiguration: FromInterface: %s", err.Error())
	}
	asc.SubscriptionID = subscriptionID

	account, err := GetInterfaceValue[string](fmap, "account")
	if err != nil {
		return fmt.Errorf("AzureStorageConfiguration: FromInterface: %s", err.Error())
	}
	asc.Account = account

	container, err := GetInterfaceValue[string](fmap, "container")
	if err != nil {
		return fmt.Errorf("AzureStorageConfiguration: FromInterface: %s", err.Error())
	}
	asc.Container = container

	path, err := GetInterfaceValue[string](fmap, "path")
	if err != nil {
		return fmt.Errorf("AzureStorageConfiguration: FromInterface: %s", err.Error())
	}
	asc.Path = path

	cloud, err := GetInterfaceValue[string](fmap, "cloud")
	if err != nil {
		return fmt.Errorf("AzureStorageConfiguration: FromInterface: %s", err.Error())
	}
	asc.Cloud = cloud

	configurerAny, err := GetInterfaceValue[any](fmap, "configurer")
	if err != nil {
		return fmt.Errorf("AzureStorageConfiguration: FromInterface: %s", err.Error())
	}
	configurer, err := AzureConfigurerFromInterface(configurerAny)
	if err != nil {
		return fmt.Errorf("AzureStorageConfiguration: FromInterface: %s", err.Error())
	}
	asc.Configurer = configurer

	return nil
}

func (asc *AzureStorageConfiguration) ToConfigDTO() ConfigDTO {
	return asc.ToAzureStorageConfigurationDTO()
}

func (asc *AzureStorageConfiguration) ToAzureStorageConfigurationDTO() *AzureStorageConfigurationDTO {
	var configurer AzureConfigurerDTO
	if asc.Configurer != nil {
		configurer = asc.Configurer.ToDTO()
	}
	return &AzureStorageConfigurationDTO{
		SubscriptionID: asc.SubscriptionID,
		Account:        asc.Account,
		Container:      asc.Container,
		Path:           asc.Path,
		Cloud:          asc.Cloud,
		Configurer:     configurer,
	}
}

// UnmarshalJSON assumes data is save as an BigQueryConfigurationDTO
func (asc *AzureStorageConfiguration) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	dto := &AzureStorageConfigurationDTO{}
	err = dto.FromInterface(f)
	if err != nil {
		return err
	}

	// set values from dto
	asc.setValuesFromDTO(dto)
	return nil
}

func (asc *AzureStorageConfiguration) setValuesFromDTO(ascd *AzureStorageConfigurationDTO) {
	asc.SubscriptionID = ascd.SubscriptionID
	asc.Account = ascd.Account
	asc.Container = ascd.Container
	asc.Path = ascd.Path
	asc.Cloud = ascd.Cloud
	asc.Configurer = ascd.Configurer.ToAzureConfigurer()
}

// MarshalJSON converts to BigQueryConfigurationDTO and sanitize then Marshals as JSON
func (asc *AzureStorageConfiguration) MarshalJSON() ([]byte, error) {
	dto := asc.ToAzureStorageConfigurationDTO()
	dto.Sanitize()
	return json.Marshal(dto)
}
