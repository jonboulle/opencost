package config

import (
	"fmt"
	"github.com/opencost/opencost/pkg/log"
	"github.com/opencost/opencost/pkg/util/json"
	"testing"
)

func TestAzureStorageConfiguration_Validate(t *testing.T) {
	testCases := map[string]struct {
		config   AzureStorageConfiguration
		expected error
	}{
		"valid config Azure AccessKey": {
			config: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: nil,
		},
		"access key invalid": {
			config: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					Account: "account",
				},
			},
			expected: fmt.Errorf("AzureAccessKey: missing access key"),
		},
		"missing configurer": {
			config: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer:     nil,
			},
			expected: fmt.Errorf("AzureStorageConfiguration: missing configurer"),
		},
		"missing subscriptionID": {
			config: AzureStorageConfiguration{
				SubscriptionID: "",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: fmt.Errorf("AzureStorageConfiguration: missing Subcription ID"),
		},
		"missing account": {
			config: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: fmt.Errorf("AzureStorageConfiguration: missing Account"),
		},
		"missing container": {
			config: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: fmt.Errorf("AzureStorageConfiguration: missing Container"),
		},
		"missing path": {
			config: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: nil,
		},
		"missing cloud": {
			config: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: nil,
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

func TestAzureStorageConfiguration_Equals(t *testing.T) {
	testCases := map[string]struct {
		left     AzureStorageConfiguration
		right    Config
		expected bool
	}{
		"matching config": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			right: &AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: true,
		},

		"missing both configurer": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer:     nil,
			},
			right: &AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer:     nil,
			},
			expected: true,
		},
		"missing left configurer": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer:     nil,
			},
			right: &AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: false,
		},
		"missing right configurer": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			right: &AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer:     nil,
			},
			expected: false,
		},
		"different subscriptionID": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			right: &AzureStorageConfiguration{
				SubscriptionID: "subscriptionID2",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: false,
		},
		"different account": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			right: &AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account2",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: false,
		},
		"different container": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			right: &AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container2",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: false,
		},
		"different path": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			right: &AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path2",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: false,
		},
		"different cloud": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			right: &AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud2",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			expected: false,
		},
		"different config": {
			left: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
			right: &AzureAccessKey{
				AccessKey: "accessKey",
				Account:   "account",
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

func TestAzureStorageConfiguration_JSON(t *testing.T) {
	testCases := map[string]struct {
		config AzureStorageConfiguration
	}{
		"Empty Config": {
			config: AzureStorageConfiguration{},
		},
		"Nil Configurer": {
			config: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer:     nil,
			},
		},
		"AccessKey Configurer": {
			config: AzureStorageConfiguration{
				SubscriptionID: "subscriptionID",
				Account:        "account",
				Container:      "container",
				Path:           "path",
				Cloud:          "cloud",
				Configurer: &AzureAccessKey{
					AccessKey: "accessKey",
					Account:   "account",
				},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			// test dto conversion
			dto := testCase.config.ToAzureStorageConfigurationDTO()
			dtoConfig := &AzureStorageConfiguration{}
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
			unmarshalledDTO := &AzureStorageConfigurationDTO{}
			err = json.Unmarshal(marshalledDTO, unmarshalledDTO)
			if err != nil {
				t.Errorf("failed to unmarshal configuration: %s", err.Error())
			}
			unmarshalledConfig := &AzureStorageConfiguration{}
			unmarshalledConfig.setValuesFromDTO(unmarshalledDTO)
			if !testCase.config.Equals(unmarshalledConfig) {
				t.Error("config does not equal unmarshalled config")
			}
		})
	}
}
