package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Tool struct {
	Name        string
	Description string
	Func        reflect.Value
	InputSchema map[string]any
	Example     string
	Metadata    map[string]any
	
}


func (t *Tool) metaEqual(key, expect string) bool {
	if t.Metadata == nil {
		return false
	}
	val, ok := t.Metadata[key].(string)
	return ok && val == expect
}

func (t *Tool) HasCategory(category string) bool   { return t.metaEqual("category", category) }
func (t *Tool) HasVersion(version string) bool     { return t.metaEqual("version", version) }
func (t *Tool) HasAuthor(author string) bool       { return t.metaEqual("author", author) }

func (t *Tool) HasDeprecated() bool {
	if t.Metadata == nil {
		return false
	}
	if d, ok := t.Metadata["deprecated"].(bool); ok {
		return d
	}
	return false
}

func (t *Tool) HasMetadata() bool { return t.Metadata != nil }

// New creates a new Tool from a Go function with a required name argument
func New(name string, f any, description string) (*Tool, error) {
	val := reflect.ValueOf(f)
	if val.Kind() != reflect.Func {
		return nil, fmt.Errorf("provided interface is not a function")
	}
	t := val.Type()
	if t.NumIn() != 1 || t.NumOut() != 2 {
		return nil, fmt.Errorf("function must have exactly one input and two outputs (result, error)")
	}
	inputType := t.In(0)
	inputSchema, err := generateSchema(inputType)
	if err != nil {
		return nil, fmt.Errorf("error generating schema for input type: %w", err)
	}
	return &Tool{
		Name:        name,
		Description: description,
		Func:        val,
		InputSchema: inputSchema,
	}, nil
}

// Call executes the tool with the given JSON arguments.
func (t *Tool) Call(jsonArgs string) (any, error) {
	inputType := t.Func.Type().In(0)
	inputValue := reflect.New(inputType).Interface()
	if err := json.Unmarshal([]byte(jsonArgs), &inputValue); err != nil {
		return nil, fmt.Errorf("error unmarshaling arguments: %w", err)
	}
	results := t.Func.Call([]reflect.Value{reflect.ValueOf(inputValue).Elem()})
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	return results[0].Interface(), nil
}

// generateSchema creates a JSON schema from a Go type
func generateSchema(t reflect.Type) (map[string]any, error) {
	schema := map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	}
	props := schema["properties"].(map[string]any)
	var required []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}
		prop := map[string]any{}
		switch field.Type.Kind() {
		case reflect.String:
			prop["type"] = "string"
		case reflect.Int, reflect.Int64, reflect.Float64:
			prop["type"] = "number"
		case reflect.Bool:
			prop["type"] = "boolean"
		}
		props[jsonTag] = prop
		required = append(required, jsonTag)
	}
	schema["required"] = required
	return schema, nil
}

// Describe returns a detailed summary of the tool, including its input schema and signature.
func (t *Tool) Describe() string {
	inputSchemaBytes, _ := json.MarshalIndent(t.InputSchema, "", "  ")

	meta := ""
	if t.Metadata != nil {
		metaBytes, _ := json.MarshalIndent(t.Metadata, "", "  ")
		meta = fmt.Sprintf("Metadata:\n%s\n", string(metaBytes))
	}

	example := ""
	if t.Example != "" {
		example = fmt.Sprintf("Example Input:\n%s\n", t.Example)
	}

	return fmt.Sprintf(
		"==========\nTool: %s\nDescription: %s\nFunction Signature: %s\n\nInput Schema:\n%s\n\n%s%s==========\n",
		t.Name,
		t.Description,
		t.Func.Type().String(),
		string(inputSchemaBytes),
		example,
		meta,
	)
}

// Usage
/*
tool := &Tool{
    Name:        "add",
    Description: "Adds two numbers.",
    Func:        reflect.ValueOf(Add),
    InputSchema: map[string]any{
        "type": "object",
        "properties": map[string]any{
            "a": map[string]any{"type": "number"},
            "b": map[string]any{"type": "number"},
        },
        "required": []string{"a", "b"},
    },
    Example:  `{"a": 1, "b": 2}`,
    Metadata: map[string]any{
		"author": "joe",
        "category": "math",
        "version": "1.0",
		"deprecated": false,
    },
}

fmt.Println(tool.Describe())
*/
