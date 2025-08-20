# Gogurt

Gogurt is a modular framework that allows you to create AI agents.

## Features

- **Multiple LLM Providers**: Supports OpenAI, Azure OpenAI, and Ollama.
- **Pluggable Vector Stores**: Choose between a simple in-memory vector store or a persistent ChromaDB instance.
- **Retrieval-Augmented Generation (RAG)**: Ingest documents (`.txt`, `.pdf`, `.md`) and use a configurable text splitter (`recursive`, `markdown`, `character`) to optimize context retrieval.
- **Tool Use**: Easily create and add tools for the agent to use.
- **Extensible**: Designed to be easily extended with new LLMs and tools.

## Getting Started

### Prerequisites

- **Go**: Version 1.25 or later
- **Ollama**: For running local LLMs and embedding models
- **Poppler**: Required for .pdf document parsing
  - **macOS**: brew install poppler
  - **Debian/Ubuntu**: sudo apt-get install poppler-utils
  - **Windows**: I don't know. Use a better OS?
- **Docker**: Required if you plan to use the ChromaDB vector store.

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

3.  **Create Documents Directory**:
    Create a `docs/` directory in the project root and add any `.txt`, `.pdf`, or `.md` files you want the agent to learn from.

        ```bash
        mkdir docs
        cp path/to/your/file.pdf docs/
        ```

4.  **Configure Environment**:
    Copy the example environment file and edit it to match your setup.

    ```bash
    cp .env.example .env
    ```

    See the **Configuration** section below for a detailed explanation of each variable.

---

## Configuration

All configuration is managed through the `.env` file.

| Variable                | Default                 | Description                                                              |
| ----------------------- | ----------------------- | ------------------------------------------------------------------------ |
| `LLM_PROVIDER`          | `ollama`                | The chat model provider. Options: `ollama`, `openai`, `azure`.           |
| `OLLAMA_MODEL`          | `llama3.2:3b`           | The Ollama model to use for chat generation.                             |
| `OLLAMA_EMBED_MODEL`    | `llama3.2:3b`           | The Ollama model to use for creating document embeddings.                |
| `AGENT_MAX_ITERATIONS`  | `10`                    | The maximum number of steps the agent can take to answer a query.        |
| `SPLITTER_PROVIDER`     | `recursive`             | The text splitter to use. Options: `recursive`, `markdown`, `character`. |
| `VECTOR_STORE_PROVIDER` | `simple`                | The vector store to use. Options: `simple` (in-memory), `chroma`.        |
| `CHROMA_URL`            | `http://localhost:8000` | The URL for your running ChromaDB instance.                              |
| `OPENAI_API_KEY`        | `your-api-key`          | Your API key for OpenAI.                                                 |
| `AZURE_OPENAI_...`      | `your-key`              | Your credentials for Azure OpenAI services.                              |

---

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

### Using with ChromaDB

If you set `VECTOR_STORE_PROVIDER=chroma` in your `.env` file, you must first start a ChromaDB instance using Docker:

```bash
docker run -p 8000:8000 chromadb/chroma
```

Then, run the `gogurt` application as usual. Your document embeddings will be stored persistently in the ChromaDB container.

## Example Sessions

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
{"time":"2025-08-20T00:40:55.758264-04:00","level":"INFO","msg":"No document path provided, loading from default 'docs/' directory"}
{"time":"2025-08-20T00:40:55.758811-04:00","level":"INFO","msg":"Setting up RAG pipeline..."}
{"time":"2025-08-20T00:40:55.75885-04:00","level":"INFO","msg":"Using Ollama for LLM"}
{"time":"2025-08-20T00:40:55.797159-04:00","level":"INFO","msg":"Using recursive text splitter"}
{"time":"2025-08-20T00:40:55.79732-04:00","level":"INFO","msg":"Using simple in-memory vector store"}
{"time":"2025-08-20T00:40:58.37175-04:00","level":"INFO","msg":"RAG pipeline setup complete.","documents_loaded":3,"chunks_created":52}
{"time":"2025-08-20T00:40:58.371791-04:00","level":"INFO","msg":"Chat session started. Type 'exit' to end."}
You: what is gogurt?
AI: According to the context, Gogurt is a modular framework that allows you to create AI agents.
You: who is Joseph O'Brien?
AI: Based on the context, Joseph O'Brien is a developer with a presence on Toptal.
You: exit
{"time":"2025-08-20T00:41:46.183651-04:00","level":"INFO","msg":"Ending chat session."}
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
