package character

import (
	"gogurt/internal/types"
	"reflect"
	"testing"
)

func TestCharSplitter_SplitDocuments(t *testing.T) {
	splitter := New(10, 2)
	doc := types.Document{PageContent: "abcdefghijklmnopqrstuvwxyz"}
	expected := []types.Document{
		{PageContent: "abcdefghij"},
		{PageContent: "ijklmnopqr"},
		{PageContent: "qrstuvwxyz"},
	}

	result := splitter.SplitDocuments([]types.Document{doc})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SplitDocuments() = %v, want %v", result, expected)
	}
}
