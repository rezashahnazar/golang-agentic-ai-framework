package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testConfig := `openai:
  api_key: "test-key"
  base_url: "https://api.openai.com/v1"
`
	err := os.WriteFile("test_config.yaml", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}
	defer os.Remove("test_config.yaml")

	cfg := LoadConfig("test_config.yaml")

	if cfg.OpenAI.APIKey == "" {
		t.Error("expected openai.api_key to be loaded from config.yaml")
	}
	if cfg.OpenAI.APIKey != "test-key" {
		t.Errorf("expected api_key 'test-key', got '%s'", cfg.OpenAI.APIKey)
	}
	if cfg.OpenAI.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("expected base_url 'https://api.openai.com/v1', got '%s'", cfg.OpenAI.BaseURL)
	}
}
