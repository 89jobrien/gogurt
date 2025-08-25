package golang

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

func Split(content string) []string {
	if strings.TrimSpace(content) == "" {
		return []string{}
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return []string{content}
	}

	var chunks []string

	if file.Name != nil {
		chunks = append(chunks, "package "+file.Name.Name)
	}

	for _, decl := range file.Decls {
		chunks = append(chunks, printNode(fset, decl))
	}

	return chunks
}

func printNode(fset *token.FileSet, node ast.Node) string {
	var buf bytes.Buffer
	cfg := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
	if err := cfg.Fprint(&buf, fset, node); err != nil {
		return ""
	}
	return strings.TrimSpace(buf.String())
}
