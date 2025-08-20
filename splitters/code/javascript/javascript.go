package javascript

import "regexp"

// parses a .js or .ts file and splits it by functions and classes.
func Split(content string) []string {
	re := regexp.MustCompile(`(?m)^(\s*)(export |default |async )*(function|class|const|let|var)\s`)
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
	chunks = append(chunks, content[start:])

	return chunks
}