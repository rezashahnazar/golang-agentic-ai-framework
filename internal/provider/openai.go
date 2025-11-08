package provider

import (
	"fmt"
	"net/http"

	"agentic-ai-framework/internal/config"
	"agentic-ai-framework/internal/strategy"
	"agentic-ai-framework/internal/transport"
	"agentic-ai-framework/internal/types"
)

type OpenAIModel struct {
	name       string
	parameters []string
}

func (m *OpenAIModel) Name() string {
	return m.name
}

func (m *OpenAIModel) AvailableRequestParameters() []string {
	return m.parameters
}

type OpenAIChatCompletionsProvider struct {
	config            map[string]any
	availableModels   []Model
	modelParameters   map[string][]string
	name              string
	apiKey            string
	baseURL           string
	httpClient        *http.Client
}

func NewOpenAIChatCompletionsProvider(configFile string) *OpenAIChatCompletionsProvider {
	cfg := config.LoadConfig(configFile)

	if cfg.OpenAI.APIKey == "" {
		panic("openai.api_key is required in config file")
	}

	baseURL := cfg.OpenAI.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	provider := &OpenAIChatCompletionsProvider{
		name:   "OpenAI Chat Completions",
		apiKey: cfg.OpenAI.APIKey,
		baseURL: baseURL,
		httpClient: transport.NewClient(transport.DefaultTimeout),
		availableModels: []Model{
			&OpenAIModel{name: "gpt-4.1", parameters: []string{"temperature", "top_p"}},
			&OpenAIModel{name: "gpt-5", parameters: []string{}},
		},
		modelParameters: map[string][]string{
			"gpt-4.1": {"temperature", "top_p"},
			"gpt-5":   {},
		},
		config: map[string]any{
			"api_key":  cfg.OpenAI.APIKey,
			"base_url": baseURL,
		},
	}

	return provider
}

func (p *OpenAIChatCompletionsProvider) Name() string {
	return p.name
}

func (p *OpenAIChatCompletionsProvider) AvailableModels() []Model {
	return p.availableModels
}

func (p *OpenAIChatCompletionsProvider) AvailableRequestParameters(modelName string) []string {
	if params, exists := p.modelParameters[modelName]; exists {
		return params
	}
	return []string{}
}

func (p *OpenAIChatCompletionsProvider) GetModel(modelName string) (Model, error) {
	for _, model := range p.availableModels {
		if model.Name() == modelName {
			return model, nil
		}
	}
	return nil, fmt.Errorf("model %s not found in provider %s", modelName, p.Name())
}

func (p *OpenAIChatCompletionsProvider) Config() map[string]any {
	return p.config
}

func (p *OpenAIChatCompletionsProvider) GenerateText(prompt string, modelName string, requestParameters map[string]any) (types.GenerateTextResult, error) {
	if err := ValidateModel(p.availableModels, modelName, p.Name()); err != nil {
		panic(err.Error())
	}

	availableParams := p.AvailableRequestParameters(modelName)
	if err := ValidateRequestParameters(availableParams, requestParameters, modelName); err != nil {
		panic(err.Error())
	}

	chatRequest := strategy.ChatCompletionsRequest{
		Model: modelName,
		Messages: []strategy.ChatMessage{
			{Role: "user", Content: prompt},
		},
		RequestParams: requestParameters,
	}

	requestBody := strategy.BuildChatCompletionsRequestBody(chatRequest)

	cfg := strategy.ChatCompletionsConfig{
		BaseURL:    p.baseURL,
		Endpoint:   "/chat/completions",
		APIKey:     p.apiKey,
		HTTPClient: p.httpClient,
	}

	response, statusCode, err := strategy.ExecuteChatCompletionsRequest(cfg, requestBody)
	if err != nil {
		return types.GenerateTextResult{}, err
	}

	return strategy.ParseChatCompletionsResponse(response, statusCode)
}
