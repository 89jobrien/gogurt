package pipes

import (
	"encoding/json"
	"gogurt/internal/tools"
)

type ToolPipe struct {
	Registry    *tools.Registry
	Names       []string
	InitialJSON string
}

func NewToolPipe(reg *tools.Registry, names []string, initialJSON string) (any, error) {
	var current any = initialJSON
	var err error
	for _, name := range names {
		current, err = reg.Call(name, current.(string))
		if err != nil {
			return nil, err
		}
		// Optionally serialize for JSON input where expected
		b, _ := json.Marshal(current)
		current = string(b)
	}
	return current, nil
}

func (tp *ToolPipe) Run() (any, error) {
	return NewToolPipe(tp.Registry, tp.Names, tp.InitialJSON)
}
