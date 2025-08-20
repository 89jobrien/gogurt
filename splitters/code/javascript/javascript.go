package javascript

import "regexp"

// Split parses a .js or .ts file and splits it by top-level
// functions, classes, and variable declarations.
// Each chunk starts at a match for: "function", "class", "const", "let", "var",
// possibly preceded by "export", "default", and "async" modifiers.
// Blank lines before a match stay with the previous chunk.
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
