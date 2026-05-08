package session

import "time"

type Status string

const (
	StatusActive    Status = "active"
	StatusPaused    Status = "paused"
	StatusCompleted Status = "completed"
	StatusError     Status = "error"
)

type Session struct {
	SessionID     string            `json:"session_id"`
	Version       string            `json:"version"`
	Title         string            `json:"title,omitempty"`
	Status        Status            `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Source        *Source           `json:"source,omitempty"`
	AgentsConfig  []AgentConfig     `json:"agents_config"`
	Rounds        []Round           `json:"rounds"`
	Interventions []Intervention    `json:"interventions"`
	Summaries     map[string]string `json:"summaries"`
	Totals        Totals            `json:"totals"`
}

type Source struct {
	Type           string `json:"type"`
	Path           string `json:"path,omitempty"`
	Content        string `json:"content,omitempty"`
	ContentHash    string `json:"content_hash,omitempty"`
	ContentPreview string `json:"content_preview,omitempty"`
}

type AgentConfig struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Model       string `json:"model"`
	Stance      string `json:"stance"`
	Order       int    `json:"order"`
	Color       string `json:"color"`
}

type Round struct {
	RoundIndex  int        `json:"round_index"`
	StartedAt   time.Time  `json:"started_at,omitempty"`
	CompletedAt time.Time  `json:"completed_at,omitempty"`
	Responses   []Response `json:"responses"`
}

type Response struct {
	AgentName    string     `json:"agent_name"`
	Model        string     `json:"model"`
	RespondingTo *string    `json:"responding_to"`
	StanceTag    *string    `json:"stance_tag"`
	Content      string     `json:"content"`
	Tokens       TokenUsage `json:"tokens"`
	CostUSD      float64    `json:"cost_usd"`
	Timestamp    time.Time  `json:"timestamp"`
	Status       string     `json:"status"`
}

type TokenUsage struct {
	Prompt     int `json:"prompt"`
	Completion int `json:"completion"`
}

type Intervention struct {
	AfterRound int         `json:"after_round"`
	Type       string      `json:"type"`
	Content    string      `json:"content"`
	FocusRange *FocusRange `json:"focus_range"`
	Timestamp  time.Time   `json:"timestamp"`
}

type FocusRange struct {
	StartLine int `json:"start_line"`
	EndLine   int `json:"end_line"`
}

type Totals struct {
	RoundsCompleted       int     `json:"rounds_completed"`
	TotalTokens           int     `json:"total_tokens"`
	TotalPromptTokens     int     `json:"total_prompt_tokens"`
	TotalCompletionTokens int     `json:"total_completion_tokens"`
	TotalCostUSD          float64 `json:"total_cost_usd"`
}
