package tools

import (
	"reflect"
)

// AutoInputSchema generates a JSON Schema for a struct type.
func AutoInputSchema(t reflect.Type) map[string]any {
	properties := make(map[string]any)
	required := []string{}

	// Dereference pointer types if necessary
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		name := jsonTag
		if commaIdx := findComma(jsonTag); commaIdx != -1 {
			name = jsonTag[:commaIdx]
		}
		typemap := map[string]any{}
		switch field.Type.Kind() {
		case reflect.String:
			typemap["type"] = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			typemap["type"] = "integer"
		case reflect.Float32, reflect.Float64:
			typemap["type"] = "number"
		case reflect.Bool:
			typemap["type"] = "boolean"
		default:
			typemap["type"] = "string"
		}
		properties[name] = typemap
		required = append(required, name)
	}
	return map[string]any{
		"type":       "object",
		"properties": properties,
		"required":   required,
	}
}

func findComma(tag string) int {
	for i, c := range tag {
		if c == ',' {
			return i
		}
	}
	return -1
}
