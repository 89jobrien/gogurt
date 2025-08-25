package tools

import (
	"errors"
	"fmt"
	"strings"
)

type RegistryStruct struct {
	Count       int
	ToolNames   []string
	Categories  []string
	HasDups     bool
	HasCategory map[string]bool
}

type Registry struct {
	tools map[string]*Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]*Tool),
	}
}

func (r *Registry) Register(tool *Tool) error {
	if tool == nil {
		return errors.New("cannot register nil tool")
	}
	if !isValidToolName(tool.Name) {
		return fmt.Errorf("tool name invalid: %q", tool.Name)
	}
	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool %q already registered", tool.Name)
	}
	r.tools[tool.Name] = tool
	return nil
}

func isValidToolName(name string) bool {
	return name != "" && !strings.ContainsAny(name, " \t\n")
}

func (r *Registry) RegisterBatch(tools []*Tool) []error {
	var errs []error
	for _, tool := range tools {
		if err := r.Register(tool); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (r *Registry) Get(name string) *Tool {
	if !isValidToolName(name) {
		return nil
	}
	return r.tools[name]
}

func (r *Registry) Call(name string, jsonArgs string) (any, error) {
	tool := r.Get(name)
	if tool == nil {
		return nil, fmt.Errorf("tool %q not found", name)
	}
	return tool.Call(jsonArgs)
}

func (r *Registry) List() []string {
	return toolNames(r.tools)
}

func toolNames(m map[string]*Tool) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	return names
}

func (r *Registry) ListTools() []*Tool {
	return toolPointers(r.tools)
}

func toolPointers(m map[string]*Tool) []*Tool {
	list := make([]*Tool, 0, len(m))
	for _, t := range m {
		list = append(list, t)
	}
	return list
}

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

// Stats reports an overview of registry contents for analytics/user feedback.
func (r *Registry) Stats() RegistryStruct {
	toolNames := r.List()
	catSet := make(map[string]struct{})
	hasCategory := make(map[string]bool)
	for _, t := range r.ListTools() {
		cat := ""
		if t.Metadata != nil {
			if c, ok := t.Metadata["category"].(string); ok {
				cat = c
			}
		}
		if cat != "" {
			catSet[cat] = struct{}{}
			hasCategory[cat] = true
		}
	}
	categories := make([]string, 0, len(catSet))
	for c := range catSet {
		categories = append(categories, c)
	}
	return RegistryStruct{
		Count:       len(toolNames),
		ToolNames:   toolNames,
		Categories:  categories,
		HasDups:     false, // Registry does not allow dups, always false.
		HasCategory: hasCategory,
	}
}