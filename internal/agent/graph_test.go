package agent

import (
	"context"
	"errors"
	"gogurt/internal/types"
	"testing"
)

// Mock agent implementation
type mockGraphAgent struct {
	name    string
	child   Agent
	initErr error
	added   []Agent
}

func (a *mockGraphAgent) Init(ctx context.Context, config types.AgentConfig) error { return a.initErr }
func (a *mockGraphAgent) Invoke(ctx context.Context, input any) (any, error)        { return nil, nil }
func (a *mockGraphAgent) InvokeAsync(ctx context.Context, input any) (<-chan any, <-chan error) {
	return nil, nil
}
func (a *mockGraphAgent) Delegate(ctx context.Context, task any) (any, error)   { return nil, nil }
func (a *mockGraphAgent) Planner() Planner                                     { return nil }
func (a *mockGraphAgent) State() any                                           { return nil }
func (a *mockGraphAgent) Capabilities() []string                               { return nil }
func (a *mockGraphAgent) Describe() *types.AgentDescription                    { return nil }

// Implements AddChild if needed
func (a *mockGraphAgent) AddChild(child Agent) {
	a.added = append(a.added, child)
	a.child = child
}

func TestNewAgentGraph_SingleAgent(t *testing.T) {
	ctx := context.Background()
	mock := &mockGraphAgent{name: "A"}
	registry := map[string]func() Agent{
		"A": func() Agent { return mock },
	}
	cfg := types.AgentConfig{Name: "A"}
	graph, err := NewAgentGraph(ctx, cfg, registry)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if graph.Root == nil || graph.Agents["A"] != mock {
		t.Errorf("Root agent mismatch: got %v, want %v", graph.Root, mock)
	}
	if len(graph.Agents) != 1 {
		t.Errorf("Agents length = %d, want 1", len(graph.Agents))
	}
}

func TestNewAgentGraph_MissingRegistry(t *testing.T) {
	ctx := context.Background()
	registry := map[string]func() Agent{}
	cfg := types.AgentConfig{Name: "missing"}
	graph, err := NewAgentGraph(ctx, cfg, registry)
	if err == nil {
		t.Fatalf("expected error for missing agent, got nil")
	}
	if graph != nil {
		t.Errorf("expected nil graph for missing agent, got %v", graph)
	}
	if !errors.Is(err, err) && !containsString(err.Error(), "agent not registered") {
		t.Errorf("error mismatch: got %v", err)
	}
}

func containsString(s, substr string) bool { return s != "" && substr != "" && (len(s) == len(substr) || len(s) > 0 && len(substr) > 0 && (s == substr || (len(substr) < len(s) && s[len(s)-len(substr):] == substr))) }

func TestNewAgentGraph_InitFails(t *testing.T) {
	ctx := context.Background()
	fail := &mockGraphAgent{name: "bad", initErr: errors.New("fail")}
	registry := map[string]func() Agent{
		"bad": func() Agent { return fail },
	}
	cfg := types.AgentConfig{Name: "bad"}
	graph, err := NewAgentGraph(ctx, cfg, registry)
	if err == nil {
		t.Fatalf("expected error for bad init, got nil")
	}
	if graph != nil {
		t.Errorf("expected nil graph for bad init, got %v", graph)
	}
	if !containsString(err.Error(), "fail") {
		t.Errorf("error mismatch: got %v", err)
	}
}

func TestNewAgentGraph_WithChildren(t *testing.T) {
	ctx := context.Background()
	parent := &mockGraphAgent{name: "parent"}
	child := &mockGraphAgent{name: "child"}

	registry := map[string]func() Agent{
		"parent": func() Agent { return parent },
		"child":  func() Agent { return child },
	}
	// Parent with child, AddChild interface supported
	cfg := types.AgentConfig{
		Name: "parent",
		Children: []types.AgentConfig{
			{Name: "child"},
		},
	}
	graph, err := NewAgentGraph(ctx, cfg, registry)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if graph.Root != parent {
		t.Errorf("Root = %v, want %v", graph.Root, parent)
	}
	if graph.Agents["parent"] != parent || graph.Agents["child"] != child {
		t.Errorf("Agent mapping failed: %v", graph.Agents)
	}
	if parent.child != child {
		t.Errorf("AddChild failed: parent.child = %v, want %v", parent.child, child)
	}
}

func TestNewAgentGraph_DeepTree(t *testing.T) {
	ctx := context.Background()
	factory := func(name string) func() Agent {
		return func() Agent {
			return &mockGraphAgent{name: name}
		}
	}
	registry := map[string]func() Agent{
		"A": factory("A"),
		"B": factory("B"),
		"C": factory("C"),
	}
	cfg := types.AgentConfig{
		Name: "A",
		Children: []types.AgentConfig{
			{Name: "B", Children: []types.AgentConfig{
				{Name: "C"},
			}},
		},
	}
	graph, err := NewAgentGraph(ctx, cfg, registry)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if graph.Root == nil || graph.Agents["A"] == nil || graph.Agents["B"] == nil || graph.Agents["C"] == nil {
		t.Errorf("Agent tree missing: %+v", graph)
	}
	if len(graph.Agents) != 3 {
		t.Errorf("Agents has %d nodes, want 3", len(graph.Agents))
	}
}

func TestNewAgentGraph_InterfaceCompliance(t *testing.T) {
	var _ Agent = &mockGraphAgent{}
}