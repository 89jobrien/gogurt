package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LLMProvider string
	OllamaHost string
	OllamaModel string
	AzureOpenAIEndpoint string
	AzureOpenAIAPIKey string
	AzureDeployment string
	OpenAIAPIKey string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		LLMProvider: getEnv("LLM_PROVIDER", "openai"),
		OllamaHost: getEnv("OLLAMA_HOST", "http://localhost:11434"),
		OllamaModel: getEnv("OLLAMA_MODEL", "llama3"),
		AzureOpenAIEndpoint: getEnv("AZURE_OPENAI_ENDPOINT", ""),
		AzureOpenAIAPIKey: getEnv("AZURE_OPENAI_API_KEY", ""),
		AzureDeployment: getEnv("AZURE_OPENAI_DEPLOYMENT_NAME", ""),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}