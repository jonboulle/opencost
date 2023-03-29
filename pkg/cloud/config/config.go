package config

// Config allows for nested configurations which encapsulate their functionality to be validated and compared easily
type Config interface {
	Validate() error
	Equals(Config) bool
}

// KeyedConfig is a top level Config which uses its public values as a unique identifier allowing duplicates to be identified
type KeyedConfig interface {
	Config
	Key() string
}

type ConfigDTO interface {
	Sanitize()
}
