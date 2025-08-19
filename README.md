# Gogurt

Gogurt is a modular framework that allows you to create AI agents.

## Features

- **Multiple LLM Providers**: Supports OpenAI, Azure OpenAI, and Ollama.
- **Tool Use**: Easily create and add tools for the agent to use.
- **Retrieval-Augmented Generation (RAG)**: Ask questions about your own documents. The agent can read a local text file to provide context-aware answers.
- **Extensible**: Designed to be easily extended with new LLMs and tools.

## Getting Started

### Prerequisites

- Go 1.22 or later
- For local inference: Ensure you have a running Ollama instance with a chat model and a model that supports embeddings

### Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/89jobrien/gogurt.git
    cd gogurt
    ```

2.  Install dependencies:

    ```bash
    go mod tidy
    ```

3.  **Create a Knowledge Base**
    Create a file named `docs.txt` in the root of the project. This will be the document your agent reads from. Add some text to it, for example:

    ```text
    Gogurt is a powerful, Go-based AI agent framework.
    It was created in 2025 and is designed for building modular and extensible AI applications.
    The framework supports multiple Large Language Models, including providers like Ollama, OpenAI, and Azure.
    ```

4.  Create a `.env` file by copying the example:

    ```bash
    cp .env.example .env
    ```

5.  Edit the `.env` file to configure the LLM provider. For the RAG pipeline, ensure `LLM_PROVIDER`, `OLLAMA_EMBED_MODEL` and `OLLAMA_MODEL` are set correctly.

    ```env
    LLM_PROVIDER=ollama
    OLLAMA_MODEL=llama3.2:3b
    OLLAMA_EMBED_MODEL=nomic-embed-text:latest
    ```

## Usage

Run the application from your terminal. It will start an interactive chat session.

```bash
go run main.go
```

The application will first set up the RAG pipeline by loading and processing `docs.txt`. Once ready, you can start asking questions.

### Example Session

```
$ go run main.go
{"time":"2025-08-19T19:37:00.000-04:00","level":"INFO","msg":"Setting up RAG pipeline..."}
{"time":"2025-08-19T19:37:00.000-04:00","level":"INFO","msg":"RAG pipeline setup complete."}
{"time":"2025-08-19T19:37:00.000-04:00","level":"INFO","msg":"Chat session started. Type 'exit' to end."}
You: What is Gogurt?
AI: Gogurt is a powerful, Go-based AI agent framework created in 2025. It's designed for building modular AI applications and supports multiple Large Language Models.
You: exit
{"time":"2025-08-19T19:38:00.000-04:00","level":"INFO","msg":"Ending chat session."}
```

## Project Structure

```ascii
/gogurt
|-- /agent
|   |-- agent.go
|-- /config
|   |-- config.go
|-- /documentloaders
|   |-- /text
|       |-- text.go
|-- /embeddings
|   |-- /ollama
|   |   |-- ollama.go
|   |-- embeddings.go
|-- /llm
|   |-- /azure
|   |-- /ollama
|   |-- /openai
|-- /textsplitter
|   |-- /character
|       |-- character.go
|-- /tools
|   |-- tools.go
|   |-- weather.go
|-- /types
|   |-- types.go
|-- /vectorstores
|   |-- /simple
|       |-- simple.go
|-- .env.example
|-- docs.txt
|-- go.mod
|-- go.sum
|-- main.go
|-- README.md
```
