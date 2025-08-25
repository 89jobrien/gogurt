package python

import "regexp"

// parses a .py file and splits it by classes and functions
func Split(content string) []string {
	re := regexp.MustCompile(`(?m)^(def |class )`)
	matches := re.FindAllStringIndex(content, -1)

	if len(matches) == 0 {
		return []string{content}
	}

	var chunks []string
	start := 0
	for _, match := range matches {
		end := match[0]
		if start < end {
			chunks = append(chunks, content[start:end])
		}
		start = end
	}
	// Add the last chunk.
	chunks = append(chunks, content[start:])

	return chunks
}