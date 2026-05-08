package openrouter

type ModelInfo struct {
	ID                     string
	Name                   string
	ContextLength          int
	PromptCostPerToken     float64
	CompletionCostPerToken float64
}
