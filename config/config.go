package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	LLMProvider         string
	OllamaHost          string
	OllamaModel         string
	OllamaEmbedModel    string
	AzureOpenAIEndpoint string
	AzureOpenAIAPIKey   string
	AzureDeployment     string
	OpenAIAPIKey        string
	AgentMaxIterations int
	SplitterProvider   string
	VectorStoreProvider string
	ChromaURL           string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	maxIter, _ := strconv.Atoi(getEnv("AGENT_MAX_ITERATIONS", "10"))

	return &Config{
		LLMProvider:         getEnv("LLM_PROVIDER", "openai"),
		OllamaHost:          getEnv("OLLAMA_HOST", "http://localhost:11434"),
		OllamaModel:         getEnv("OLLAMA_MODEL", "llama3.2:3b"),
		OllamaEmbedModel:    getEnv("OLLAMA_EMBED_MODEL", "llama3.2:3b"),
		AzureOpenAIEndpoint: getEnv("AZURE_OPENAI_ENDPOINT", ""),
		AzureOpenAIAPIKey:   getEnv("AZURE_OPENAI_API_KEY", ""),
		AzureDeployment:     getEnv("AZURE_OPENAI_DEPLOYMENT_NAME", ""),
		OpenAIAPIKey:        getEnv("OPENAI_API_KEY", ""),
		AgentMaxIterations: maxIter,
		SplitterProvider:   getEnv("SPLITTER_PROVIDER", "recursive"),
		VectorStoreProvider: getEnv("VECTOR_STORE_PROVIDER", "faiss"),
		ChromaURL:           getEnv("CHROMA_URL", "http://localhost:8000"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}