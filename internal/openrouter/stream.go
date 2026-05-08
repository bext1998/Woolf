package openrouter

type StreamEvent struct {
	Content string
	Done    bool
	Usage   *Usage
	Error   error
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}
