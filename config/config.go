package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
    LLMProvider          string
    OllamaHost           string
    OllamaModel          string
    OllamaEmbedModel     string
    AzureOpenAIEndpoint  string
    AzureOpenAIAPIKey    string
    AzureDeployment      string
    OpenAIAPIKey         string
    AgentMaxIterations   int
    SplitterProvider     string
    VectorStoreProvider  string
    ChromaURL            string
    ChromaSpace         string
    ChromaCollection    string
    ChromaTenant        string
    ChromaDatabase      string
    ChromaEFConstruction int
    ChromaEFSearch       int
    ChromaMaxNeighbors   int
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    maxIterStr := getEnv("AGENT_MAX_ITERATIONS", "10")
    maxIter, err := strconv.Atoi(maxIterStr)
    if err != nil {
        log.Printf("Invalid AGENT_MAX_ITERATIONS: %v; using default 10.", err)
        maxIter = 10
    }
    efConstruction, _ := strconv.Atoi(getEnv("CHROMA_EF_CONSTRUCTION", "100"))
    efSearch, _ := strconv.Atoi(getEnv("CHROMA_EF_SEARCH", "100"))
    maxNeighbors, _ := strconv.Atoi(getEnv("CHROMA_MAX_NEIGHBORS", "16"))

    return &Config{
        LLMProvider:          getEnv("LLM_PROVIDER", "openai"),
        OllamaHost:           getEnv("OLLAMA_HOST", "http://localhost:11434"),
        OllamaModel:          getEnv("OLLAMA_MODEL", "llama3.2:3b"),
        OllamaEmbedModel:     getEnv("OLLAMA_EMBED_MODEL", "llama3.2:3b"),
        AzureOpenAIEndpoint:  getEnv("AZURE_OPENAI_ENDPOINT", ""),
        AzureOpenAIAPIKey:    getEnv("AZURE_OPENAI_API_KEY", ""),
        AzureDeployment:      getEnv("AZURE_OPENAI_DEPLOYMENT_NAME", ""),
        OpenAIAPIKey:         getEnv("OPENAI_API_KEY", ""),
        AgentMaxIterations:   maxIter,
        SplitterProvider:     getEnv("SPLITTER_PROVIDER", "recursive"),
        VectorStoreProvider:  getEnv("VECTOR_STORE_PROVIDER", "faiss"),
        ChromaURL:            getEnv("CHROMA_URL", "http://localhost:8000"),
        ChromaSpace:          getEnv("CHROMA_SPACE", "cosine"),
        ChromaCollection:     getEnv("CHROMA_COLLECTION", "GogurtCol"),
        ChromaTenant:         getEnv("CHROMA_TENANT", "joe"),
        ChromaDatabase:       getEnv("CHROMA_DATABASE", "GogurtDB"),
        ChromaEFConstruction: efConstruction,
        ChromaEFSearch:       efSearch,
        ChromaMaxNeighbors:   maxNeighbors,
    }
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
