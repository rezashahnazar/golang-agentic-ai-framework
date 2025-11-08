package runtime

import (
	"errors"
	"testing"

	"agentic-ai-framework/internal/provider"
	"agentic-ai-framework/internal/types"
)

type mockProvider struct {
	shouldError bool
}

func (m *mockProvider) Name() string {
	return "MockProvider"
}

func (m *mockProvider) AvailableModels() []provider.Model {
	return nil
}

func (m *mockProvider) GetModel(modelName string) (provider.Model, error) {
	return nil, nil
}

func (m *mockProvider) AvailableRequestParameters(modelName string) []string {
	return nil
}

func (m *mockProvider) Config() map[string]any {
	return nil
}

func (m *mockProvider) GenerateText(prompt string, modelName string, requestParameters map[string]any) (types.GenerateTextResult, error) {
	if m.shouldError {
		return types.GenerateTextResult{}, errors.New("mock provider error")
	}
	return types.NewGenerateTextResult("Mock response", types.NewTokenUsage(10, 20, 30)), nil
}

func TestGenerateText(t *testing.T) {
	t.Run("successful generation", func(t *testing.T) {
		provider := &mockProvider{shouldError: false}

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("unexpected panic: %v", r)
			}
		}()

		result := GenerateText(provider, "test prompt", "gpt-4", map[string]any{})

		if result.TextContent() != "Mock response" {
			t.Errorf("expected 'Mock response', got '%s'", result.TextContent())
		}

		if result.Usage().PromptTokens() != 10 {
			t.Errorf("expected 10 prompt tokens, got %d", result.Usage().PromptTokens())
		}

		if result.Usage().CompletionTokens() != 20 {
			t.Errorf("expected 20 completion tokens, got %d", result.Usage().CompletionTokens())
		}

		if result.Usage().TotalTokens() != 30 {
			t.Errorf("expected 30 total tokens, got %d", result.Usage().TotalTokens())
		}
	})

	t.Run("provider error causes panic", func(t *testing.T) {
		provider := &mockProvider{shouldError: true}

		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic when provider returns error")
			}
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected panic to be error, got %T", r)
			}
			if err.Error() != "mock provider error" {
				t.Errorf("expected 'mock provider error', got '%s'", err.Error())
			}
		}()

		GenerateText(provider, "test prompt", "gpt-4", map[string]any{})
	})
}
