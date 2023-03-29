package config

import (
	"fmt"
	"github.com/opencost/opencost/pkg/log"
)

const (
	AzureAccessKeyType string = "AzureAccessKey"
)

type AzureConfigurerDTO struct {
	Type   string         `json:"type"`
	Values map[string]any `json:"values,omitempty"`
}

// Sanitize redacts sensitive information from the Values
func (a *AzureConfigurerDTO) Sanitize() {
	switch a.Type {
	case AzureAccessKeyType:
		if value, ok := a.Values["accessKey"]; ok && value != "" {
			a.Values["accessKey"] = Redacted
		}
	}
}

func (a *AzureConfigurerDTO) ToAzureConfigurer() AzureConfigurer {
	var configurer AzureConfigurer
	switch a.Type {
	case AzureAccessKeyType:
		configurer = &AzureAccessKey{}
	default:
		log.Errorf("AzureConfigurerDTO: ToAzureConfigurer: invalid type: %s", a.Type)
		return nil
	}
	err := configurer.FromInterface(a.Values)
	if err != nil {
		log.Errorf("AzureConfigurerDTO: ToAzureConfigurer: failed to convert to AzureConfigurer: %s", err.Error())
	}
	return configurer
}

func (a *AzureConfigurerDTO) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	typeOf, err := GetInterfaceValue[string](fmap, "type")
	if err != nil {
		return fmt.Errorf("AzureConfigurerDTO: %s", err.Error())
	}
	a.Type = typeOf

	_, ok := fmap["values"]
	if ok {
		values, err := GetInterfaceValue[map[string]any](fmap, "values")
		if err != nil {
			return fmt.Errorf("AzureConfigurerDTO: %s", err.Error())
		}
		a.Values = values
	}

	return nil
}

func AzureConfigurerFromInterface(f any) (AzureConfigurer, error) {
	if f == nil {
		return nil, nil
	}
	fmap, ok := f.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("AzureConfigurerFromInterface: could not case interface as map")
	}

	var configurer AzureConfigurer
	if _, ok = fmap["accesskey"]; ok {
		configurer = &AzureAccessKey{}
		err := configurer.FromInterface(fmap)
		if err != nil {
			return nil, err
		}
		return configurer, nil
	}
	if _, ok = fmap["type"]; ok {
		dto := &AzureConfigurerDTO{}
		err := dto.FromInterface(fmap)
		if err != nil {
			return nil, err
		}
		return dto.ToAzureConfigurer(), nil
	}

	return nil, fmt.Errorf("AzureConfigurerFromInterface: interface is not a valid configurer format")
}

func (aak *AzureAccessKey) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	accessKey, err := GetInterfaceValue[string](fmap, "accessKey")
	if err != nil {
		return fmt.Errorf("AzureAccessKey: %s", err.Error())
	}
	aak.AccessKey = accessKey

	account, err := GetInterfaceValue[string](fmap, "account")
	if err != nil {
		return fmt.Errorf("AzureAccessKey: %s", err.Error())
	}
	aak.Account = account

	return nil
}

func (aak *AzureAccessKey) ToDTO() AzureConfigurerDTO {
	values := make(map[string]any, 2)
	values["accessKey"] = aak.AccessKey
	values["account"] = aak.Account
	return AzureConfigurerDTO{
		Type:   AzureAccessKeyType,
		Values: values,
	}
}
