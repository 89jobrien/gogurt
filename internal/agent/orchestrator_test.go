package agent

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"gogurt/internal/types"
)

func makeCallResult(val string, err error) *types.AgentCallResult {
	return &types.AgentCallResult{Output: val, Error: err}
}

func errorsEqual(a, b error) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Error() == b.Error()
}

// func callResultEqual(a, b *types.AgentCallResult) bool {
// 	if (a == nil) != (b == nil) {
// 		return false
// 	}
// 	if a == nil && b == nil {
// 		return true
// 	}
// 	return a.Output == b.Output && errorsEqual(a.Error, b.Error)
// }

func agentCallResultSliceEqual(got, want []*types.AgentCallResult) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		a, b := got[i], want[i]
		if (a == nil) != (b == nil) {
			return false
		}
		if a == nil && b == nil {
			continue
		}
		if a.Output != b.Output || !errorsEqual(a.Error, b.Error) {
			return false
		}
	}
	return true
}

func TestOrchestrator_RunParallel(t *testing.T) {
	ctx := context.Background()
	success := &mockAgent{invokeRes: makeCallResult("success", nil), invokeErr: nil}
	fail := &mockAgent{invokeRes: nil, invokeErr: errors.New("fail")}
	invalid := &mockAgent{invokeRes: "not-call-result", invokeErr: nil}

	tests := []struct {
		name    string
		agents  []Agent
		input   string
		want    []*types.AgentCallResult
		wantErr bool
	}{
		{
			name:    "all success",
			agents:  []Agent{success, success},
			input:   "input",
			want:    []*types.AgentCallResult{makeCallResult("success", nil), makeCallResult("success", nil)},
			wantErr: false,
		},
		{
			name:    "one error/one success",
			agents:  []Agent{success, fail},
			input:   "input",
			want:    []*types.AgentCallResult{makeCallResult("success", nil), {Error: errors.New("fail")}},
			wantErr: true,
		},
		{
			name:    "all error",
			agents:  []Agent{fail, fail},
			input:   "input",
			want:    []*types.AgentCallResult{{Error: errors.New("fail")}, {Error: errors.New("fail")}},
			wantErr: true,
		},
		{
			name:    "invalid result type",
			agents:  []Agent{invalid},
			input:   "input",
			want:    []*types.AgentCallResult{{Error: fmt.Errorf("invalid AgentCallResult (got %T)", "not-call-result")}},
			wantErr: true,
		},
		{
			name:    "zero agents",
			agents:  []Agent{},
			input:   "input",
			want:    []*types.AgentCallResult{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Orchestrator{Agents: tt.agents}
			got, err := o.RunParallel(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunParallel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !agentCallResultSliceEqual(got, tt.want) {
				t.Errorf("RunParallel() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

type chainedMock struct {
	responses []any
	errors    []error
	index     int
}

func (m *chainedMock) Invoke(ctx context.Context, input any) (any, error) {
	if m.index >= len(m.responses) {
		return nil, nil
	}
	resp := m.responses[m.index]
	var err error
	if m.errors != nil && m.index < len(m.errors) {
		err = m.errors[m.index]
	}
	m.index++
	return resp, err
}

// Satisfy the Agent interface
func (m *chainedMock) Capabilities() []string                                   { return nil }
func (m *chainedMock) Delegate(ctx context.Context, task any) (any, error)      { return nil, nil }
func (m *chainedMock) Init(ctx context.Context, config types.AgentConfig) error { return nil }
func (m *chainedMock) InvokeAsync(ctx context.Context, input any) (<-chan any, <-chan error) {
	return nil, nil
}
func (m *chainedMock) Planner() Planner                  { return nil }
func (m *chainedMock) State() any                        { return nil }
func (m *chainedMock) Describe() *types.AgentDescription { return nil }

func TestOrchestrator_RunPiped(t *testing.T) {
	ctx := context.Background()
	chained := func(results []any, errs []error) Agent {
		return &chainedMock{responses: results, errors: errs}
	}
	tests := []struct {
		name    string
		agents  []Agent
		input   string
		want    *types.AgentCallResult
		wantErr bool
	}{
		{
			name:    "all success",
			agents:  []Agent{chained([]any{makeCallResult("foo", nil), makeCallResult("bar", nil)}, nil)},
			input:   "input",
			want:    makeCallResult("foo", nil),
			wantErr: false,
		},
		{
			name:    "error in chain",
			agents:  []Agent{chained([]any{makeCallResult("ok", nil), makeCallResult("fail", errors.New("bad"))}, []error{nil, errors.New("bad")})},
			input:   "input",
			want:    &types.AgentCallResult{Output: "fail", Error: errors.New("bad")},
			wantErr: true,
		},
		{
			name:    "invalid result type",
			agents:  []Agent{chained([]any{"wrong-type"}, nil)},
			input:   "input",
			want:    &types.AgentCallResult{Error: fmt.Errorf("invalid AgentCallResult (got %T)", "wrong-type")},
			wantErr: true,
		},
		{
			name:    "zero agents",
			agents:  []Agent{},
			input:   "input",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Orchestrator{Agents: tt.agents}
			got, err := o.RunPiped(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				return
			}
			if !agentCallResultSliceEqual([]*types.AgentCallResult{got}, []*types.AgentCallResult{tt.want}) {
				t.Errorf("RunPiped() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
