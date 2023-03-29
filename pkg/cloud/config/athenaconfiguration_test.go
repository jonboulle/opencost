package config

import (
	"fmt"
	"github.com/opencost/opencost/pkg/log"
	"github.com/opencost/opencost/pkg/util/json"
	"testing"
)

func TestAthenaConfiguration_Validate(t *testing.T) {
	testCases := map[string]struct {
		config   AthenaConfiguration
		expected error
	}{
		"valid config access key": {
			config: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			expected: nil,
		},
		"valid config service account": {
			config: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: &AWSServiceAccount{},
			},
			expected: nil,
		},
		"access key invalid": {
			config: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID: "id",
				},
			},
			expected: fmt.Errorf("AthenaConfiguration: AWSAccessKey: missing Secret"),
		},
		"missing configurer": {
			config: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: nil,
			},
			expected: fmt.Errorf("AthenaConfiguration: missing configurer"),
		},
		"missing bucket": {
			config: AthenaConfiguration{
				Bucket:     "",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: &AWSServiceAccount{},
			},
			expected: fmt.Errorf("AthenaConfiguration: missing bucket"),
		},
		"missing region": {
			config: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: &AWSServiceAccount{},
			},
			expected: fmt.Errorf("AthenaConfiguration: missing region"),
		},
		"missing database": {
			config: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: &AWSServiceAccount{},
			},
			expected: fmt.Errorf("AthenaConfiguration: missing database"),
		},
		"missing table": {
			config: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: &AWSServiceAccount{},
			},
			expected: fmt.Errorf("AthenaConfiguration: missing table"),
		},
		"missing workgroup": {
			config: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "",
				Account:    "account",
				Configurer: &AWSServiceAccount{},
			},
			expected: nil,
		},
		"missing account": {
			config: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "",
				Configurer: &AWSServiceAccount{},
			},
			expected: fmt.Errorf("AthenaConfiguration: missing account"),
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

func TestAthenaConfiguration_Equals(t *testing.T) {
	testCases := map[string]struct {
		left     AthenaConfiguration
		right    Config
		expected bool
	}{
		"matching config": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			expected: true,
		},
		"different configurer": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: &AWSServiceAccount{},
			},
			expected: false,
		},
		"missing both configurer": {
			left: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: nil,
			},
			right: &AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: nil,
			},
			expected: true,
		},
		"missing left configurer": {
			left: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: nil,
			},
			right: &AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: &AWSServiceAccount{},
			},
			expected: false,
		},
		"missing right configurer": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: nil,
			},
			expected: false,
		},
		"different bucket": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AthenaConfiguration{
				Bucket:    "bucket2",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			expected: false,
		},
		"different region": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region2",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			expected: false,
		},
		"different database": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database2",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			expected: false,
		},
		"different table": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table2",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			expected: false,
		},
		"different workgroup": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup2",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			expected: false,
		},
		"different account": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account2",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			expected: false,
		},
		"different config": {
			left: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
			right: &AWSAccessKey{
				ID:     "id",
				Secret: "secret",
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

func TestAthenaConfiguration_JSON(t *testing.T) {
	testCases := map[string]struct {
		config AthenaConfiguration
	}{
		"Empty Config": {
			config: AthenaConfiguration{},
		},
		"AWSAccessKey": {
			config: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAccessKey{
					ID:     "id",
					Secret: "secret",
				},
			},
		},

		"AWSServiceAccount": {
			config: AthenaConfiguration{
				Bucket:     "bucket",
				Region:     "region",
				Database:   "database",
				Table:      "table",
				Workgroup:  "workgroup",
				Account:    "account",
				Configurer: &AWSServiceAccount{},
			},
		},
		"RoleArnAWSAccessKey": {
			config: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAssumeRoleConfigurer{
					configurer: &AWSAccessKey{
						ID:     "id",
						Secret: "secret",
					},
					roleARN: "12345",
				},
			},
		},
		"RoleArnAWSServiceAccount": {
			config: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAssumeRoleConfigurer{
					configurer: &AWSServiceAccount{},
					roleARN:    "12345",
				},
			},
		},
		"RoleArnNil": {
			config: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAssumeRoleConfigurer{
					configurer: nil,
					roleARN:    "12345",
				},
			},
		},
		"RoleArnRoleArnAWSServiceAccount": {
			config: AthenaConfiguration{
				Bucket:    "bucket",
				Region:    "region",
				Database:  "database",
				Table:     "table",
				Workgroup: "workgroup",
				Account:   "account",
				Configurer: &AWSAssumeRoleConfigurer{
					configurer: &AWSAssumeRoleConfigurer{
						roleARN:    "12345",
						configurer: &AWSServiceAccount{},
					},
					roleARN: "12345",
				},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			// test dto conversion
			dto := testCase.config.ToAthenaConfigurationDTO()
			dtoConfig := &AthenaConfiguration{}
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
			unmarshalledDTO := &AthenaConfigurationDTO{}
			err = json.Unmarshal(marshalledDTO, unmarshalledDTO)
			if err != nil {
				t.Errorf("failed to unmarshal configuration: %s", err.Error())
			}
			unmarshalledConfig := &AthenaConfiguration{}
			unmarshalledConfig.setValuesFromDTO(unmarshalledDTO)
			if !testCase.config.Equals(unmarshalledConfig) {
				t.Error("config does not equal unmarshalled config")
			}
		})
	}
}
