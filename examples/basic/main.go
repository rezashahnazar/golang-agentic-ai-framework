package main

import (
	"fmt"

	"agentic-ai-framework/internal/provider"
	"agentic-ai-framework/internal/runtime"
)

func main() {
	p := provider.NewOpenAIChatCompletionsProvider("config.yaml")

	fmt.Println("Provider:", p.Name())
	fmt.Println("Available Models:")
	models := p.AvailableModels()
	for _, model := range models {
		fmt.Println("  -", model.Name())
	}
	fmt.Println()

	fmt.Println("Generating text with", models[0].Name(), "...")

	result := runtime.GenerateText(p, "Hello! How are you?", models[0].Name(), map[string]any{
		"temperature": 0.7,
		"top_p":       0.9,
	})

	fmt.Println("\nGenerated text:")
	fmt.Println(result.TextContent())
	fmt.Println("\nToken usage:")
	fmt.Printf("  Prompt tokens: %d\n", result.Usage().PromptTokens())
	fmt.Printf("  Completion tokens: %d\n", result.Usage().CompletionTokens())
	fmt.Printf("  Total tokens: %d\n", result.Usage().TotalTokens())
}
