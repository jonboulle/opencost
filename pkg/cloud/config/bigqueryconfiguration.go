package config

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
)

type BigQueryConfiguration struct {
	ProjectID  string
	Dataset    string
	Table      string
	Configurer GCPConfigurer
}

func (bqc *BigQueryConfiguration) Validate() error {

	if bqc.Configurer == nil {
		return fmt.Errorf("BigQueryConfig: missing configurer")
	}

	err := bqc.Configurer.Validate()
	if err != nil {
		return fmt.Errorf("BigQueryConfig: issue with GCP Configurer: %s", err.Error())
	}

	if bqc.ProjectID == "" {
		return fmt.Errorf("BigQueryConfig: missing ProjectID")
	}

	if bqc.Dataset == "" {
		return fmt.Errorf("BigQueryConfig: missing Dataset")
	}

	if bqc.Table == "" {
		return fmt.Errorf("BigQueryConfig: missing Table")
	}

	return nil
}

func (bqc *BigQueryConfiguration) Equals(config Config) bool {
	if config == nil {
		return false
	}
	thatConfig, ok := config.(*BigQueryConfiguration)
	if !ok {
		return false
	}

	if bqc.Configurer != nil {
		if !bqc.Configurer.Equals(thatConfig.Configurer) {
			return false
		}
	} else {
		if thatConfig.Configurer != nil {
			return false
		}
	}

	if bqc.ProjectID != thatConfig.ProjectID {
		return false
	}

	if bqc.Dataset != thatConfig.Dataset {
		return false
	}

	if bqc.Table != thatConfig.Table {
		return false
	}

	return true
}

// Key uses the Usage Project Id as the Provider Key for GCP
func (bqc *BigQueryConfiguration) Key() string {
	return fmt.Sprintf("%s/%s", bqc.ProjectID, bqc.GetBillingDataDataset())
}

func (bqc *BigQueryConfiguration) GetBillingDataDataset() string {
	return fmt.Sprintf("%s.%s", bqc.Dataset, bqc.Table)
}

func (bqc *BigQueryConfiguration) GetBigQueryClient(ctx context.Context) (*bigquery.Client, error) {
	clientOpts, err := bqc.Configurer.CreateGCPClientOption()
	if err != nil {
		return nil, err
	}
	return bigquery.NewClient(ctx, bqc.ProjectID, clientOpts)
}
