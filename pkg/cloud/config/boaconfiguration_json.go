package config

import (
	"encoding/json"
	"fmt"
)

type BoaConfigurationDTO struct {
	Region     string               `json:"region"`
	Account    string               `json:"account"`
	Configurer AlibabaConfigurerDTO `json:"configurer"`
}

// Sanitize redacts sensitive information in configuration
func (bc *BoaConfigurationDTO) Sanitize() {
	bc.Configurer.Sanitize()
}

func (bc *BoaConfigurationDTO) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = bc.FromInterface(f)
	if err != nil {
		return err
	}

	return nil
}

func (bc *BoaConfigurationDTO) ToBoaConfiguration() *BoaConfiguration {
	return &BoaConfiguration{
		Region:     bc.Region,
		Account:    bc.Account,
		Configurer: bc.Configurer.ToAlibabaConfigurer(),
	}
}

func (bc *BoaConfigurationDTO) setBoaConfiguration(boaConfig *BoaConfiguration) {
	boaConfig.Region = bc.Region
	boaConfig.Account = bc.Account
	boaConfig.Configurer = bc.Configurer.ToAlibabaConfigurer()
}

func (bc *BoaConfigurationDTO) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	region, err := GetInterfaceValue[string](fmap, "region")
	if err != nil {
		return fmt.Errorf("BoaConfigurationDTO: FromInterface: %s", err.Error())
	}
	bc.Region = region

	account, err := GetInterfaceValue[string](fmap, "account")
	if err != nil {
		return fmt.Errorf("BoaConfigurationDTO: FromInterface: %s", err.Error())
	}
	bc.Account = account

	if confAny, ok := fmap["configurer"]; ok {
		var configurerDTO AlibabaConfigurerDTO
		err = configurerDTO.FromInterface(confAny)
		if err != nil {
			return fmt.Errorf("BoaConfigurationDTO: FromInterface: %s", err.Error())
		}
		bc.Configurer = configurerDTO
	}

	return nil
}

func (bc *BoaConfiguration) ToConfigDTO() ConfigDTO {
	return bc.ToBoaConfigurationDTO()
}

func (bc *BoaConfiguration) ToBoaConfigurationDTO() *BoaConfigurationDTO {
	var configurer AlibabaConfigurerDTO
	if bc.Configurer != nil {
		configurer = bc.Configurer.ToDTO()
	}
	return &BoaConfigurationDTO{
		Region:     bc.Region,
		Account:    bc.Account,
		Configurer: configurer,
	}
}

// UnmarshalJSON assumes data is save as an BoaConfigurationDTO
func (bc *BoaConfiguration) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	dto := &BoaConfigurationDTO{}
	err = dto.FromInterface(f)
	if err != nil {
		return err
	}

	// set values from dto
	dto.setBoaConfiguration(bc)
	return nil
}

// MarshalJSON converts to BoaConfigurationDTO and sanitize then Marshals as JSON
func (bc *BoaConfiguration) MarshalJSON() ([]byte, error) {
	dto := bc.ToBoaConfigurationDTO()
	dto.Sanitize()
	return json.Marshal(dto)
}
