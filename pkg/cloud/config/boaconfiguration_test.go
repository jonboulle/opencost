package config

import (
	"fmt"
	"testing"
)

func TestBoaConfiguration_Validate(t *testing.T) {
	testCases := map[string]struct {
		config   BoaConfiguration
		expected error
	}{
		"valid config Azure AccessKey": {
			config: BoaConfiguration{
				Account: "Account Number",
				Region:  "Region",
				Configurer: &AlibabaAccessKey{
					AccessKeyId:     "accessKeyID",
					AccessKeySecret: "accessKeySecret",
				},
			},
			expected: nil,
		},
		"access key invalid": {
			config: BoaConfiguration{
				Account: "Account Number",
				Region:  "Region",
				Configurer: &AlibabaAccessKey{
					AccessKeySecret: "accessKeySecret",
				},
			},
			expected: fmt.Errorf("AlibabaAccessKey: missing Access key ID"),
		},
		"access secret invalid": {
			config: BoaConfiguration{
				Account: "Account Number",
				Region:  "Region",
				Configurer: &AlibabaAccessKey{
					AccessKeyId: "accessKeyId",
				},
			},
			expected: fmt.Errorf("AlibabaAccessKey: missing Access Key secret"),
		},
		"missing configurer": {
			config: BoaConfiguration{
				Account:    "Account Number",
				Region:     "Region",
				Configurer: nil,
			},
			expected: fmt.Errorf("BoaConfiguration: missing configurer"),
		},
		"missing Account": {
			config: BoaConfiguration{
				Account: "",
				Region:  "Region",
				Configurer: &AlibabaAccessKey{
					AccessKeyId:     "accessKeyID",
					AccessKeySecret: "accessKeySecret",
				},
			},
			expected: fmt.Errorf("BoaConfiguration: missing account"),
		},
		"missing Region": {
			config: BoaConfiguration{
				Account: "Account",
				Configurer: &AlibabaAccessKey{
					AccessKeyId:     "accessKeyID",
					AccessKeySecret: "accessKeySecret",
				},
			},
			expected: fmt.Errorf("BoaConfiguration: missing region"),
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
