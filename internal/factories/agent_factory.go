package factories

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"gogurt/internal/agent"
	"gogurt/internal/types"
)

// Loads AgentConfig from a JSON file.
func LoadAgentConfigFromFile(filename string) (types.AgentConfig, error) {
	var cfg types.AgentConfig
	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Factory for composing agents recursively from config
func ComposeAgentFromConfig(
	ctx context.Context,
	cfg types.AgentConfig,
) (agent.Agent, error) {
	reg, ok := agent.RegisteredAgents[cfg.Name]
	if !ok {
		return nil, fmt.Errorf("no such agent type: %s", cfg.Name)
	}
	a := reg()
	if err := a.Init(ctx, cfg); err != nil {
		return nil, err
	}
	// Recursively compose child agents, if any (e.g. for pipelines/graphs)
	for _, childCfg := range cfg.Children {
		child, err := ComposeAgentFromConfig(ctx, childCfg)
		if err != nil {
			return nil, fmt.Errorf("compose child agent '%s' failed: %w", childCfg.Name, err)
		}
		// How to attach depends on your agent structure:
		// E.g. if agent.Agent supports AddChild(child agent.Agent)
		if ac, ok := a.(interface{ AddChild(agent.Agent) }); ok {
			ac.AddChild(child)
		}
	}
	return a, nil
}
