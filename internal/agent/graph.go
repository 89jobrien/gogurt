package agent

import (
	"context"
	"fmt"
	"gogurt/internal/types"
	"sync"
)

type AgentGraph struct {
	Root   Agent
	Agents map[string]Agent
	Nodes  map[string]string
	Edges  map[string][]string
}

// NewAgentGraph builds each agent (and its children) asynchronously
func NewAgentGraph(ctx context.Context, config types.AgentConfig, registry map[string]func() Agent) (*AgentGraph, error) {
	agents := make(map[string]Agent)
	var mu sync.Mutex
	var build func(cfg types.AgentConfig) (Agent, error)

	errorsCh := make(chan error, 128)

	build = func(cfg types.AgentConfig) (Agent, error) {
		factory, ok := registry[cfg.Name]
		if !ok {
			return nil, fmt.Errorf("agent not registered: %s", cfg.Name)
		}
		a := factory()
		if err := a.Init(ctx, cfg); err != nil {
			return nil, err
		}
		childAgents := make([]Agent, len(cfg.Children))

		if len(cfg.Children) > 0 {
			childErrs := make([]error, len(cfg.Children))
			var childWg sync.WaitGroup

			for i, childCfg := range cfg.Children {
				childWg.Add(1)
				go func(i int, childCfg types.AgentConfig) {
					defer childWg.Done()
					childAgent, err := build(childCfg)
					if err != nil {
						childErrs[i] = err
						errorsCh <- err
						return
					}
					childAgents[i] = childAgent
				}(i, childCfg)
			}
			childWg.Wait()
			for _, err := range childErrs {
				if err != nil {
					return nil, err
				}
			}

			if aWithChildren, ok := a.(interface{ AddChild(Agent) }); ok {
				for _, child := range childAgents {
					aWithChildren.AddChild(child)
				}
			}
		}
		mu.Lock()
		agents[cfg.Name] = a
		mu.Unlock()
		return a, nil
	}

	root, err := build(config)
	close(errorsCh)
	for e := range errorsCh {
		if e != nil {
			return nil, e
		}
	}
	if err != nil {
		return nil, err
	}
	return &AgentGraph{Root: root, Agents: agents}, nil
}