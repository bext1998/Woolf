package agents

import (
	"fmt"
	"strings"
)

type Role struct {
	Name             string   `yaml:"name"`
	DisplayName      string   `yaml:"display_name"`
	Model            string   `yaml:"model"`
	Stance           string   `yaml:"stance"`
	Temperature      float64  `yaml:"temperature"`
	MaxTokens        int      `yaml:"max_tokens"`
	FocusAreas       []string `yaml:"focus_areas"`
	SystemPrompt     string   `yaml:"system_prompt"`
	ResponseTemplate string   `yaml:"response_template"`
	Color            string   `yaml:"color"`
	FallbackModel    string   `yaml:"fallback_model"`
}

func (r Role) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("CFG-003: role name is required")
	}
	if strings.TrimSpace(r.DisplayName) == "" {
		return fmt.Errorf("CFG-003: role %s display_name is required", r.Name)
	}
	if strings.TrimSpace(r.Model) == "" {
		return fmt.Errorf("CFG-003: role %s model is required", r.Name)
	}
	if strings.TrimSpace(r.SystemPrompt) == "" {
		return fmt.Errorf("CFG-003: role %s system_prompt is required", r.Name)
	}
	switch r.Stance {
	case "", "critique", "support", "neutral":
	default:
		return fmt.Errorf("CFG-003: role %s stance must be critique, support, or neutral", r.Name)
	}
	return nil
}
