package config

import (
	"fmt"

	"github.com/opencost/opencost/pkg/log"
)

const (
	AlibabaAccessKeyType string = "AlibabaAccessKey"
)

type AlibabaConfigurerDTO struct {
	Type   string         `json:"type"`
	Values map[string]any `json:"values,omitempty"`
}

// Sanitize redacts sensitive information from the Values
func (a *AlibabaConfigurerDTO) Sanitize() {
	switch a.Type {
	case AlibabaAccessKeyType:
		if value, ok := a.Values["secret"]; ok && value != "" {
			a.Values["secret"] = Redacted
		}
	}
}

func (a *AlibabaConfigurerDTO) ToAlibabaConfigurer() AlibabaConfigurer {
	var configurer AlibabaConfigurer
	switch a.Type {
	case AlibabaAccessKeyType:
		configurer = &AlibabaAccessKey{}
	default:
		log.Errorf("AlibabaConfigurerDTO: ToAlibabaConfigurer: invalid type: %s", a.Type)
		return nil
	}
	err := configurer.FromInterface(a.Values)
	if err != nil {
		log.Errorf("AlibabaConfigurerDTO: ToAlibabaConfigurer: failed to convert to AlibabaConfigurer: %s", err.Error())
	}
	return configurer
}

func (a *AlibabaConfigurerDTO) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	typeOf, err := GetInterfaceValue[string](fmap, "type")
	if err != nil {
		return fmt.Errorf("AlibabaConfigurerDTO: %s", err.Error())
	}
	a.Type = typeOf

	_, ok := fmap["values"]
	if ok {
		values, err := GetInterfaceValue[map[string]any](fmap, "values")
		if err != nil {
			return fmt.Errorf("AlibabaConfigurerDTO: %s", err.Error())
		}
		a.Values = values
	}

	return nil
}

func (ak *AlibabaAccessKey) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	id, err := GetInterfaceValue[string](fmap, "id")
	if err != nil {
		return fmt.Errorf("AlibabaAccessKey: %s", err.Error())
	}
	ak.AccessKeyId = id

	secret, err := GetInterfaceValue[string](fmap, "secret")
	if err != nil {
		return fmt.Errorf("AlibabaAccessKey: %s", err.Error())
	}
	ak.AccessKeySecret = secret

	return nil
}

func (ak *AlibabaAccessKey) ToDTO() AlibabaConfigurerDTO {
	values := make(map[string]any, 2)
	values["id"] = ak.AccessKeyId
	values["secret"] = ak.AccessKeySecret
	return AlibabaConfigurerDTO{
		Type:   AlibabaAccessKeyType,
		Values: values,
	}
}
