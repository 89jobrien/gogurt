package prompts

import (
	"testing"
)

func TestPromptTemplate_Format(t *testing.T) {
	templateStr := "Context: {{.context}} Question: {{.question}}"
	pt, err := NewPromptTemplate(templateStr)
	if err != nil {
		t.Fatalf("failed to create prompt template: %v", err)
	}

	testCases := []struct {
		name     string
		data     map[string]string
		expected string
	}{
		{
			name: "Simple case",
			data: map[string]string{
				"context":  "Gogurt is a Go framework.",
				"question": "What is Gogurt?",
			},
			expected: "Context: Gogurt is a Go framework. Question: What is Gogurt?",
		},
		{
			name: "Empty context",
			data: map[string]string{
				"context":  "",
				"question": "What is Gogurt?",
			},
			expected: "Context:  Question: What is Gogurt?",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := pt.Format(tc.data)
			if err != nil {
				t.Errorf("Format() returned an error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Format() = %q, want %q", result, tc.expected)
			}
		})
	}
}