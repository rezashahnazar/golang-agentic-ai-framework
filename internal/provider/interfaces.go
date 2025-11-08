package provider

import (
	"agentic-ai-framework/internal/types"
)

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
	GenerateText(prompt string, modelName string, requestParameters map[string]any) (types.GenerateTextResult, error)
}
