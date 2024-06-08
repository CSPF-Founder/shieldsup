package iputils

import "testing"

func TestGetIPCountIfRange(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Valid CIDR",
			input:    "192.168.56.1/24",
			expected: 256,
		},
		{
			name:     "Invalid CIDR",
			input:    "192.168.56.1/42",
			expected: 0,
		},
		{
			name:     "CIDR /32",
			input:    "192.168.56.1/32",
			expected: 1,
		},
		{
			name:     "CIDR /16",
			input:    "192.168.0.1/16",
			expected: 65536,
		},
		{
			name:     "CIDR /25",
			input:    "192.168.0.1/25",
			expected: 128,
		},
		{
			name:     "CIDR /31",
			input:    "192.168.56.130/31",
			expected: 2,
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function being tested
			result, _ := GetIPCountIfRange(tt.input)

			// Check the result
			if result != tt.expected {
				t.Errorf("Expected: %d, Got: %d", tt.expected, result)
			}
		})
	}
}
