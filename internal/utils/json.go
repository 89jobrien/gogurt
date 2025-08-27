package utils

import "strings"

// extractJSONArray finds and returns the first valid JSON array from a string.
func ExtractJSONArray(s string) string {
	start := strings.Index(s, "[")
	if start == -1 {
		return ""
	}

	balance := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
			case '[':
				balance++
			case ']':
				balance--
				if balance == 0 {
					return s[start : i+1]
			}
		}
	}

	return ""
}
