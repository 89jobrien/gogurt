package state

type InMemState struct {
    current AgentState
}

func (m *InMemState) Save(state AgentState) error {
    m.current = state
    return nil
}

func (m *InMemState) Load() (AgentState, error) {
    return m.current, nil
}
