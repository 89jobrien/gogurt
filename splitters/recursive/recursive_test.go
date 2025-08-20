package recursive

import (
	"gogurt/types"
	"reflect"
	"testing"
)

func TestRecursiveSplitter_SplitDocuments(t *testing.T) {
	splitter := New(20, 5)
	doc := types.Document{PageContent: "This is a long sentence for testing the recursive splitter."}
	expected := []types.Document{
		{PageContent: "This is a long"},
		{PageContent: "long sentence for"},
		{PageContent: "for testing the"},
		{PageContent: "the recursive"},
		{PageContent: "splitter."},
	}

	result := splitter.SplitDocuments([]types.Document{doc})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SplitDocuments() = %v, want %v", result, expected)
	}
}