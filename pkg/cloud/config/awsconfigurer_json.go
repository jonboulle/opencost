package config

import (
	"fmt"
	"github.com/opencost/opencost/pkg/log"
)

const Redacted = "REDACTED"

const (
	AWSAccessKeyType      string = "AWSAccessKey"
	AWSServiceAccountType        = "AWSServiceAccount"
	AWSAssumeRoleType            = "AWSAssumeRole"
)

type AWSConfigurerDTO struct {
	Type   string         `json:"type"`
	Values map[string]any `json:"values,omitempty"`
}

// Sanitize redacts sensitive information from the Values
func (a *AWSConfigurerDTO) Sanitize() {
	switch a.Type {
	case AWSAccessKeyType:
		if value, ok := a.Values["secret"]; ok && value != "" {
			a.Values["secret"] = Redacted
		}
	case AWSAssumeRoleType:
		if value, ok := a.Values["configurer"]; ok && value != nil {
			configurer, _ := value.(AWSConfigurerDTO)
			configurer.Sanitize()
			a.Values["configurer"] = configurer
		}
	}
}

func (a *AWSConfigurerDTO) ToAWSConfigurer() AWSConfigurer {
	var configurer AWSConfigurer
	switch a.Type {
	case AWSAccessKeyType:
		configurer = &AWSAccessKey{}
	case AWSServiceAccountType:
		configurer = &AWSServiceAccount{}
	case AWSAssumeRoleType:
		configurer = &AWSAssumeRoleConfigurer{}
	default:
		log.Errorf("AWSConfigurerDTO: ToAWSConfigurer: invalid type: %s", a.Type)
		return nil
	}
	err := configurer.FromInterface(a.Values)
	if err != nil {
		log.Errorf("AWSConfigurerDTO: ToAWSConfigurer: failed to convert to AWSConfigurer: %s", err.Error())
	}
	return configurer
}

func (a *AWSConfigurerDTO) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	typeOf, err := GetInterfaceValue[string](fmap, "type")
	if err != nil {
		return fmt.Errorf("AWSConfigurerDTO: %s", err.Error())
	}
	a.Type = typeOf

	_, ok := fmap["values"]
	if ok {
		values, err := GetInterfaceValue[map[string]any](fmap, "values")
		if err != nil {
			return fmt.Errorf("AWSConfigurerDTO: %s", err.Error())
		}
		a.Values = values
	}

	return nil
}

func AWSConfigurerFromInterface(f any) (AWSConfigurer, error) {
	if f == nil {
		return nil, nil
	}
	fmap, ok := f.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("AWSConfigurerFromInterface: could not case interface as map")
	}
	if len(fmap) == 0 {
		return &AWSServiceAccount{}, nil
	}
	if _, ok = fmap["roleARN"]; ok {
		awsConfigurer := &AWSAssumeRoleConfigurer{}
		err := awsConfigurer.FromInterface(fmap)
		if err != nil {
			return nil, err
		}
		return awsConfigurer, nil
	}
	if _, ok = fmap["id"]; ok {
		awsConfigurer := &AWSAccessKey{}
		err := awsConfigurer.FromInterface(fmap)
		if err != nil {
			return nil, err
		}
		return awsConfigurer, nil
	}
	if _, ok = fmap["type"]; ok {
		dto := &AWSConfigurerDTO{}
		err := dto.FromInterface(fmap)
		if err != nil {
			return nil, err
		}
		return dto.ToAWSConfigurer(), nil
	}

	return nil, fmt.Errorf("AWSConfigurerFromInterface: interface is not a valid configurer format")
}

func (ak *AWSAccessKey) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	id, err := GetInterfaceValue[string](fmap, "id")
	if err != nil {
		return fmt.Errorf("AWSAccessKey: %s", err.Error())
	}
	ak.ID = id

	secret, err := GetInterfaceValue[string](fmap, "secret")
	if err != nil {
		return fmt.Errorf("AWSAccessKey: %s", err.Error())
	}
	ak.Secret = secret

	return nil
}

func (ak *AWSAccessKey) ToDTO() AWSConfigurerDTO {
	values := make(map[string]any, 2)
	values["id"] = ak.ID
	values["secret"] = ak.Secret
	return AWSConfigurerDTO{
		Type:   AWSAccessKeyType,
		Values: values,
	}
}

func (asa *AWSServiceAccount) FromInterface(f interface{}) error {
	return nil
}

func (asa *AWSServiceAccount) ToDTO() AWSConfigurerDTO {
	return AWSConfigurerDTO{
		Type: AWSServiceAccountType,
	}
}

func (arc *AWSAssumeRoleConfigurer) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	roleArn, err := GetInterfaceValue[string](fmap, "roleARN")
	if err != nil {
		return fmt.Errorf("AWSAssumeRoleConfigurer: %s", err.Error())
	}
	arc.roleARN = roleArn

	configurerAny, err := GetInterfaceValue[any](fmap, "configurer")
	if err != nil {
		return fmt.Errorf("AWSAssumeRoleConfigurer: %s", err.Error())
	}

	if configurerAny != nil {
		// If configurer is a DTO convert to AWSConfgirurer
		if configurerDTO, ok := configurerAny.(AWSConfigurerDTO); ok {
			arc.configurer = configurerDTO.ToAWSConfigurer()
		} else {
			awsConfigurer, err := AWSConfigurerFromInterface(configurerAny)
			if err != nil {
				return fmt.Errorf("AWSAssumeRoleConfigurer: %s", err.Error())
			}
			arc.configurer = awsConfigurer
		}
	}
	return nil
}

func (arc *AWSAssumeRoleConfigurer) ToDTO() AWSConfigurerDTO {
	values := make(map[string]any, 2)
	values["roleARN"] = arc.roleARN
	var configurer AWSConfigurerDTO
	if arc.configurer != nil {
		configurer = arc.configurer.ToDTO()
	}
	values["configurer"] = configurer
	return AWSConfigurerDTO{
		Type:   AWSAssumeRoleType,
		Values: values,
	}
}
