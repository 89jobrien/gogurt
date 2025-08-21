package state

import (
	"maps"
	"sync"
)

type InMemoryState struct {
    data map[string]any
    mu   sync.RWMutex
}

// Constructor
func NewInMemoryState() *InMemoryState {
    return &InMemoryState{
        data: make(map[string]any),
    }
}

// Keys returns all keys in the state
func (s *InMemoryState) Keys() []string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    keys := make([]string, 0, len(s.data))
    for k := range s.data {
        keys = append(keys, k)
    }
    return keys
}

// Get returns value for key, or nil if not present
func (s *InMemoryState) Get(key string) (any, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.data[key]
    if !ok {
        return nil, nil
    }
    return v, nil
}

// Set sets value for key
func (s *InMemoryState) Set(key string, value any) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = value
    return nil
}

// Del removes a key
func (s *InMemoryState) Del(key string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.data, key)
    return nil
}

// Clone produces a deep copy
func (s *InMemoryState) Clone() (AgentState, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    cloned := make(map[string]any, len(s.data))
    // For simple types; deep copy needed for complex values
	maps.Copy(cloned, s.data)
    return &InMemoryState{data: cloned}, nil
}

// Serialize returns a copy as map[string]any
func (s *InMemoryState) Serialize() (map[string]any, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    copied := make(map[string]any, len(s.data))
    maps.Copy(copied, s.data)
    return copied, nil
}