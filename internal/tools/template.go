package tools

import "reflect"

type MyArgs struct {
	Field string `json:"Field"`
}

func MyTool(args MyArgs) (string, error) {
	{ /* ... */
	}
	return "...", nil
}

var MyToolInstance = &Tool{
	Name:        "...",
	Description: "...",
	Func:        reflect.ValueOf(MyTool),
	InputSchema: AutoInputSchema(reflect.TypeOf(MyArgs{})),
	Example:     `{"Field":"example"}`,
	Metadata:    map[string]any{"category": "..."},
}
