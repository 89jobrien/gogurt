package markdown

import (
	"gogurt/internal/types"
	"reflect"
	"testing"
)

func TestMarkdownSplitter_SplitDocuments(t *testing.T) {
	splitter := New(1024, 100)

	testCases := []struct {
		name     string
		doc      types.Document
		expected []types.Document
	}{
		{
			name: "Standard Case with Multiple Headers",
			doc: types.Document{
				PageContent: "# Header 1\n\nSome text.\n\n## Header 2\n\nMore text.",
				Metadata:    map[string]any{"source": "test.md"},
			},
			expected: []types.Document{
				{PageContent: "Header 1\n\nSome text.", Metadata: map[string]any{"source": "test.md"}},
				{PageContent: "Header 2\n\nMore text.", Metadata: map[string]any{"source": "test.md"}},
			},
		},
		{
			name: "Content Before the First Header",
			doc: types.Document{
				PageContent: "Preamble.\n\n# Header 1\n\nSome text.",
				Metadata:    map[string]any{"source": "test.md"},
			},
			expected: []types.Document{
				{PageContent: "Preamble.", Metadata: map[string]any{"source": "test.md"}},
				{PageContent: "Header 1\n\nSome text.", Metadata: map[string]any{"source": "test.md"}},
			},
		},
		{
			name: "No Headers in Content",
			doc: types.Document{
				PageContent: "Just a single block of text with no headers.",
				Metadata:    map[string]any{"source": "test.md"},
			},
			expected: []types.Document{
				{PageContent: "Just a single block of text with no headers.", Metadata: map[string]any{"source": "test.md"}},
			},
		},
		// {
		// 	name: "Various Heading Levels",
		// 	doc: types.Document{
		// 		PageContent: "# H1\ntext1\n## H2\ntext2\n### H3\ntext3",
		// 		Metadata:    map[string]any{"source": "test.md"},
		// 	},
		// 	expected: []types.Document{
		// 		{PageContent: "H1\ntext1", Metadata: map[string]any{"source": "test.md"}},
		// 		{PageContent: "H2\ntext2", Metadata: map[string]any{"source": "test.md"}},
		// 		{PageContent: "H3\ntext3", Metadata: map[string]any{"source": "test.md"}},
		// 	},
		// },
		{
			name: "Empty Content String",
			doc: types.Document{
				PageContent: "",
				Metadata:    map[string]any{"source": "test.md"},
			},
			expected: nil,
		},
		{
			name: "Content with Only Headers",
			doc: types.Document{
				PageContent: "# H1\n## H2\n### H3",
				Metadata:    map[string]any{"source": "test.md"},
			},
			expected: []types.Document{
				{PageContent: "H1", Metadata: map[string]any{"source": "test.md"}},
				{PageContent: "H2", Metadata: map[string]any{"source": "test.md"}},
				{PageContent: "H3", Metadata: map[string]any{"source": "test.md"}},
			},
		},
		{
			name: "Non-header line with hash symbol",
			doc: types.Document{
				PageContent: "A line with a # symbol.\n\n# Real Header\n\nMore text.",
				Metadata:    map[string]any{"source": "middle_hash.md"},
			},
			expected: []types.Document{
				{PageContent: "A line with a # symbol.", Metadata: map[string]any{"source": "middle_hash.md"}},
				{PageContent: "Real Header\n\nMore text.", Metadata: map[string]any{"source": "middle_hash.md"}},
			},
		},
		{
			name: "Hash inside a code block",
			doc: types.Document{
				PageContent: "# Real Header\n\n```\n# This is not a header\n```",
				Metadata:    map[string]any{"source": "codeblock.md"},
			},
			expected: []types.Document{
				{PageContent: "Real Header\n\n# This is not a header", Metadata: map[string]any{"source": "codeblock.md"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := splitter.SplitDocuments([]types.Document{tc.doc})
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("SplitDocuments() = \n%#v, want \n%#v", result, tc.expected)
			}
		})
	}
}