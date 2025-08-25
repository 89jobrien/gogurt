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

// List returns a slice of all registered tool names.
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// RegisterBatch registers all provided tools, returning a slice of errors for duplicates.
func (r *Registry) RegisterBatch(tools []*Tool) []error {
	var errs []error
	for _, tool := range tools {
		if err := r.Register(tool); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// ListTools returns a slice of all registered tool pointers.
func (r *Registry) ListTools() []*Tool {
	tools := make([]*Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// PrintAllDescs pretty prints metadata for all tools in the registry.
func (r *Registry) PrintAllDescs() {
	for _, tool := range r.ListTools() {
		fmt.Println(tool.Describe())
	}
}

func (r *Registry) GetByCategory(category string) []*Tool {
	var matches []*Tool
	for _, tool := range r.tools {
		if tool.HasCategory(category) {
			matches = append(matches, tool)
		}
	}
	return matches
}
