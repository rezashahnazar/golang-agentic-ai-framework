package provider_test

import (
	"fmt"
	"os"
	"testing"

	"agentic-ai-framework/internal/config"
	"agentic-ai-framework/internal/provider"
	"agentic-ai-framework/internal/runtime"
	"agentic-ai-framework/internal/types"
)

func TestNewOpenAIChatCompletionsProvider(t *testing.T) {
	testConfig := `openai:
  api_key: "test-key"
  base_url: "https://api.openai.com/v1"
`
	err := os.WriteFile("test_config.yaml", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config.yaml")

	p := provider.NewOpenAIChatCompletionsProvider("test_config.yaml")
	if p.Name() != "OpenAI Chat Completions" {
		t.Errorf("expected provider name 'OpenAI Chat Completions', got '%s'", p.Name())
	}
	models := p.AvailableModels()
	if len(models) == 0 {
		t.Error("expected at least one available model")
	}
	providerConfig := p.Config()
	if providerConfig["api_key"] == "" {
		t.Error("expected config to contain api_key")
	}
	if providerConfig["base_url"] == "" {
		t.Error("expected config to contain base_url")
	}
	testYamlConfig := config.LoadConfig("test_config.yaml")
	if providerConfig["api_key"] != testYamlConfig.OpenAI.APIKey {
		t.Error("provider api_key should match test config api_key")
	}
	if providerConfig["base_url"] != testYamlConfig.OpenAI.BaseURL && testYamlConfig.OpenAI.BaseURL != "" {
		t.Error("provider base_url should match test config base_url when provided")
	}
}

func TestProviderAvailableModels(t *testing.T) {
	testConfig := `openai:
  api_key: "test-key"
  base_url: "https://api.openai.com/v1"
`
	err := os.WriteFile("test_config.yaml", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config.yaml")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	p := provider.NewOpenAIChatCompletionsProvider("test_config.yaml")
	models := p.AvailableModels()
	expectedModels := []string{"gpt-4.1", "gpt-5"}
	if len(models) != len(expectedModels) {
		t.Errorf("expected %d models, got %d", len(expectedModels), len(models))
	}
	for i, expectedName := range expectedModels {
		if models[i].Name() != expectedName {
			t.Errorf("expected model %d to be '%s', got '%s'", i, expectedName, models[i].Name())
		}
	}
}

func TestProviderImplementsInterface(t *testing.T) {
	testConfig := `openai:
  api_key: "test-key"
  base_url: "https://api.openai.com/v1"
`
	err := os.WriteFile("test_config.yaml", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config.yaml")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	p := provider.NewOpenAIChatCompletionsProvider("test_config.yaml")
	var _ provider.Provider = p
}

func TestGenerateTextIntegration(t *testing.T) {
	if _, err := os.Stat("../../config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: config.yaml not found")
	}
	p := provider.NewOpenAIChatCompletionsProvider("../../config.yaml")
	var result types.GenerateTextResult
	var panicked bool
	var panicErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
				if err, ok := r.(error); ok {
					panicErr = err
				} else {
					panicErr = fmt.Errorf("%v", r)
				}
			}
		}()
		result = runtime.GenerateText(p, "Say 'test' and nothing else", "gpt-4.1", map[string]any{
			"temperature": 0.0,
			"top_p":       0.9,
		})
	}()
	if panicked {
		t.Fatalf("GenerateText panicked with error: %v", panicErr)
	}
	if result.TextContent() == "" {
		t.Error("expected non-empty text content")
	}
}

func TestProviderValidatesModelName(t *testing.T) {
	testConfig := `openai:
  api_key: "test-key"
  base_url: "https://api.openai.com/v1"
`
	err := os.WriteFile("test_config.yaml", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config.yaml")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when using invalid model name")
		}
		if _, ok := r.(string); !ok {
			t.Errorf("expected panic message to be a string, got %T", r)
		}
	}()
	prv := provider.NewOpenAIChatCompletionsProvider("test_config.yaml")
	prv.GenerateText("test", "invalid-model-name", map[string]any{})
}

func TestProviderAvailableRequestParameters(t *testing.T) {
	testConfig := `openai:
  api_key: "test-key"
  base_url: "https://api.openai.com/v1"
`
	err := os.WriteFile("test_config.yaml", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config.yaml")

	p := provider.NewOpenAIChatCompletionsProvider("test_config.yaml")
	gpt41Params := p.AvailableRequestParameters("gpt-4.1")
	expectedGPT41Params := []string{"temperature", "top_p"}
	if len(gpt41Params) != len(expectedGPT41Params) {
		t.Errorf("expected gpt-4.1 to have %d parameters, got %d", len(expectedGPT41Params), len(gpt41Params))
	}
	hasTemperature := false
	hasTopP := false
	for _, param := range gpt41Params {
		if param == "temperature" {
			hasTemperature = true
		}
		if param == "top_p" {
			hasTopP = true
		}
	}
	if !hasTemperature {
		t.Error("expected gpt-4.1 to support temperature parameter")
	}
	if !hasTopP {
		t.Error("expected gpt-4.1 to support top_p parameter")
	}
	gpt5Params := p.AvailableRequestParameters("gpt-5")
	if len(gpt5Params) != 0 {
		t.Errorf("expected gpt-5 to have no parameters, got %v", gpt5Params)
	}
}

func TestProviderValidatesRequestParameters(t *testing.T) {
	testConfig := `openai:
  api_key: "test-key"
  base_url: "https://api.openai.com/v1"
`
	err := os.WriteFile("test_config.yaml", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config.yaml")
	t.Run("gpt-5 rejects any parameter", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic when using any parameter for gpt-5")
			}
			errMsg, ok := r.(string)
			if !ok {
				t.Errorf("expected panic message to be a string, got %T", r)
			}
			if errMsg == "" {
				t.Error("expected non-empty panic message")
			}
		}()
		prv := provider.NewOpenAIChatCompletionsProvider("test_config.yaml")
		prv.GenerateText("test", "gpt-5", map[string]any{
			"temperature": 0.7,
		})
	})
	t.Run("gpt-4.1 rejects invalid parameters", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic when using invalid parameter for gpt-4.1")
			}
			errMsg, ok := r.(string)
			if !ok {
				t.Errorf("expected panic message to be a string, got %T", r)
			}
			if errMsg == "" {
				t.Error("expected non-empty panic message")
			}
		}()
		prv := provider.NewOpenAIChatCompletionsProvider("test_config.yaml")
		prv.GenerateText("test", "gpt-4.1", map[string]any{
			"max_tokens": 100,
		})
	})
	t.Run("gpt-4.1 accepts valid parameters", func(t *testing.T) {
		prv := provider.NewOpenAIChatCompletionsProvider("test_config.yaml")
		_, err := prv.GenerateText("test", "gpt-4.1", map[string]any{
			"temperature": 0.7,
			"top_p":       0.9,
		})
		if err == nil {
			t.Error("expected error due to invalid API key in test config")
		}
	})
}

func TestProviderGenerateTextDirectly(t *testing.T) {
	if _, err := os.Stat("../../config.yaml"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: config.yaml not found")
	}
	p := provider.NewOpenAIChatCompletionsProvider("../../config.yaml")
	result, err := p.GenerateText("Hello", "gpt-4.1", map[string]any{
		"temperature": 0.7,
		"top_p":       0.9,
	})
	if err != nil {
		t.Fatalf("Provider.GenerateText returned error: %v", err)
	}
	if result.TextContent() == "" {
		t.Error("expected non-empty text content")
	}
	if result.Usage().TotalTokens() == 0 {
		t.Error("expected total tokens to be greater than 0")
	}
}

func TestProviderGenerateTextWithInvalidAPIKey(t *testing.T) {
	invalidConfig := `openai:
  api_key: "invalid-key"
  base_url: "https://api.openai.com/v1"
`
	err := os.WriteFile("test_config_invalid.yaml", []byte(invalidConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config_invalid.yaml")
	p := provider.NewOpenAIChatCompletionsProvider("test_config_invalid.yaml")
	_, err = p.GenerateText("Hello", "gpt-4.1", map[string]any{})
	if err == nil {
		t.Error("expected error with invalid API key")
	}
}

func TestProviderGenerateTextHandlesEmptyChoices(t *testing.T) {
	invalidConfig := `openai:
  api_key: "test-key"
  base_url: "https://invalid-url-that-will-fail.com/v1"
`
	err := os.WriteFile("test_config_empty_choices.yaml", []byte(invalidConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config_empty_choices.yaml")
	p := provider.NewOpenAIChatCompletionsProvider("test_config_empty_choices.yaml")
	_, err = p.GenerateText("test", "gpt-4.1", map[string]any{})
	if err == nil {
		t.Error("expected error when API call fails or returns empty response")
	}
	if err != nil && err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestGenerateTextHelperPropagatesErrors(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected GenerateText to panic when provider returns error")
		}
		errMsg, ok := r.(error)
		if !ok {
			t.Fatalf("expected panic to be error, got %T", r)
		}
		if errMsg.Error() == "" {
			t.Error("expected non-empty error message")
		}
	}()
	invalidConfig := `openai:
  api_key: "invalid-key"
  base_url: "https://api.openai.com/v1"
`
	err := os.WriteFile("test_config_error.yaml", []byte(invalidConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config_error.yaml")
	p := provider.NewOpenAIChatCompletionsProvider("test_config_error.yaml")
	runtime.GenerateText(p, "test", "gpt-4.1", map[string]any{})
}
