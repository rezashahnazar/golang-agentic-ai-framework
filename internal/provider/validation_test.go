package provider

import (
	"testing"
)

func TestValidateModel(t *testing.T) {
	t.Run("valid model", func(t *testing.T) {
		models := []Model{
			&mockModel{name: "gpt-4"},
			&mockModel{name: "gpt-3.5"},
		}

		err := ValidateModel(models, "gpt-4", "TestProvider")
		if err != nil {
			t.Errorf("expected no error for valid model, got %v", err)
		}
	})

	t.Run("invalid model", func(t *testing.T) {
		models := []Model{
			&mockModel{name: "gpt-4"},
			&mockModel{name: "gpt-3.5"},
		}

		err := ValidateModel(models, "invalid-model", "TestProvider")
		if err == nil {
			t.Fatal("expected error for invalid model")
		}

		expected := "model invalid-model is not available in provider TestProvider"
		if err.Error() != expected {
			t.Errorf("expected '%s', got '%s'", expected, err.Error())
		}
	})
}

func TestValidateRequestParameters(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		availableParams := []string{"temperature", "top_p", "max_tokens"}
		requestParams := map[string]any{
			"temperature": 0.7,
			"top_p":       0.9,
		}

		err := ValidateRequestParameters(availableParams, requestParams, "gpt-4")
		if err != nil {
			t.Errorf("expected no error for valid parameters, got %v", err)
		}
	})

	t.Run("invalid parameter", func(t *testing.T) {
		availableParams := []string{"temperature", "top_p"}
		requestParams := map[string]any{
			"temperature": 0.7,
			"max_tokens":  100,
		}

		err := ValidateRequestParameters(availableParams, requestParams, "gpt-4")
		if err == nil {
			t.Fatal("expected error for invalid parameter")
		}

		expected := "request parameter 'max_tokens' is not available for model gpt-4. Available parameters: [temperature top_p]"
		if err.Error() != expected {
			t.Errorf("expected '%s', got '%s'", expected, err.Error())
		}
	})

	t.Run("empty request parameters", func(t *testing.T) {
		availableParams := []string{"temperature", "top_p"}
		requestParams := map[string]any{}

		err := ValidateRequestParameters(availableParams, requestParams, "gpt-4")
		if err != nil {
			t.Errorf("expected no error for empty parameters, got %v", err)
		}
	})
}

type mockModel struct {
	name string
}

func (m *mockModel) Name() string {
	return m.name
}

func (m *mockModel) AvailableRequestParameters() []string {
	return []string{"temperature", "top_p"}
}
