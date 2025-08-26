package file_tools

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"gogurt/internal/tools"
)

const workingDir = "WORKING_DIR"

func init() {
	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		os.Mkdir(workingDir, 0755)
	}
}

// --- ReadFile Tool ---

type ReadFileArgs struct {
	Filename string `json:"filename"`
}

func ReadFile(args ReadFileArgs) (string, error) {
	path := filepath.Join(workingDir, args.Filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", path)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

var ReadFileTool = &tools.Tool{
	Name:        "read_file",
	Description: "Reads the content of a file from the working directory.",
	Func:        reflect.ValueOf(ReadFile),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(ReadFileArgs{})),
	Example:     `{"filename":"example.txt"}`,
	Metadata:    map[string]any{"category": "file"},
}

// --- WriteFile Tool ---

type WriteFileArgs struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func WriteFile(args WriteFileArgs) (string, error) {
	path := filepath.Join(workingDir, args.Filename)
	err := os.WriteFile(path, []byte(args.Content), 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Successfully wrote to %s", args.Filename), nil
}

var WriteFileTool = &tools.Tool{
	Name:        "write_file",
	Description: "Writes content to a file in the working directory.",
	Func:        reflect.ValueOf(WriteFile),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(WriteFileArgs{})),
	Example:     `{"filename":"example.txt", "content":"Hello, World!"}`,
	Metadata:    map[string]any{"category": "file"},
}

// --- ListFiles Tool ---

type ListFilesArgs struct {
	Path string `json:"path"`
}

func ListFiles(args ListFilesArgs) ([]string, error) {
	path := workingDir
	if args.Path != "" {
		path = filepath.Join(workingDir, args.Path)
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	return filenames, nil
}

var ListFilesTool = &tools.Tool{
	Name:        "list_files",
	Description: "Lists files in a directory within the working directory.",
	Func:        reflect.ValueOf(ListFiles),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(ListFilesArgs{})),
	Example:     `{"path":"."}`,
	Metadata:    map[string]any{"category": "file"},
}