package config

import (
	"encoding/json"
	"fmt"
)

type AthenaConfigurationDTO struct {
	Bucket     string           `json:"bucket"`
	Region     string           `json:"region"`
	Database   string           `json:"database"`
	Table      string           `json:"table"`
	Workgroup  string           `json:"workgroup"`
	Account    string           `json:"account"`
	Configurer AWSConfigurerDTO `json:"configurer"`
}

// Sanitize redacts sensitive information in configuration
func (a *AthenaConfigurationDTO) Sanitize() {
	a.Configurer.Sanitize()
}

func (a *AthenaConfigurationDTO) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = a.FromInterface(f)
	if err != nil {
		return err
	}

	return nil
}

func (a *AthenaConfigurationDTO) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	bucket, err := GetInterfaceValue[string](fmap, "bucket")
	if err != nil {
		return fmt.Errorf("AthenaConfigurationDTO: FromInterface: %s", err.Error())
	}
	a.Bucket = bucket

	region, err := GetInterfaceValue[string](fmap, "region")
	if err != nil {
		return fmt.Errorf("AthenaConfigurationDTO: FromInterface: %s", err.Error())
	}
	a.Region = region

	database, err := GetInterfaceValue[string](fmap, "database")
	if err != nil {
		return fmt.Errorf("AthenaConfigurationDTO: FromInterface: %s", err.Error())
	}
	a.Database = database

	table, err := GetInterfaceValue[string](fmap, "table")
	if err != nil {
		return fmt.Errorf("AthenaConfigurationDTO: FromInterface: %s", err.Error())
	}
	a.Table = table

	workgroup, err := GetInterfaceValue[string](fmap, "workgroup")
	if err != nil {
		return fmt.Errorf("AthenaConfigurationDTO: FromInterface: %s", err.Error())
	}
	a.Workgroup = workgroup

	account, err := GetInterfaceValue[string](fmap, "account")
	if err != nil {
		return fmt.Errorf("AthenaConfigurationDTO: FromInterface: %s", err.Error())
	}
	a.Account = account

	if confAny, ok := fmap["configurer"]; ok {
		var configurerDTO AWSConfigurerDTO
		err = configurerDTO.FromInterface(confAny)
		if err != nil {
			return fmt.Errorf("AthenaConfigurationDTO: FromInterface: %s", err.Error())
		}
		a.Configurer = configurerDTO
	}

	return nil
}

func (ac *AthenaConfiguration) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	bucket, err := GetInterfaceValue[string](fmap, "bucket")
	if err != nil {
		return fmt.Errorf("AthenaConfiguration: FromInterface: %s", err.Error())
	}
	ac.Bucket = bucket

	region, err := GetInterfaceValue[string](fmap, "region")
	if err != nil {
		return fmt.Errorf("AthenaConfiguration: FromInterface: %s", err.Error())
	}
	ac.Region = region

	database, err := GetInterfaceValue[string](fmap, "database")
	if err != nil {
		return fmt.Errorf("AthenaConfiguration: FromInterface: %s", err.Error())
	}
	ac.Database = database

	table, err := GetInterfaceValue[string](fmap, "table")
	if err != nil {
		return fmt.Errorf("AthenaConfiguration: FromInterface: %s", err.Error())
	}
	ac.Table = table

	workgroup, err := GetInterfaceValue[string](fmap, "workgroup")
	if err != nil {
		return fmt.Errorf("AthenaConfiguration: FromInterface: %s", err.Error())
	}
	ac.Workgroup = workgroup

	account, err := GetInterfaceValue[string](fmap, "account")
	if err != nil {
		return fmt.Errorf("AthenaConfiguration: FromInterface: %s", err.Error())
	}
	ac.Account = account

	configurerAny, err := GetInterfaceValue[any](fmap, "configurer")
	if err != nil {
		return fmt.Errorf("AthenaConfiguration: FromInterface: %s", err.Error())
	}
	awsConfigurer, err := AWSConfigurerFromInterface(configurerAny)
	if err != nil {
		return fmt.Errorf("AthenaConfiguration: FromInterface: %s", err.Error())
	}
	ac.Configurer = awsConfigurer

	return nil
}

func (ac *AthenaConfiguration) ToConfigDTO() ConfigDTO {
	return ac.ToAthenaConfigurationDTO()
}

func (ac *AthenaConfiguration) ToAthenaConfigurationDTO() *AthenaConfigurationDTO {
	var configurer AWSConfigurerDTO
	if ac.Configurer != nil {
		configurer = ac.Configurer.ToDTO()
	}
	return &AthenaConfigurationDTO{
		Bucket:     ac.Bucket,
		Region:     ac.Region,
		Database:   ac.Database,
		Table:      ac.Table,
		Workgroup:  ac.Workgroup,
		Account:    ac.Account,
		Configurer: configurer,
	}
}

// UnmarshalJSON assumes data is save as an AthenaConfigurationDTO
func (ac *AthenaConfiguration) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	dto := &AthenaConfigurationDTO{}
	err = dto.FromInterface(f)
	if err != nil {
		return err
	}

	// set values from dto
	ac.setValuesFromDTO(dto)
	return nil
}

func (ac *AthenaConfiguration) setValuesFromDTO(a *AthenaConfigurationDTO) {
	ac.Bucket = a.Bucket
	ac.Region = a.Region
	ac.Database = a.Database
	ac.Table = a.Table
	ac.Workgroup = a.Workgroup
	ac.Account = a.Account
	ac.Configurer = a.Configurer.ToAWSConfigurer()
}

// MarshalJSON converts to AthenaConfigurationDTO and sanitize then Marshals as JSON
func (ac *AthenaConfiguration) MarshalJSON() ([]byte, error) {
	dto := ac.ToAthenaConfigurationDTO()
	dto.Sanitize()
	return json.Marshal(dto)
}
