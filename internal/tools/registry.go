package tools

import (
	"context"
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

// Async Register
func (r *Registry) ARegister(ctx context.Context, tool *Tool) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- r.Register(tool)
	}()
	return errCh
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

// Async RegisterBatch
func (r *Registry) ARegisterBatch(ctx context.Context, tools []*Tool) <-chan []error {
	errCh := make(chan []error, 1)
	go func() {
		defer close(errCh)
		errCh <- r.RegisterBatch(tools)
	}()
	return errCh
}

func (r *Registry) Get(name string) *Tool {
	if !isValidToolName(name) {
		return nil
	}
	return r.tools[name]
}

// Async Get
func (r *Registry) AGet(ctx context.Context, name string) <-chan *Tool {
	out := make(chan *Tool, 1)
	go func() {
		defer close(out)
		out <- r.Get(name)
	}()
	return out
}

func (r *Registry) Call(name string, jsonArgs string) (any, error) {
	tool := r.Get(name)
	if tool == nil {
		return nil, fmt.Errorf("tool %q not found", name)
	}
	return tool.Call(jsonArgs)
}

// Async Call
func (r *Registry) ACall(ctx context.Context, name string, jsonArgs string) (<-chan any, <-chan error) {
	out := make(chan any, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errCh)
		result, err := r.Call(name, jsonArgs)
		if err != nil {
			errCh <- err
		} else {
			out <- result
		}
	}()
	return out, errCh
}

func (r *Registry) List() []string {
	return toolNames(r.tools)
}

// Async List
func (r *Registry) AList(ctx context.Context) <-chan []string {
	out := make(chan []string, 1)
	go func() {
		defer close(out)
		out <- r.List()
	}()
	return out
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

// Async ListTools
func (r *Registry) AListTools(ctx context.Context) <-chan []*Tool {
	out := make(chan []*Tool, 1)
	go func() {
		defer close(out)
		out <- r.ListTools()
	}()
	return out
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

// Async PrintAllDescs (returns descriptions)
func (r *Registry) APrintAllDescs(ctx context.Context) <-chan []string {
	out := make(chan []string, 1)
	go func() {
		defer close(out)
		descs := []string{}
		for _, tool := range r.ListTools() {
			descs = append(descs, tool.Describe())
		}
		out <- descs
	}()
	return out
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

// Async GetByCategory
func (r *Registry) AGetByCategory(ctx context.Context, category string) <-chan []*Tool {
	out := make(chan []*Tool, 1)
	go func() {
		defer close(out)
		out <- r.GetByCategory(category)
	}()
	return out
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

// Async Stats
func (r *Registry) AStats(ctx context.Context) <-chan RegistryStruct {
	out := make(chan RegistryStruct, 1)
	go func() {
		defer close(out)
		out <- r.Stats()
	}()
	return out
}