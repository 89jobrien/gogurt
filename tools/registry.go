package tools

import "fmt"

type Registry struct {
    tools map[string]*Tool
}

func NewRegistry() *Registry {
    return &Registry{
        tools: make(map[string]*Tool),
    }
}

// Register a tool under a name. Returns error if name is duplicate.
func (r *Registry) Register(tool *Tool) error {
    if _, exists := r.tools[tool.Name]; exists {
        return fmt.Errorf("tool %q already registered", tool.Name)
    }
    r.tools[tool.Name] = tool
    return nil
}

// Get returns a registered tool by name, or nil if not found.
func (r *Registry) Get(name string) *Tool {
    return r.tools[name]
}

// Call looks up a tool by name and calls it with jsonArgs.
func (r *Registry) Call(name string, jsonArgs string) (any, error) {
    tool := r.Get(name)
    if tool == nil {
        return nil, fmt.Errorf("tool %q not found", name)
    }
    return tool.Call(jsonArgs)
}