package simple

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     []float32
		expected float64
	}{
		{
			name:     "Identical vectors",
			a:        []float32{1, 1, 1},
			b:        []float32{1, 1, 1},
			expected: 1.0,
		},
		{
			name:     "Opposite vectors",
			a:        []float32{1, 1, 1},
			b:        []float32{-1, -1, -1},
			expected: -1.0,
		},
		{
			name:     "Orthogonal vectors",
			a:        []float32{1, 0},
			b:        []float32{0, 1},
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cosineSimilarity(tc.a, tc.b)
			if math.Abs(result-tc.expected) > 1e-9 {
				t.Errorf("cosineSimilarity() = %v, want %v", result, tc.expected)
			}
		})
	}
}
