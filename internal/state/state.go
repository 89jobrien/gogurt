package state

// AgentState is interface for agent or module state
type AgentState interface {
	Keys() []string
	Get(key string) (any, error)
	Set(key string, value any) error
	Del(key string) error
	Clone() (AgentState, error)
	Serialize() (map[string]any, error)
}

type PersistentState interface {
	Save(state AgentState) error
	Load() (AgentState, error)
}
