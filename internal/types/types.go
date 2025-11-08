package types

type GenerateTextResult struct {
	textContent string
	tokenUsage  TokenUsage
}

func (r *GenerateTextResult) TextContent() string {
	return r.textContent
}

func (r *GenerateTextResult) Usage() TokenUsage {
	return r.tokenUsage
}

type TokenUsage struct {
	promptTokens     int
	completionTokens int
	totalTokens      int
}

func (t TokenUsage) PromptTokens() int {
	return t.promptTokens
}

func (t TokenUsage) CompletionTokens() int {
	return t.completionTokens
}

func (t TokenUsage) TotalTokens() int {
	return t.totalTokens
}

func NewGenerateTextResult(textContent string, tokenUsage TokenUsage) GenerateTextResult {
	return GenerateTextResult{
		textContent: textContent,
		tokenUsage:  tokenUsage,
	}
}

func NewTokenUsage(promptTokens, completionTokens, totalTokens int) TokenUsage {
	return TokenUsage{
		promptTokens:     promptTokens,
		completionTokens: completionTokens,
		totalTokens:      totalTokens,
	}
}
