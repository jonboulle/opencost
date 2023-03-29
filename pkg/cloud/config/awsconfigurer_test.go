package config

import "testing"

func TestAWSConfigurerDTO_Sanitize(t *testing.T) {

	testCases := map[string]struct {
		input    AWSConfigurer
		expected AWSConfigurer
	}{
		"Access Key": {
			input: &AWSAccessKey{
				ID:     "ID",
				Secret: "Secret",
			},
			expected: &AWSAccessKey{
				ID:     "ID",
				Secret: Redacted,
			},
		},
		"Service Account": {
			input:    &AWSServiceAccount{},
			expected: &AWSServiceAccount{},
		},
		"Master Payer Access Key": {
			input: &AWSAssumeRoleConfigurer{
				configurer: &AWSAccessKey{
					ID:     "ID",
					Secret: "Secret",
				},
				roleARN: "role arn",
			},
			expected: &AWSAssumeRoleConfigurer{
				configurer: &AWSAccessKey{
					ID:     "ID",
					Secret: Redacted,
				},
				roleARN: "role arn",
			},
		},
		"Master Payer Service Account": {
			input: &AWSAssumeRoleConfigurer{
				configurer: &AWSServiceAccount{},
				roleARN:    "role arn",
			},
			expected: &AWSAssumeRoleConfigurer{
				configurer: &AWSServiceAccount{},
				roleARN:    "role arn",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Convert to DTO for sanitization
			dto := tc.input.ToDTO()

			// Sanitize configurer
			dto.Sanitize()

			// Convert back for equality check
			sanitizedConfigurer := dto.ToAWSConfigurer()
			if !tc.expected.Equals(sanitizedConfigurer) {
				t.Error("configurer was not as expected after Sanitization")
			}

		})
	}
}
