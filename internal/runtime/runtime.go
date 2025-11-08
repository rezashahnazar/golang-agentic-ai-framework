package runtime

import (
	"agentic-ai-framework/internal/provider"
	"agentic-ai-framework/internal/types"
)

func GenerateText(p provider.Provider, prompt string, modelName string, requestParameters map[string]any) types.GenerateTextResult {
	response, err := p.GenerateText(prompt, modelName, requestParameters)
	if err != nil {
		panic(err)
	}
	return response
}

