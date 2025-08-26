package state

import (
	"context"
	"maps"
	"sync"
)

// AgentState defines a generic agent/module state interface.
type AgentState interface {
	Keys() []string
	Get(key string) (any, error)
	Set(key string, value any) error
	Del(key string) error
	Clone() (AgentState, error)
	Serialize() (map[string]any, error)

	AKeys(ctx context.Context) <-chan []string
	AGet(ctx context.Context, key string) (<-chan any, <-chan error)
	ASet(ctx context.Context, key string, value any) <-chan error
	ADel(ctx context.Context, key string) <-chan error
	AClone(ctx context.Context) (<-chan AgentState, <-chan error)
	ASerialize(ctx context.Context) (<-chan map[string]any, <-chan error)
}

// PersistentState enables persistence of AgentState.
type PersistentState interface {
	Save(state AgentState) error
	Load() (AgentState, error)
	
	ASave(ctx context.Context, state AgentState) <-chan error
	ALoad(ctx context.Context) (<-chan AgentState, <-chan error)
}

// MemoryState can be used for both volatile state and in-memory persistence.
type MemoryState struct {
	mu     sync.RWMutex
	data   map[string]any
	backup map[string]any // for in-memory persistence of state snapshot
}

// NewMemoryState creates a new MemoryState.
func NewMemoryState() *MemoryState {
	return &MemoryState{data: make(map[string]any)}
}

// AgentState implementation
func (s *MemoryState) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return mapKeys(s.data)
}
func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (s *MemoryState) Get(key string) (any, error) {
	if key == "" {
		return nil, ErrInvalidKey
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	if !ok {
		return nil, nil
	}
	return val, nil
}

func (s *MemoryState) Set(key string, value any) error {
	if key == "" {
		return ErrInvalidKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *MemoryState) Del(key string) error {
	if key == "" {
		return ErrInvalidKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}

func (s *MemoryState) Clone() (AgentState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cloned := make(map[string]any, len(s.data))
	maps.Copy(cloned, s.data)
	return &MemoryState{data: cloned}, nil
}

func (s *MemoryState) Serialize() (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copied := make(map[string]any, len(s.data))
	maps.Copy(copied, s.data)
	return copied, nil
}

// Async AgentState methods
func (s *MemoryState) AKeys(ctx context.Context) <-chan []string {
	out := make(chan []string, 1)
	go func() {
		defer close(out)
		out <- s.Keys()
	}()
	return out
}

func (s *MemoryState) AGet(ctx context.Context, key string) (<-chan any, <-chan error) {
	valCh := make(chan any, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(valCh)
		defer close(errCh)
		val, err := s.Get(key)
		if err != nil {
			errCh <- err
		} else {
			valCh <- val
		}
	}()
	return valCh, errCh
}

func (s *MemoryState) ASet(ctx context.Context, key string, value any) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- s.Set(key, value)
	}()
	return errCh
}

func (s *MemoryState) ADel(ctx context.Context, key string) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- s.Del(key)
	}()
	return errCh
}

func (s *MemoryState) AClone(ctx context.Context) (<-chan AgentState, <-chan error) {
	stateCh := make(chan AgentState, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(stateCh)
		defer close(errCh)
		cloned, err := s.Clone()
		if err != nil {
			errCh <- err
		} else {
			stateCh <- cloned
		}
	}()
	return stateCh, errCh
}

func (s *MemoryState) ASerialize(ctx context.Context) (<-chan map[string]any, <-chan error) {
	mapCh := make(chan map[string]any, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(mapCh)
		defer close(errCh)
		m, err := s.Serialize()
		if err != nil {
			errCh <- err
		} else {
			mapCh <- m
		}
	}()
	return mapCh, errCh
}

// PersistentState implementation
func (s *MemoryState) Save(state AgentState) error {
	if state == nil {
		return ErrNilState
	}
	serialized, err := state.Serialize()
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.backup = make(map[string]any, len(serialized))
	maps.Copy(s.backup, serialized)
	return nil
}

func (s *MemoryState) Load() (AgentState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.backup == nil {
		return nil, ErrNilState
	}
	restored := make(map[string]any, len(s.backup))
	maps.Copy(restored, s.backup)
	return &MemoryState{data: restored}, nil
}

// Async PersistentState methods
func (s *MemoryState) ASave(ctx context.Context, state AgentState) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- s.Save(state)
	}()
	return errCh
}

func (s *MemoryState) ALoad(ctx context.Context) (<-chan AgentState, <-chan error) {
	stateCh := make(chan AgentState, 1)
	errCh := make(chan error, 1)
	go func() {
		defer close(stateCh)
		defer close(errCh)
		state, err := s.Load()
		if err != nil {
			errCh <- err
		} else {
			stateCh <- state
		}
	}()
	return stateCh, errCh
}

// Error definitions for validation.
var (
	ErrInvalidKey = &StateError{"invalid key"}
	ErrNilState   = &StateError{"nil AgentState"}
)

// StateError is returned for invalid operations on state.
type StateError struct {
	msg string
}
func (e *StateError) Error() string {
	return e.msg
}