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

The application can load documents in two ways:

**Option 1:** Load all files from the `docs/` directory (this is the default)
**Note:** If you have a `docs/` directory in your project root, you can run the application without any arguments to load all supported files (`.txt`, `.pdf`) from it.

```bash
# Make sure you have a docs/ directory with files in it
go run main.go
```

**Option 2:** Load a specific file or directory
**Note:** You can also provide a path to a specific file or directory as a command-line argument.

```bash
# Load a single PDF
go run main.go path/to/your/document.pdf

# Load all files from a different directory
go run main.go path/to/another/directory/
```

### Example Sessions

**Using `docs.txt`:**

```plaintext
$ go run main.go
{"time":"2025-08-19T19:48:34.394173-04:00","level":"INFO","msg":"Using Ollama LLM"}
{"time":"2025-08-19T19:48:34.394322-04:00","level":"INFO","msg":"Setting up RAG pipeline..."}
{"time":"2025-08-19T19:48:35.07571-04:00","level":"INFO","msg":"RAG pipeline setup complete."}
{"time":"2025-08-19T19:48:35.075724-04:00","level":"INFO","msg":"Chat session started. Type 'exit' to end."}
You: what is gogurt?
AI: According to the context, Gogurt is a modular framework that allows you to create AI agents.
You: who develops it?
AI: According to the context, Gogurt is being developed by an engineer named Joe.
You: what is the file structure?
AI: According to the context, the file structure of Gogurt is as follows:

/gogurt
|-- /agent
| |-- agent.go
|-- /config
| |-- config.go
|-
| |-- config.go
|-- /documentloaders
| |-- /text
| |-- text.go
|-- /embeddings
| |-- /ol

This appears to be a hierarchical directory structure, with nested directories for different components of the Gogurt framework.
You: exit
{"time":"2025-08-19T19:50:13.055-04:00","level":"INFO","msg":"Ending chat session."}
```

**Using `docs/` directory:**

```plaintext
{"time":"2025-08-19T20:40:43.291317-04:00","level":"INFO","msg":"No document path provided, loading from default 'docs/' directory"}
{"time":"2025-08-19T20:40:43.291952-04:00","level":"INFO","msg":"Using Ollama LLM"}
{"time":"2025-08-19T20:40:43.291987-04:00","level":"INFO","msg":"Setting up RAG pipeline..."}
{"time":"2025-08-19T20:40:45.940107-04:00","level":"INFO","msg":"RAG pipeline setup complete.","documents_loaded":1,"chunks_created":142}
{"time":"2025-08-19T20:40:45.940156-04:00","level":"INFO","msg":"Chat session started. Type 'exit' to end."}
You: who is Joseph O'Brien?
AI: Based on the provided context, Joseph O'Brien is a developer (presumably a software developer with a presence on Toptal.)
You: exit
{"time":"2025-08-19T20:41:31.415392-04:00","level":"INFO","msg":"Ending chat session."}
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
