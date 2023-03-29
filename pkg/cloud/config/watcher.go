package config

import (
	"strings"

	"github.com/opencost/opencost/pkg/cloud"
	"github.com/opencost/opencost/pkg/log"
)

// ConfigSource is an Enum of the sources int value of the Source determines its priority
type ConfigSource int

const (
	UnknownSource ConfigSource = iota
	ConfigControllerSource
	MultiCloudSource
	ConfigFileSource
	HelmSource
)

func GetConfigSource(str string) ConfigSource {
	switch str {
	case "configController":
		return ConfigControllerSource
	case "configfile":
		return ConfigFileSource
	case "helm":
		return HelmSource
	case "multicloud":
		return MultiCloudSource
	default:
		return UnknownSource
	}
}

func (cs ConfigSource) String() string {
	switch cs {
	case ConfigControllerSource:
		return "configController"
	case ConfigFileSource:
		return "configfile"
	case HelmSource:
		return "helm"
	case MultiCloudSource:
		return "multicloud"
	case UnknownSource:
		return "unknown"
	default:
		return "unknown"
	}
}

func ExtractConfigFromProviders(provider cloud.Provider) *cloud.ProviderConfig {
	if provider == nil {
		log.Errorf("cannot extract config from nil provider")
		return nil
	}
	switch p := provider.(type) {
	case *cloud.CSVProvider:
		return ExtractConfigFromProviders(p.CustomProvider)
	case *cloud.CustomProvider:
		return p.Config
	case *cloud.GCP:
		return p.Config
	case *cloud.AWS:
		return p.Config
	case *cloud.Azure:
		return p.Config
	case *cloud.Alibaba:
		return p.Config
	default:
		log.Errorf("failed to extract config from provider")
		return nil
	}
}

// ConvertAwsAthenaInfoToConfig takes a legacy config and generates a Config based on the presence of properties to match
// legacy behavior
func ConvertAwsAthenaInfoToConfig(aai cloud.AwsAthenaInfo) KeyedConfig {
	if aai.IsEmpty() {
		return nil
	}

	var configurer AWSConfigurer
	if aai.ServiceKeyName == "" && aai.ServiceKeySecret == "" {
		configurer = &AWSServiceAccount{}
	} else {
		configurer = &AWSAccessKey{
			ID:     aai.ServiceKeyName,
			Secret: aai.ServiceKeySecret,
		}
	}

	// Wrap configurer with AWSAssumeRoleConfigurer if MasterPayerArn is set
	if aai.MasterPayerARN != "" {
		configurer = &AWSAssumeRoleConfigurer{
			configurer: configurer,
			roleARN:    aai.MasterPayerARN,
		}
	}

	var config KeyedConfig
	if aai.AthenaTable != "" || aai.AthenaDatabase != "" {
		config = &AthenaConfiguration{
			Bucket:     aai.AthenaBucketName,
			Region:     aai.AthenaRegion,
			Database:   aai.AthenaDatabase,
			Table:      aai.AthenaTable,
			Workgroup:  aai.AthenaWorkgroup,
			Account:    aai.AccountID,
			Configurer: configurer,
		}
	} else {
		config = &S3Configuration{
			Bucket:        aai.AthenaBucketName,
			Region:        aai.AthenaRegion,
			Account:       aai.AccountID,
			AWSConfigurer: configurer,
		}
	}

	return config
}

func ConvertAlibabaInfoToConfig(acc cloud.AlibabaInfo) KeyedConfig {
	if acc.IsEmpty() {
		return nil
	}
	var configurer AlibabaConfigurer

	configurer = &AlibabaAccessKey{
		AccessKeyId:     acc.AlibabaServiceKeyName,
		AccessKeySecret: acc.AlibabaServiceKeySecret,
	}

	return &BoaConfiguration{
		Account:    acc.AlibabaAccountID,
		Region:     acc.AlibabaClusterRegion,
		Configurer: configurer,
	}
}
func ConvertBigQueryConfigToConfig(bqc cloud.BigQueryConfig) KeyedConfig {
	if bqc.IsEmpty() {
		return nil
	}

	BillingDataDataset := strings.Split(bqc.BillingDataDataset, ".")
	dataset := BillingDataDataset[0]
	var table string
	if len(BillingDataDataset) > 1 {
		table = BillingDataDataset[1]
	}

	config := &BigQueryConfiguration{
		ProjectID:  bqc.ProjectID,
		Dataset:    dataset,
		Table:      table,
		Configurer: nil,
	}

	if len(bqc.Key) != 0 {
		config.Configurer = &GCPKey{
			key: bqc.Key,
		}
	}

	return config
}

func ConvertAzureStorageConfigToConfig(asc cloud.AzureStorageConfig) KeyedConfig {
	if asc.IsEmpty() {
		return nil
	}

	var configurer AzureConfigurer
	configurer = &AzureAccessKey{
		AccessKey: asc.AccessKey,
		Account:   asc.AccountName,
	}

	return &AzureStorageConfiguration{
		SubscriptionID: asc.SubscriptionId,
		Account:        asc.AccountName,
		Container:      asc.ContainerName,
		Path:           asc.ContainerPath,
		Cloud:          asc.AzureCloud,
		Configurer:     configurer,
	}
}
