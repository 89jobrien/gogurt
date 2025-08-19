package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gogurt/agent"
	"gogurt/llm/azure"
	"gogurt/llm/ollama"
	"gogurt/llm/openai"
	"gogurt/tools"
	"gogurt/types"
)

// WeatherInput defines the input schema for the GetWeather tool.
type WeatherInput struct {
	City string `json:"city"`
}

// GetWeather is a tool that returns the weather for a city.
func GetWeather(input WeatherInput) (string, error) {
	if input.City == "New York" {
		return "The weather in New York is sunny.", nil
	}
	return fmt.Sprintf("I don't know the weather for %s", input.City), nil
}

// getLLM initializes and returns the selected LLM client.
func getLLM() types.LLM {
	llmProvider := os.Getenv("LLM_PROVIDER")
	switch llmProvider {
	case "azure":
		fmt.Println("Using Azure OpenAI LLM")
		return azure.New()
	case "ollama":
		fmt.Println("Using Ollama LLM")
		ollamaLLM, err := ollama.New()
		if err != nil {
			log.Fatalf("failed to create ollama client: %v", err)
		}
		return ollamaLLM
	default:
		fmt.Println("Using OpenAI LLM")
		return openai.New()
	}
}

func main() {
	// Initialize the LLM based on the environment variable
	llmClient := getLLM()

	// Create a new tool for getting the weather
	weatherTool, err := tools.New(GetWeather, "Get the weather for a city")
	if err != nil {
		log.Fatalf("failed to create weather tool: %v", err)
	}

	// Create a new agent
	aiAgent := agent.New(llmClient, weatherTool)

	// Run the agent with a prompt
	prompt := "What's the weather like in New York?"
	response, err := aiAgent.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("agent run failed: %v", err)
	}

	fmt.Println("AI Response:", response)
}