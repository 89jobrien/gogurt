package golang

import (
	"go/parser"
	"go/token"
)

// parses a .go file and splits it by top-level functions and types
func Split(content string) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return []string{content}
	}

	var chunks []string
	for _, decl := range file.Decls {
		start := fset.Position(decl.Pos()).Offset
		end := min(fset.Position(decl.End()).Offset, len(content))
		chunks = append(chunks, content[start:end])
	}
	return chunks
}