# Golang Agentic AI Framework

[![CI](https://github.com/rezashahnazar/golang-agentic-ai-framework/actions/workflows/ci.yml/badge.svg)](https://github.com/rezashahnazar/golang-agentic-ai-framework/actions/workflows/ci.yml)

**Status**: Under development.

A scalable Golang framework for creation, execution, and management of Agentic AI systems.

## Current Implementation

### Core Architecture

The framework provides a standardized interface for AI providers with model-specific parameter validation:

- **Provider Interface**: Unified interface for all AI providers
- **Model-Specific Parameters**: Each model defines its own available request parameters
- **YAML Configuration**: Provider credentials configured via `config.yaml`
- **Error Handling**: Comprehensive error handling with detailed API error messages

### Features

- **OpenAI Provider**: Full implementation of OpenAI Chat Completions API
- **Model Support**:
  - `gpt-4.1`: Supports `temperature` and `top_p` parameters
  - `gpt-5`: No parameters supported
- **Parameter Validation**: Automatic validation of request parameters against model capabilities
- **Token Usage Tracking**: Tracks prompt, completion, and total tokens

## Setup

1. Copy the example config file:

   ```bash
   cp config.yaml.example config.yaml
   ```

2. Edit `config.yaml` with your OpenAI credentials:

   ```yaml
   openai:
     api_key: "your-api-key-here"
     base_url: "https://api.openai.com/v1"
   ```

3. Build the project:
   ```bash
   go build ./...
   ```

## Usage

### Run the example

```bash
go run ./examples/basic
```

### Example Code

```go
package main

import (
    "fmt"
    runtime "agentic-ai-framework/internal/runtime"
    "agentic-ai-framework/internal/provider"
)

func main() {
    p := provider.NewOpenAIChatCompletionsProvider("config.yaml")

    result := runtime.GenerateText(
        p,
        "Hello! How are you?",
        "gpt-4.1",
        map[string]any{
            "temperature": 0.7,
            "top_p": 0.9,
        },
    )

    fmt.Println(result.TextContent())
    fmt.Printf("Tokens used: %d\n", result.Usage().TotalTokens())
}
```

### Provider Interface

```go
type Model interface {
    Name() string
    AvailableRequestParameters() []string
}

type Provider interface {
    Name() string
    AvailableModels() []Model
    GetModel(modelName string) (Model, error)
    AvailableRequestParameters(modelName string) []string
    Config() map[string]any
    GenerateText(prompt string, modelName string, requestParameters map[string]any) (GenerateTextResult, error)
}
```

### Working with Models

You can work with models in two ways:

**Method 1: Using Model objects (type-safe)**

```go
p := provider.NewOpenAIChatCompletionsProvider("config.yaml")
model, err := p.GetModel("gpt-4.1")
if err != nil { panic(err) }
params := model.AvailableRequestParameters()
result := runtime.GenerateText(p, "Hello", model.Name(), map[string]any{ "temperature": 0.7 })
```

**Method 2: Using string names (convenience)**

```go
p := provider.NewOpenAIChatCompletionsProvider("config.yaml")
params := p.AvailableRequestParameters("gpt-4.1")
result := runtime.GenerateText(p, "Hello", "gpt-4.1", map[string]any{ "temperature": 0.7 })
```

## Testing

Run all tests:

```bash
go test ./... -v
```

Run specific test:

```bash
go test -v ./internal/provider -run TestProviderValidatesRequestParameters
```

### CI/CD

This project uses GitHub Actions for continuous integration:

- **Automated Testing**: All tests run on push/PR to main branch
- **Demo Execution**: Example program runs with real API credentials
- **Environment Setup**: Uses `OPENAI_API_KEY` secret and `OPENAI_BASE_URL` variable

Workflow: `.github/workflows/ci.yml`

## Project Structure

```
.
├── .github/
│   └── workflows/
│       └── ci.yml             # GitHub Actions CI/CD
├── examples/
│   └── basic/
│       └── main.go            # Example program
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── provider/
│   │   ├── interfaces.go
│   │   ├── interfaces_test.go
│   │   ├── openai.go
│   │   ├── openai_test.go
│   │   ├── validation.go
│   │   └── validation_test.go
│   ├── runtime/
│   │   ├── runtime.go
│   │   └── runtime_test.go
│   ├── strategy/
│   │   ├── chatcompletions.go
│   │   └── chatcompletions_test.go
│   ├── transport/
│   │   ├── client.go
│   │   └── client_test.go
│   └── types/
│       ├── types.go
│       └── types_test.go
├── config.yaml
├── config.yaml.example
├── go.mod
├── go.sum
└── README.md
```

## Development Status

**Phase 1 - Simple Core**: Complete

- Basic OpenAI provider implementation
- Model-specific parameter validation
- YAML configuration support
- Comprehensive test coverage (19 tests, 100% package coverage)
- Error handling and validation
- Clean architecture with separated concerns (provider, strategy, transport, runtime)

## Next Steps

- Additional AI providers (Anthropic, Ollama, etc.)
- Streaming response support
- Tool/function calling
- Conversation memory management
- Workflow patterns

---
