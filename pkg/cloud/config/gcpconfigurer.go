package config

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
)

type GCPConfigurer interface {
	Config
	CreateGCPClientOption() (option.ClientOption, error)
	FromInterface(any) error
	ToDTO() GCPConfigurerDTO
}

type GCPKey struct {
	key map[string]string
}

func (gkc *GCPKey) Validate() error {
	if gkc.key == nil || len(gkc.key) == 0 {
		return fmt.Errorf("GCPKey: missing Key")
	}

	return nil
}

func (gkc *GCPKey) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*GCPKey)
	if !ok {
		return false
	}

	if len(gkc.key) != len(thatConfig.key) {
		return false
	}

	for k, v := range gkc.key {
		if thatConfig.key[k] != v {
			return false
		}
	}

	return true
}

func (gkc *GCPKey) CreateGCPClientOption() (option.ClientOption, error) {
	err := gkc.Validate()
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(gkc.key)
	if err != nil {
		return nil, fmt.Errorf("GCPKey: failed to marshal key: %s", err.Error())
	}
	clientOption := option.WithCredentialsJSON(b)

	// The creation of the BigQuery Client is where FAILED_CONNECTION CloudConnectionStatus is recorded for GCP
	return clientOption, nil
}

// WorkloadIdentity is a configurer which does
type WorkloadIdentity struct{}

func (wi *WorkloadIdentity) Validate() error {
	return nil
}

func (wi *WorkloadIdentity) Equals(config Config) bool {
	if config == nil {
		return false
	}
	_, ok := config.(*WorkloadIdentity)
	if !ok {
		return false
	}

	return true
}

func (wi *WorkloadIdentity) CreateGCPClientOption() (option.ClientOption, error) {
	return nil, nil
}
