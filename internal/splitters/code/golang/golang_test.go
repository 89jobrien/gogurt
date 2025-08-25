package golang

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "Standard Go file with multiple functions and types",
			content: `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}

type MyStruct struct {
	Field int
}`,
			expected: []string{
				"package main",
				"import \"fmt\"",
				"func main() {\n    fmt.Println(\"Hello, World!\")\n}",
				"type MyStruct struct {\n    Field int\n}",
			},
		},
		{
			name:    "File with only one function",
			content: "package main\n\nfunc hello() {}",
			expected: []string{
				"package main",
				"func hello() {}",
			},
		},
		{
			name: "File with multiple imports",
			content: `package main

import (
	"fmt"
	"os"
)

// comment above
func main() {
	fmt.Println("hi")
}

type Foo struct{}`,
			expected: []string{
				"package main",
				"import (\n    \"fmt\"\n    \"os\"\n)",
				"// comment above\nfunc main() {\n    fmt.Println(\"hi\")\n}",
				"type Foo struct{}",
			},
		},
		{
			name:    "File with tabs (ensure normalization to spaces)",
			content: "package main\n\nfunc main() {\n\tfmt.Println(\"hello with tab\")\n}\n",
			expected: []string{
				"package main",
				"func main() {\n    fmt.Println(\"hello with tab\")\n}",
			},
		},
		{
			name: "File with multi-const block",
			content: `package main

const (
	A = 1
	B = 2
	C = 3
)`,
			expected: []string{
				"package main",
				"const (\n    A   = 1\n    B   = 2\n    C   = 3\n)",
			},
		},
		{
			name: "File with multi-var block",
			content: `package main

var (
	x = 10
	y = 20
)`,
			expected: []string{
				"package main",
				"var (\n    x   = 10\n    y   = 20\n)",
			},
		},
		{
			name:     "Empty file",
			content:  "",
			expected: []string{},
		},
		{
			name:    "File with syntax errors",
			content: "package main\n\nfunc hello() {",
			expected: []string{
				"package main\n\nfunc hello() {",
			},
		},
		{
			name: "File with comments",
			content: `package main

// This is a comment
func hello() {}`,
			expected: []string{
				"package main",
				"// This is a comment\nfunc hello() {}",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Split(tc.content)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Split() =\n%#v,\nwant\n%#v", result, tc.expected)
			}
		})
	}
}
