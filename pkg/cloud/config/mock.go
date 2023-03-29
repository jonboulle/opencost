package config

import (
	"fmt"
)

type MockConfig struct {
}

func (mc *MockConfig) Validate() error {
	return nil
}

func (mc *MockConfig) Equals(config Config) bool {
	_, ok := config.(*MockConfig)
	return ok
}

// MockKeyedConfig implements KeyedConfig it only requires a key to be valid, there is an additional property allowing
// MockKeyedConfig with the same key to not be equal
type MockKeyedConfig struct {
	key      string
	property string
	valid    bool
}

func NewMockKeyedConfig(key, property string, valid bool) KeyedConfig {
	return &MockKeyedConfig{
		key:      key,
		property: property,
		valid:    valid,
	}
}

func (mkc *MockKeyedConfig) Validate() error {
	if !mkc.valid {
		return fmt.Errorf("MockKeyedConfig: set to invalid")
	}
	if mkc.key == "" {
		return fmt.Errorf("MockKeyedConfig: missing key")
	}
	return nil
}

func (mkc *MockKeyedConfig) Equals(config Config) bool {
	that, ok := config.(*MockKeyedConfig)
	if !ok {
		return false
	}

	if mkc.key != that.key {
		return false
	}

	if mkc.property != that.property {
		return false
	}

	if mkc.valid != that.valid {
		return false
	}

	return true
}

func (mkc *MockKeyedConfig) Key() string {
	return mkc.key
}
