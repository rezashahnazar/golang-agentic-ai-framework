package strategy

import (
	"fmt"
	"net/http"

	"agentic-ai-framework/internal/transport"
	"agentic-ai-framework/internal/types"
)

type ChatCompletionsRequest struct {
	Model            string
	Messages         []ChatMessage
	RequestParams    map[string]any
}

type ChatMessage struct {
	Role    string
	Content string
}

type ChatCompletionsResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

type ChatCompletionsConfig struct {
	BaseURL    string
	Endpoint   string
	APIKey     string
	HTTPClient *http.Client
}

func BuildChatCompletionsRequestBody(req ChatCompletionsRequest) map[string]any {
	messages := make([]map[string]string, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	requestBody := map[string]any{
		"model":    req.Model,
		"messages": messages,
	}

	for key, value := range req.RequestParams {
		requestBody[key] = value
	}

	return requestBody
}

func ExecuteChatCompletionsRequest(config ChatCompletionsConfig, requestBody map[string]any) (ChatCompletionsResponse, int, error) {
	url := config.BaseURL + config.Endpoint
	headers := map[string]string{
		"Authorization": "Bearer " + config.APIKey,
	}

	ctx, cancel := transport.CreateRequestContext(transport.DefaultTimeout)
	defer cancel()

	req, err := transport.CreateJSONRequest(ctx, "POST", url, requestBody, headers)
	if err != nil {
		return ChatCompletionsResponse{}, 0, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := transport.ExecuteRequest(config.HTTPClient, req)
	if err != nil {
		return ChatCompletionsResponse{}, 0, err
	}

	statusCode := resp.StatusCode
	bodyBytes, err := transport.ReadResponseBody(resp)
	if err != nil {
		return ChatCompletionsResponse{}, statusCode, err
	}

	var responseBody ChatCompletionsResponse
	err = transport.DecodeJSONResponse(bodyBytes, &responseBody)
	if err != nil {
		return ChatCompletionsResponse{}, statusCode, fmt.Errorf("failed to decode response: %v", err)
	}

	return responseBody, statusCode, nil
}

func ParseChatCompletionsResponse(response ChatCompletionsResponse, statusCode int) (types.GenerateTextResult, error) {
	if statusCode != http.StatusOK {
		if response.Error.Message != "" {
			return types.GenerateTextResult{}, fmt.Errorf("API error (status %d): %s (type: %s, code: %s)", statusCode, response.Error.Message, response.Error.Type, response.Error.Code)
		}
		return types.GenerateTextResult{}, fmt.Errorf("API request failed with status %d", statusCode)
	}

	if response.Error.Message != "" {
		return types.GenerateTextResult{}, fmt.Errorf("API error: %s (type: %s)", response.Error.Message, response.Error.Type)
	}

	if len(response.Choices) == 0 {
		return types.GenerateTextResult{}, fmt.Errorf("no choices in API response")
	}

	usage := types.NewTokenUsage(
		response.Usage.PromptTokens,
		response.Usage.CompletionTokens,
		response.Usage.TotalTokens,
	)

	result := types.NewGenerateTextResult(
		response.Choices[0].Message.Content,
		usage,
	)

	return result, nil
}
