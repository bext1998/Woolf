package agents

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
