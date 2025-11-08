package provider

import (
	"testing"
)

func TestProviderInterfaceCompliance(t *testing.T) {
	var _ Provider = &OpenAIChatCompletionsProvider{}
}

func TestModelInterfaceCompliance(t *testing.T) {
	var _ Model = &OpenAIModel{}
}
