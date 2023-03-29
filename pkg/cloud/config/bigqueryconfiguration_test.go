package config

import (
	"fmt"
	"github.com/opencost/opencost/pkg/log"
	"github.com/opencost/opencost/pkg/util/json"
	"testing"
)

func TestBigQueryConfiguration_Validate(t *testing.T) {
	testCases := map[string]struct {
		config   BigQueryConfiguration
		expected error
	}{
		"valid config GCP key": {
			config: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			expected: nil,
		},
		"valid config WorkloadIdentity": {
			config: BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: &WorkloadIdentity{},
			},
			expected: nil,
		},
		"access key invalid": {
			config: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: nil,
				},
			},
			expected: fmt.Errorf("BigQueryConfig: issue with GCP Configurer: GCPKey: missing Key"),
		},
		"missing configurer": {
			config: BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: nil,
			},
			expected: fmt.Errorf("BigQueryConfig: missing configurer"),
		},
		"missing projectID": {
			config: BigQueryConfiguration{
				ProjectID: "",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			expected: fmt.Errorf("BigQueryConfig: missing ProjectID"),
		},
		"missing dataset": {
			config: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			expected: fmt.Errorf("BigQueryConfig: missing Dataset"),
		},
		"missing table": {
			config: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			expected: fmt.Errorf("BigQueryConfig: missing Table"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := testCase.config.Validate()
			actualString := "nil"
			if actual != nil {
				actualString = actual.Error()
			}
			expectedString := "nil"
			if testCase.expected != nil {
				expectedString = testCase.expected.Error()
			}
			if actualString != expectedString {
				t.Errorf("errors do not match: Actual: '%s', Expected: '%s", actualString, expectedString)
			}
		})
	}
}

func TestBigQueryConfiguration_Equals(t *testing.T) {
	testCases := map[string]struct {
		left     BigQueryConfiguration
		right    Config
		expected bool
	}{
		"matching config": {
			left: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			right: &BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			expected: true,
		},
		"different configurer": {
			left: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			right: &BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: &WorkloadIdentity{},
			},
			expected: false,
		},
		"missing both configurer": {
			left: BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: nil,
			},
			right: &BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: nil,
			},
			expected: true,
		},
		"missing left configurer": {
			left: BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: nil,
			},
			right: &BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: &WorkloadIdentity{},
			},
			expected: false,
		},
		"missing right configurer": {
			left: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			right: &BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: nil,
			},
			expected: false,
		},
		"different projectID": {
			left: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			right: &BigQueryConfiguration{
				ProjectID: "projectID2",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			expected: false,
		},
		"different dataset": {
			left: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			right: &BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset2",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			expected: false,
		},
		"different table": {
			left: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			right: &BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table2",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			expected: false,
		},
		"different config": {
			left: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
			right: &GCPKey{

				key: map[string]string{
					"key":  "key",
					"key1": "key2",
				},
			},
			expected: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := testCase.left.Equals(testCase.right)
			if actual != testCase.expected {
				t.Errorf("incorrect result: Actual: '%t', Expected: '%t", actual, testCase.expected)
			}
		})
	}
}

func TestBigQueryConfiguration_JSON(t *testing.T) {
	testCases := map[string]struct {
		config BigQueryConfiguration
	}{
		"Empty Config": {
			config: BigQueryConfiguration{},
		},
		"Nil Configurer": {
			config: BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: nil,
			},
		},
		"GCPKeyConfigurer": {
			config: BigQueryConfiguration{
				ProjectID: "projectID",
				Dataset:   "dataset",
				Table:     "table",
				Configurer: &GCPKey{
					key: map[string]string{
						"key":  "key",
						"key1": "key2",
					},
				},
			},
		},
		"WorkLoadIdentityConfigurer": {
			config: BigQueryConfiguration{
				ProjectID:  "projectID",
				Dataset:    "dataset",
				Table:      "table",
				Configurer: &WorkloadIdentity{},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			// test dto conversion
			dto := testCase.config.ToBigQueryConfigurationDTO()
			dtoConfig := &BigQueryConfiguration{}
			dtoConfig.setValuesFromDTO(dto)
			if !testCase.config.Equals(dtoConfig) {
				t.Error("config is not equal not equal after dto conversion")
			}

			// test JSON Marshalling
			marshalledDTO, err := json.Marshal(dto)
			if err != nil {
				t.Errorf("failed to marshal configuration: %s", err.Error())
			}
			log.Info(string(marshalledDTO))
			unmarshalledDTO := &BigQueryConfigurationDTO{}
			err = json.Unmarshal(marshalledDTO, unmarshalledDTO)
			if err != nil {
				t.Errorf("failed to unmarshal configuration: %s", err.Error())
			}
			unmarshalledConfig := &BigQueryConfiguration{}
			unmarshalledConfig.setValuesFromDTO(unmarshalledDTO)
			if !testCase.config.Equals(unmarshalledConfig) {
				t.Error("config does not equal unmarshalled config")
			}
		})
	}
}
