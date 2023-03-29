package config

import (
	"fmt"
	"github.com/opencost/opencost/pkg/util/json"
)

type BigQueryConfigurationDTO struct {
	ProjectID  string           `json:"projectID"`
	Dataset    string           `json:"dataset"`
	Table      string           `json:"table"`
	Configurer GCPConfigurerDTO `json:"configurer"`
}

// Sanitize redacts sensitive information in configuration
func (bqc *BigQueryConfigurationDTO) Sanitize() {
	bqc.Configurer.Sanitize()
}

func (bqc *BigQueryConfigurationDTO) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = bqc.FromInterface(f)
	if err != nil {
		return err
	}

	return nil
}

func (bqcd *BigQueryConfigurationDTO) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	projectID, err := GetInterfaceValue[string](fmap, "projectID")
	if err != nil {
		return fmt.Errorf("BigQueryConfigurationDTO: FromInterface: %s", err.Error())
	}
	bqcd.ProjectID = projectID

	dataset, err := GetInterfaceValue[string](fmap, "dataset")
	if err != nil {
		return fmt.Errorf("BigQueryConfigurationDTO: FromInterface: %s", err.Error())
	}
	bqcd.Dataset = dataset

	table, err := GetInterfaceValue[string](fmap, "table")
	if err != nil {
		return fmt.Errorf("BigQueryConfigurationDTO: FromInterface: %s", err.Error())
	}
	bqcd.Table = table

	if confAny, ok := fmap["configurer"]; ok {
		var configurerDTO GCPConfigurerDTO
		err = configurerDTO.FromInterface(confAny)
		if err != nil {
			return fmt.Errorf("AthenaConfigurationDTO: FromInterface: %s", err.Error())
		}
		bqcd.Configurer = configurerDTO
	}

	return nil
}

func (bqc *BigQueryConfiguration) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	projectID, err := GetInterfaceValue[string](fmap, "projectID")
	if err != nil {
		return fmt.Errorf("BigQueryConfiguration: FromInterface: %s", err.Error())
	}
	bqc.ProjectID = projectID

	dataset, err := GetInterfaceValue[string](fmap, "dataset")
	if err != nil {
		return fmt.Errorf("BigQueryConfiguration: FromInterface: %s", err.Error())
	}
	bqc.Dataset = dataset

	table, err := GetInterfaceValue[string](fmap, "table")
	if err != nil {
		return fmt.Errorf("BigQueryConfiguration: FromInterface: %s", err.Error())
	}
	bqc.Table = table

	configurerAny, err := GetInterfaceValue[any](fmap, "configurer")
	if err != nil {
		return fmt.Errorf("BigQueryConfiguration: FromInterface: %s", err.Error())
	}
	gcpConfigurer, err := GCPConfigurerFromInterface(configurerAny)
	if err != nil {
		return fmt.Errorf("BigQueryConfiguration: FromInterface: %s", err.Error())
	}
	bqc.Configurer = gcpConfigurer

	return nil
}

func (bqc *BigQueryConfiguration) ToConfigDTO() ConfigDTO {
	return bqc.ToBigQueryConfigurationDTO()
}

func (bqc *BigQueryConfiguration) ToBigQueryConfigurationDTO() *BigQueryConfigurationDTO {
	var configurer GCPConfigurerDTO
	if bqc.Configurer != nil {
		configurer = bqc.Configurer.ToDTO()
	}
	return &BigQueryConfigurationDTO{
		ProjectID:  bqc.ProjectID,
		Dataset:    bqc.Dataset,
		Table:      bqc.Table,
		Configurer: configurer,
	}
}

// UnmarshalJSON assumes data is save as an BigQueryConfigurationDTO
func (bqc *BigQueryConfiguration) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	dto := &BigQueryConfigurationDTO{}
	err = dto.FromInterface(f)
	if err != nil {
		return err
	}

	// set values from dto
	bqc.setValuesFromDTO(dto)
	return nil
}

func (bqc *BigQueryConfiguration) setValuesFromDTO(bqcd *BigQueryConfigurationDTO) {
	bqc.ProjectID = bqcd.ProjectID
	bqc.Table = bqcd.Table
	bqc.Dataset = bqcd.Dataset
	bqc.Configurer = bqcd.Configurer.ToGCPConfigurer()
}

// MarshalJSON converts to BigQueryConfigurationDTO and sanitize then Marshals as JSON
func (bqc *BigQueryConfiguration) MarshalJSON() ([]byte, error) {
	dto := bqc.ToBigQueryConfigurationDTO()
	dto.Sanitize()
	return json.Marshal(dto)
}
