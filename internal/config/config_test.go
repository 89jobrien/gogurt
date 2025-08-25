package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("LLM_PROVIDER", "test_provider")
	os.Setenv("AGENT_MAX_ITERATIONS", "20")

	cfg := Load()

	if cfg.LLMProvider != "test_provider" {
		t.Errorf("expected LLM_PROVIDER to be 'test_provider', got %s", cfg.LLMProvider)
	}
	if cfg.AgentMaxIterations != 20 {
		t.Errorf("expected AGENT_MAX_ITERATIONS to be 20, got %d", cfg.AgentMaxIterations)
	}

	os.Unsetenv("LLM_PROVIDER")
	os.Unsetenv("AGENT_MAX_ITERATIONS")
}