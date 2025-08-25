package agent

// Supports tool/function chaining and metadata propagation
type AgentCallResult struct {
    Output   string
    Error    error
    Metadata map[string]interface{}
    Next     *AgentCallResult
}