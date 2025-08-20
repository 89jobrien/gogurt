package chroma

import (
	"context"
	"os"
	"testing"
	"time"

	chromadb "github.com/amikos-tech/chroma-go/pkg/api/v2"
)

func uniqueCollectionName(prefix string) string {
	return prefix + "-" + time.Now().UTC().Format("20060102T150405.000000000")
}

func ensureClient(t *testing.T) chromadb.Client {
	t.Helper()
	client, err := chromadb.NewHTTPClient()
	if err != nil {
		t.Skipf("Skipping: cannot create Chroma HTTP client: %v", err)
		return nil
	}
	t.Cleanup(func() {
		_ = client.Close()
	})
	return client
}

func TestCreateGetCollectionWithMetadata(t *testing.T) {
	// Optional: allow disabling integration tests easily
	if os.Getenv("CHROMA_INTEGRATION_TESTS") == "0" {
		t.Skip("Integration tests disabled by CHROMA_INTEGRATION_TESTS=0")
	}

	client := ensureClient(t)
	ctx := context.Background()

	name := uniqueCollectionName("col-meta")

	col, err := client.GetOrCreateCollection(ctx, name,
		chromadb.WithCollectionMetadataCreate(
			chromadb.NewMetadata(
				chromadb.NewStringAttribute("str", "hello"),
				chromadb.NewIntAttribute("int", 1),
				chromadb.NewFloatAttribute("float", 1.1),
			),
		),
	)
	if col == nil {
		t.Fatal("failed to create collection")
	}
	t.Logf("created collection: %+v", col)


	if err != nil {
		t.Fatalf("GetOrCreateCollection failed: %v", err)
	}

	// Get the same collection again and ensure it resolves
	col2, err := client.GetOrCreateCollection(ctx, name)
	if err != nil {
		t.Fatalf("GetOrCreateCollection (second) failed: %v", err)
	}
	if col2 == nil {
		t.Fatalf("Expected non-nil collection on second get")
	}
}

func TestAddCountQueryAndDelete(t *testing.T) {
	if os.Getenv("CHROMA_INTEGRATION_TESTS") == "0" {
		t.Skip("Integration tests disabled by CHROMA_INTEGRATION_TESTS=0")
	}

	client := ensureClient(t)
	ctx := context.Background()

	name := uniqueCollectionName("col-ops")

	col, err := client.GetOrCreateCollection(ctx, name,
		chromadb.WithCollectionMetadataCreate(
			chromadb.NewMetadata(
				chromadb.NewStringAttribute("str", "hello"),
				chromadb.NewIntAttribute("int", 1),
				chromadb.NewFloatAttribute("float", 1.1),
			),
		),
	)
	if err != nil {
		t.Fatalf("GetOrCreateCollection failed: %v", err)
	}

	// Add two docs with IDs, texts, and per-document metadata
	err = col.Add(ctx,
		chromadb.WithIDs("1", "2"),
		chromadb.WithTexts("hello world", "goodbye world"),
		chromadb.WithMetadatas(
			chromadb.NewDocumentMetadata(chromadb.NewIntAttribute("int", 1)),
			chromadb.NewDocumentMetadata(chromadb.NewStringAttribute("str", "hello")),
		),
	)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Count should be 2
	count, err := col.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count != 2 {
		t.Fatalf("Expected count=2, got %d", count)
	}

	// Query by text; expect at least one result back
	qr, err := col.Query(ctx, chromadb.WithQueryTexts("say hello"))
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	// Validate query result structure conservatively
	docGroups := qr.GetDocumentsGroups()
	if len(docGroups) == 0 {
		t.Fatalf("Expected at least one documents group; got 0")
	}
	if len(docGroups[0]) == 0 {
		t.Fatalf("Expected at least one document in first group; got 0")
	}

	// Delete both items by ID
	err = col.Delete(ctx, chromadb.WithIDsDelete("1", "2"))
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Count should be 0
	countAfter, err := col.Count(ctx)
	if err != nil {
		t.Fatalf("Count (after delete) failed: %v", err)
	}
	if countAfter != 0 {
		t.Fatalf("Expected count=0 after delete, got %d", countAfter)
	}
}

func TestIdempotentCreateAndCleanup(t *testing.T) {
	if os.Getenv("CHROMA_INTEGRATION_TESTS") == "0" {
		t.Skip("Integration tests disabled by CHROMA_INTEGRATION_TESTS=0")
	}

	client := ensureClient(t)
	ctx := context.Background()

	name := uniqueCollectionName("col-idem")

	// First creation
	col, err := client.GetOrCreateCollection(ctx, name)
	if err != nil {
		t.Fatalf("first GetOrCreateCollection failed: %v", err)
	}
	// Second call should return same collection handle without error
	col2, err := client.GetOrCreateCollection(ctx, name)
	if err != nil {
		t.Fatalf("second GetOrCreateCollection failed: %v", err)
	}
	if col2 == nil || col == nil {
		t.Fatalf("collections should be non-nil")
	}

	// Add and delete a quick doc to verify collection is functional
	err = col.Add(ctx,
		chromadb.WithIDs("x1"),
		chromadb.WithTexts("transient"),
	)
	if err != nil {
		t.Fatalf("Add transient doc failed: %v", err)
	}
	c, err := col.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if c != 1 {
		t.Fatalf("expected count=1, got %d", c)
	}
	if err := col.Delete(ctx, chromadb.WithIDsDelete("x1")); err != nil {
		t.Fatalf("Delete transient doc failed: %v", err)
	}
}

