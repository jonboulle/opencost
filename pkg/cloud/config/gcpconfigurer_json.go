package config

import (
	"fmt"
	"github.com/opencost/opencost/pkg/log"
)

const (
	GCPAccessKey        string = "GCPAccessKey"
	GCPWorkLoadIdentity        = "GCPWorkLoadIdentity"
)

type GCPConfigurerDTO struct {
	Type   string         `json:"type"`
	Values map[string]any `json:"values,omitempty"`
}

// Sanitize redacts sensitive information from the Values
func (a *GCPConfigurerDTO) Sanitize() {
	switch a.Type {
	case GCPAccessKey:
		if value, ok := a.Values["key"]; ok && value != "" {
			a.Values["key"] = Redacted
		}
	}
}

func (a *GCPConfigurerDTO) ToGCPConfigurer() GCPConfigurer {
	var configurer GCPConfigurer
	switch a.Type {
	case GCPAccessKey:
		configurer = &GCPKey{}
	case GCPWorkLoadIdentity:
		configurer = &WorkloadIdentity{}
	default:
		log.Errorf("GCPConfigurerDTO: ToGCPConfigurer: invalid type: %s", a.Type)
		return nil
	}
	err := configurer.FromInterface(a.Values)
	if err != nil {
		log.Errorf("GCPConfigurerDTO: ToGCPConfigurer: failed to convert to GCPConfigurer: %s", err.Error())
	}
	return configurer
}

func (a *GCPConfigurerDTO) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	typeOf, err := GetInterfaceValue[string](fmap, "type")
	if err != nil {
		return fmt.Errorf("GCPConfigurerDTO: %s", err.Error())
	}
	a.Type = typeOf

	_, ok := fmap["values"]
	if ok {
		values, err := GetInterfaceValue[map[string]any](fmap, "values")
		if err != nil {
			return fmt.Errorf("GCPConfigurerDTO: %s", err.Error())
		}
		a.Values = values
	}

	return nil
}

func GCPConfigurerFromInterface(f any) (GCPConfigurer, error) {
	if f == nil {
		return nil, nil
	}
	fmap, ok := f.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("GCPConfigurerFromInterface: could not case interface as map")
	}
	if len(fmap) == 0 {
		return &WorkloadIdentity{}, nil
	}
	if _, ok = fmap["key"]; ok {
		gcpConfigurer := &GCPKey{}
		err := gcpConfigurer.FromInterface(fmap)
		if err != nil {
			return nil, err
		}
		return gcpConfigurer, nil
	}
	if _, ok = fmap["type"]; ok {
		dto := &GCPConfigurerDTO{}
		err := dto.FromInterface(fmap)
		if err != nil {
			return nil, err
		}
		return dto.ToGCPConfigurer(), nil
	}

	return nil, fmt.Errorf("GCPConfigurerFromInterface: interface is not a valid configurer format")
}

func (gk *GCPKey) FromInterface(f interface{}) error {
	fmap := f.(map[string]interface{})

	keyInterface, err := GetInterfaceValue[map[string]any](fmap, "key")
	if err != nil {
		return fmt.Errorf("GCPKey: %s", err.Error())
	}
	key := make(map[string]string, len(keyInterface))
	for k, v := range keyInterface {
		value, ok := v.(string)
		if !ok {
			fmt.Errorf("GCPKey: non-string in key %s", err.Error())
		}
		key[k] = value
	}
	gk.key = key

	return nil
}

func (gk *GCPKey) ToDTO() GCPConfigurerDTO {
	values := make(map[string]any, 1)
	keyInterface := make(map[string]any, len(gk.key))
	for k, v := range gk.key {
		keyInterface[k] = v
	}
	values["key"] = keyInterface
	return GCPConfigurerDTO{
		Type:   GCPAccessKey,
		Values: values,
	}
}

func (wi *WorkloadIdentity) FromInterface(f interface{}) error {
	return nil
}

func (wi *WorkloadIdentity) ToDTO() GCPConfigurerDTO {
	return GCPConfigurerDTO{
		Type: GCPWorkLoadIdentity,
	}
}
