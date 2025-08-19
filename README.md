
# Gogurt

Gogurt is a Go-based AI agent framework that allows you to create AI agents that can use tools to perform tasks.

## Features

-   **Multiple LLM Providers**: Supports OpenAI, Azure OpenAI, and Ollama.
-   **Tool Use**: Easily create and add tools for the agent to use.
-   **Extensible**: Designed to be easily extended with new LLMs and tools.

## Getting Started

### Prerequisites

-   Go 1.22 or later
-   An API key for your chosen LLM provider (OpenAI, Azure) or a running Ollama instance.

### Installation

1.  Clone the repository:
    ```bash
    git clone [https://github.com/89jobrien/gogurt.git](https://github.com/89jobrien/gogurt.git)
    cd gogurt
    ```

2.  Install dependencies:
    ```bash
    go mod tidy
    ```

3.  Create a `.env` file by copying the example:
    ```bash
    cp .env.example .env
    ```

4.  Edit the `.env` file to add your API keys and configure the LLM provider.

### Usage

Run the application from your terminal:

```bash
go run main.go
```

## Project Structure

```ascii
/gogurt
|-- /agent
|   |-- agent.go
|-- /config
|   |-- config.go
|-- /llm
|   |-- /azure
|   |   |-- azure.go
|   |-- /ollama
|   |   |-- ollama.go
|   |-- /openai
|   |   |-- openai.go
|-- /tools
|   |-- tools.go
|   |-- weather.go
|-- /types
|   |-- types.go
|-- .env.example
|-- .gitignore
|-- go.mod
|-- go.sum
|-- main.go
|-- README.md
```
