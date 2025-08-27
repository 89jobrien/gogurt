package utils

import "testing"

func TestExtractJSONArray(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple JSON array",
			input:    `[{"key":"value"}]`,
			expected: `[{"key":"value"}]`,
		},
		{
			name:     "JSON array with surrounding text",
			input:    `Here is the JSON: [{"key":"value"}] Thanks!`,
			expected: `[{"key":"value"}]`,
		},
		{
			name:     "Nested JSON array",
			input:    `[[1, 2], [3, 4]]`,
			expected: `[[1, 2], [3, 4]]`,
		},
		{
			name:     "No JSON array",
			input:    `{"key":"value"}`,
			expected: "",
		},
		{
			name:     "Malformed JSON array",
			input:    `[{"key":"value"}`,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ExtractJSONArray(tc.input)
			if result != tc.expected {
				t.Errorf("ExtractJSONArray() = %q, want %q", result, tc.expected)
			}
		})
	}
}