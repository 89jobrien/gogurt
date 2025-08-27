package file_tools

import (
	"fmt"
	"gogurt/internal/tools"
	"os"
	"path/filepath"
	"reflect"
)

// --- ReadScratchpad Tool ---

type ReadScratchpadArgs struct {
	Filename string `json:"filename"`
}

func ReadScratchpad(args ReadScratchpadArgs) (string, error) {
	if args.Filename == "" {
		args.Filename = "scratchpad.txt"
	}
	return ReadFile(ReadFileArgs{Filename: args.Filename})
}

var ReadScratchpadTool = &tools.Tool{
	Name:        "read_scratchpad",
	Description: "Reads the content of a file from the working directory, defaulting to scratchpad.txt.",
	Func:        reflect.ValueOf(ReadScratchpad),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(ReadScratchpadArgs{})),
	Example:     `{"filename":"my_notes.txt"}`,
	Metadata:    map[string]any{"category": "file"},
}

// --- SaveToScratchpad Tool ---

type SaveToScratchpadArgs struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
}

func SaveToScratchpad(args SaveToScratchpadArgs) (string, error) {
	if args.Filename == "" {
		args.Filename = "scratchpad.txt"
	}
	path := filepath.Join(workingDir, args.Filename)
	err := os.WriteFile(path, []byte(args.Content), 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Successfully saved to %s", args.Filename), nil
}

var SaveToScratchpadTool = &tools.Tool{
	Name:        "save_to_scratchpad",
	Description: "Saves content to a file in the working directory, defaulting to scratchpad.txt.",
	Func:        reflect.ValueOf(SaveToScratchpad),
	InputSchema: tools.GenInputSchema(reflect.TypeOf(SaveToScratchpadArgs{})),
	Example:     `{"content":"some important info", "filename":"my_notes.txt"}`,
	Metadata:    map[string]any{"category": "file"},
}