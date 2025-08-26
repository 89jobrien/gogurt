package file_tools

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func setup(t *testing.T) {
	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		os.Mkdir(workingDir, 0755)
	}
	// Create a dummy file for reading
	os.WriteFile(filepath.Join(workingDir, "test.txt"), []byte("test content"), 0644)
}

func teardown(t *testing.T) {
	os.RemoveAll(workingDir)
}

func TestReadFileTool(t *testing.T) {
	setup(t)
	defer teardown(t)

	// Test successful read
	result, err := ReadFileTool.Call(`{"filename":"test.txt"}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "test content" {
		t.Errorf("expected 'test content', got '%s'", result)
	}

	// Test file not found
	_, err = ReadFileTool.Call(`{"filename":"notfound.txt"}`)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestWriteFileTool(t *testing.T) {
	setup(t)
	defer teardown(t)

	// Test successful write
	_, err := WriteFileTool.Call(`{"filename":"newfile.txt", "content":"new content"}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	content, err := os.ReadFile(filepath.Join(workingDir, "newfile.txt"))
	if err != nil {
		t.Fatalf("expected no error reading file, got %v", err)
	}
	if string(content) != "new content" {
		t.Errorf("expected 'new content', got '%s'", string(content))
	}
}

func TestListFilesTool(t *testing.T) {
	setup(t)
	defer teardown(t)

	// Test listing files
	result, err := ListFilesTool.Call(`{"path":"."}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := []string{"test.txt"}
	// This is a bit tricky because result is any. We need to cast it.
	resultSlice, ok := result.([]string)
	if !ok {
		t.Fatalf("result is not a []string")
	}

	if !reflect.DeepEqual(resultSlice, expected) {
		t.Errorf("expected %v, got %v", expected, resultSlice)
	}
}