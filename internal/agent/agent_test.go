package agent

import (
	"context"
	"errors"
	"fmt"
	"gogurt/internal/types"
	"strings"
	"testing"
)

// --- Dummy Planner and AgentDescription for interface compliance ---
// type dummyPlanner struct{}
// type dummyAgentDesc struct{}

// --- Mock Agents ---

type mockAgent struct {
	initErr     error
	invokeRes   any
	invokeErr   error
	asyncVal    any
	asyncErr    error
	delegateRes any
	delegateErr error
}

func (m *mockAgent) Init(ctx context.Context, config types.AgentConfig) error { return m.initErr }
func (m *mockAgent) Invoke(ctx context.Context, input any) (any, error) {
	return m.invokeRes, m.invokeErr
}
func (m *mockAgent) InvokeAsync(ctx context.Context, input any) (<-chan any, <-chan error) {
	ch := make(chan any, 1)
	ech := make(chan error, 1)
	if m.asyncVal != nil {
		ch <- m.asyncVal
	}
	if m.asyncErr != nil {
		ech <- m.asyncErr
	}
	close(ch)
	close(ech)
	return ch, ech
}
func (m *mockAgent) Delegate(ctx context.Context, task any) (any, error) {
	return m.delegateRes, m.delegateErr
}
func (m *mockAgent) Planner() Planner       { return nil }
func (m *mockAgent) State() any             { return nil }
func (m *mockAgent) Capabilities() []string { return []string{"testcap"} }
func (m *mockAgent) Describe() *types.AgentDescription {
	return &types.AgentDescription{Name: "mock", Capabilities: m.Capabilities(), Tools: []string{"toolX"}}
}

func makeTestAgent(initErr, invokeErr, asyncErr, delegateErr error, invokeRes, asyncVal, delegateRes any) Agent {
	return &mockAgent{
		initErr:     initErr,
		invokeErr:   invokeErr,
		invokeRes:   invokeRes,
		asyncErr:    asyncErr,
		asyncVal:    asyncVal,
		delegateErr: delegateErr,
		delegateRes: delegateRes,
	}
}

func TestRegisterAgent_NewAgent_lifecycle_success(t *testing.T) {
	orig := RegisteredAgents
	RegisteredAgents = make(AgentRegistry)
	defer func() { RegisteredAgents = orig }()

	RegisterAgent("testmock", func() Agent {
		return makeTestAgent(nil, nil, nil, nil, "result", "async", "delegated")
	})

	cfg := types.AgentConfig{Name: "testmock"}
	agent, err := NewAgent(cfg)
	if err != nil {
		t.Fatalf("NewAgent error: %v", err)
	}
	res, err := agent.Invoke(context.Background(), "input")
	if err != nil || res != "result" {
		t.Errorf("Invoke: got (%v, %v), want ('result', nil)", res, err)
	}
	err = agent.Init(context.Background(), cfg)
	if err != nil {
		t.Errorf("Init: got %v, want nil", err)
	}
	asyncRes, _ := agent.InvokeAsync(context.Background(), "input")
	val := <-asyncRes
	if val != "async" {
		t.Errorf("InvokeAsync: got %v, want 'async'", val)
	}
	// _, _ = <-asyncErrs
	// delegRes, delegErr := agent.Delegate(context.Background(), "task")
	// if delegRes != "delegated" || delegErr != nil {
	// 	t.Errorf("Delegate: got (%v, %v), want ('delegated', nil)", delegRes, delegErr)
	// }
	// if agent.Planner() != nil {
	// 	t.Errorf("Planner: got non-nil, want nil")
	// }
	if agent.State() != nil {
		t.Errorf("State: got non-nil, want nil")
	}
	// if caps := agent.Capabilities(); len(caps) != 1 || caps[0] != "testcap" {
	// 	t.Errorf("Capabilities: got %v, want [testcap]", caps)
	// }
	desc := agent.Describe()
	if desc == nil || desc.Name != "mock" || desc.Capabilities[0] != "testcap" {
		t.Errorf("Describe: got %v, want mock/testcap", desc)
	}
}

func TestRegisterAgent_overwrite(t *testing.T) {
	orig := RegisteredAgents
	RegisteredAgents = make(AgentRegistry)
	defer func() { RegisteredAgents = orig }()

	RegisterAgent("X", func() Agent { return makeTestAgent(nil, nil, nil, nil, "v1", nil, nil) })
	RegisterAgent("X", func() Agent { return makeTestAgent(nil, nil, nil, nil, "v2", nil, nil) })

	agent, err := NewAgent(types.AgentConfig{Name: "X"})
	if err != nil {
		t.Fatalf("NewAgent error: %v", err)
	}
	v, _ := agent.Invoke(context.Background(), "input")
	if v != "v2" {
		t.Errorf("Agent overwrite, got %v, want v2", v)
	}
}

func TestNewAgentNotFound(t *testing.T) {
	orig := RegisteredAgents
	RegisteredAgents = make(AgentRegistry)
	defer func() { RegisteredAgents = orig }()

	_, err := NewAgent(types.AgentConfig{Name: "notfound"})
	if err == nil || !errors.Is(err, fmt.Errorf("agent not registered: %s", "notfound")) {
		if !strings.Contains(err.Error(), "agent not registered") {
			t.Errorf("NewAgent(notfound) error = %v, want contains 'agent not registered'", err)
		}
	}
}

func TestRegisterAgent_NilFactoryPanicSafe(t *testing.T) {
	orig := RegisteredAgents
	RegisteredAgents = make(AgentRegistry)
	defer func() { RegisteredAgents = orig }()

	// Register nil factory is technically legal, but returns nil Agent on NewAgent
	RegisterAgent("nil", nil)
	agent, err := NewAgent(types.AgentConfig{Name: "nil"})
	if err != nil {
		// Should return error or nil safely.
		return
	}
	if agent == nil {
		return
	}
	t.Errorf("NewAgent(nil) = %v, want nil agent or error", agent)
}

func TestAgentRegistry_multiple_types(t *testing.T) {
	orig := RegisteredAgents
	RegisteredAgents = make(AgentRegistry)
	defer func() { RegisteredAgents = orig }()

	RegisterAgent("a", func() Agent { return makeTestAgent(nil, nil, nil, nil, "aResult", nil, nil) })
	RegisterAgent("b", func() Agent {
		return makeTestAgent(errors.New("bad init"), errors.New("bad invoke"), errors.New("bad async"), errors.New("bad delegate"), nil, nil, nil)
	})

	type testCase struct {
		name          string
		config        types.AgentConfig
		wantErr       bool
		wantInvokeErr bool
		wantInitErr   bool
	}

	cases := []testCase{
		{"A", types.AgentConfig{Name: "a"}, false, false, false},
		{"B", types.AgentConfig{Name: "b"}, false, true, true},
		{"C", types.AgentConfig{Name: "notfound"}, true, false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, err := NewAgent(tc.config)
			if tc.wantErr {
				if err == nil {
					t.Errorf("wanted error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			_, invE := a.Invoke(context.Background(), nil)
			if tc.wantInvokeErr && invE == nil {
				t.Errorf("Invoke: wanted error, got nil")
			}
			if !tc.wantInvokeErr && invE != nil {
				t.Errorf("Invoke: got error %v", invE)
			}
			if err := a.Init(context.Background(), tc.config); (err != nil) != tc.wantInitErr {
				t.Errorf("Init: error got %v, wantInitErr=%v", err, tc.wantInitErr)
			}
		})
	}
}

// Edge case: RegisterAgent empty string name
func TestRegisterAgent_empty_name(t *testing.T) {
	orig := RegisteredAgents
	RegisteredAgents = make(AgentRegistry)
	defer func() { RegisteredAgents = orig }()

	RegisterAgent("", func() Agent { return makeTestAgent(nil, nil, nil, nil, "empty", nil, nil) })
	agent, err := NewAgent(types.AgentConfig{Name: ""})
	if err != nil {
		t.Errorf("NewAgent(\"\") error: %v", err)
	}
	if agent == nil {
		t.Errorf("expected non-nil Agent for empty name registration")
	}
}

func TestRegisterAgent_NilFactoriesDoNotPanic(t *testing.T) {
	orig := RegisteredAgents
	RegisteredAgents = make(AgentRegistry)
	defer func() { RegisteredAgents = orig }()

	// RegisterAgent with nil factory should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RegisterAgent(nil) panicked: %v", r)
		}
	}()
	RegisterAgent("nilfact", nil)
}
