package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"gogurt/agent"
	"gogurt/config"
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
func getLLM(cfg *config.Config) types.LLM {
	var llm types.LLM
	var err error

	switch cfg.LLMProvider {
	case "azure":
		fmt.Println("Using Azure LLM")
		llm, err = azure.New(cfg)
	case "ollama":
		fmt.Println("Using Ollama LLM")
		llm, err = ollama.New(cfg)
	default:
		fmt.Println("Using OpenAI LLM")
		llm, err = openai.New(cfg)
	}

	if err != nil {
		log.Fatalf("failed to create LLM: %v", err)
	}
	return llm
}

func main() {
	// Initialize the configuration
	cfg := config.Load()

	// Initialize the LLM based on the configuration
	llmClient := getLLM(cfg)

	// Create a new tool for getting the weather
	weatherTool, err := tools.New(GetWeather, "Get the weather for a city")
	if err != nil {
		log.Fatalf("failed to create weather tool: %v", err)
	}

	// Create a new agent
	aiAgent := agent.New(llmClient, weatherTool)

	// Read prompt from the command line
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a prompt: ")
	prompt, _ := reader.ReadString('\n')
	prompt = strings.TrimSpace(prompt)

	// Run the agent with the prompt
	response, err := aiAgent.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("agent run failed: %v", err)
	}

	fmt.Println("AI Response:", response)
}