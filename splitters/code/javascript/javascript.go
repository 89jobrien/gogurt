package javascript

import "regexp"

func Split(content string) []string {
	re := regexp.MustCompile(`(?m)^\s*(?:export\s+)?(?:default\s+)?(?:async\s+)?(function|class|const|let|var)\s`)
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
