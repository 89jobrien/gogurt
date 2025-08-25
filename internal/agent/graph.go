package agent

import (
	"context"
	"fmt"
	"gogurt/internal/types"
)

type AgentGraph struct {
	Root   Agent
	Agents map[string]Agent
    Nodes map[string]string
    Edges map[string][]string
}

func NewAgentGraph(ctx context.Context, config types.AgentConfig, registry map[string]func() Agent) (*AgentGraph, error) {
	agents := make(map[string]Agent)
	var build func(cfg types.AgentConfig) (Agent, error)
	build = func(cfg types.AgentConfig) (Agent, error) {
		factory, ok := registry[cfg.Name]
		if !ok {
			return nil, fmt.Errorf("agent not registered: %s", cfg.Name)
		}
		a := factory()
		if err := a.Init(ctx, cfg); err != nil {
			return nil, err
		}
		for _, childCfg := range cfg.Children {
			childAgent, err := build(childCfg)
			if err != nil {
				return nil, err
			}
			if aWithChildren, ok := a.(interface{ AddChild(Agent) }); ok {
				aWithChildren.AddChild(childAgent)
			}
		}
		agents[cfg.Name] = a
		return a, nil
	}
	root, err := build(config)
	if err != nil {
		return nil, err
	}
	return &AgentGraph{Root: root, Agents: agents}, nil
}
