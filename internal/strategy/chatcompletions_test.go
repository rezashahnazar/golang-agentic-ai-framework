package strategy

import (
	"net/http"
	"testing"
)

func TestBuildChatCompletionsRequestBody(t *testing.T) {
	req := ChatCompletionsRequest{
		Model: "gpt-4",
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a helpful assistant"},
			{Role: "user", Content: "Hello"},
		},
		RequestParams: map[string]any{
			"temperature": 0.7,
			"max_tokens":  100,
		},
	}

	body := BuildChatCompletionsRequestBody(req)

	if body["model"] != "gpt-4" {
		t.Errorf("expected model 'gpt-4', got %v", body["model"])
	}

	messages, ok := body["messages"].([]map[string]string)
	if !ok {
		t.Fatal("expected messages to be []map[string]string")
	}

	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
	}

	if messages[0]["role"] != "system" || messages[0]["content"] != "You are a helpful assistant" {
		t.Error("first message incorrect")
	}

	if messages[1]["role"] != "user" || messages[1]["content"] != "Hello" {
		t.Error("second message incorrect")
	}

	if body["temperature"] != 0.7 {
		t.Errorf("expected temperature 0.7, got %v", body["temperature"])
	}

	if body["max_tokens"] != 100 {
		t.Errorf("expected max_tokens 100, got %v", body["max_tokens"])
	}
}

func TestParseChatCompletionsResponse(t *testing.T) {
	t.Run("successful parsing", func(t *testing.T) {
		response := ChatCompletionsResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "Hello world"}},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		}

		result, err := ParseChatCompletionsResponse(response, http.StatusOK)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.TextContent() != "Hello world" {
			t.Errorf("expected 'Hello world', got '%s'", result.TextContent())
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

	t.Run("empty choices", func(t *testing.T) {
		response := ChatCompletionsResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{},
		}

		_, err := ParseChatCompletionsResponse(response, http.StatusOK)
		if err == nil {
			t.Fatal("expected error for empty choices")
		}
		if err.Error() != "no choices in API response" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("API error", func(t *testing.T) {
		response := ChatCompletionsResponse{
			Error: struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			}{
				Message: "Invalid API key",
				Type:    "authentication_error",
				Code:    "invalid_api_key",
			},
		}

		_, err := ParseChatCompletionsResponse(response, http.StatusUnauthorized)
		if err == nil {
			t.Fatal("expected error for API error")
		}
		expected := "API error (status 401): Invalid API key (type: authentication_error, code: invalid_api_key)"
		if err.Error() != expected {
			t.Errorf("expected '%s', got '%s'", expected, err.Error())
		}
	})

	t.Run("HTTP error", func(t *testing.T) {
		response := ChatCompletionsResponse{}

		_, err := ParseChatCompletionsResponse(response, http.StatusInternalServerError)
		if err == nil {
			t.Fatal("expected error for HTTP error")
		}
		expected := "API request failed with status 500"
		if err.Error() != expected {
			t.Errorf("expected '%s', got '%s'", expected, err.Error())
		}
	})
}
