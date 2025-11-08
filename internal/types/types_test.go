package types

import "testing"

func TestTokenUsageAccessors(t *testing.T) {
	usage := NewTokenUsage(10, 20, 30)

	if usage.PromptTokens() != 10 {
		t.Errorf("expected PromptTokens 10, got %d", usage.PromptTokens())
	}
	if usage.CompletionTokens() != 20 {
		t.Errorf("expected CompletionTokens 20, got %d", usage.CompletionTokens())
	}
	if usage.TotalTokens() != 30 {
		t.Errorf("expected TotalTokens 30, got %d", usage.TotalTokens())
	}
}

func TestGenerateTextResultAccessors(t *testing.T) {
	usage := NewTokenUsage(5, 10, 15)
	result := NewGenerateTextResult("Hello world", usage)

	if result.TextContent() != "Hello world" {
		t.Errorf("expected TextContent 'Hello world', got '%s'", result.TextContent())
	}
	if result.Usage().PromptTokens() != 5 {
		t.Errorf("expected Usage PromptTokens 5, got %d", result.Usage().PromptTokens())
	}
}
