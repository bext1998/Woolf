package config

func Default() Config {
	return Config{
		API: APIConfig{
			BaseURL:    "https://openrouter.ai/api/v1",
			TimeoutSec: 120,
			MaxRetries: 3,
		},
		Defaults: DefaultsConfig{
			MaxRounds:         3,
			AutoSave:          true,
			Language:          "zh-TW",
			DefaultPreset:     "editorial",
			SummarizerEnabled: false,
			SummarizerModel:   "openai/gpt-4o-mini",
		},
		TUI: TUIConfig{
			Theme:          "dark",
			ShowTokenCount: true,
			ShowCost:       true,
			StreamSpeed:    "realtime",
			Editor:         "vim",
		},
		Context: ContextConfig{
			MaxWindowRatio:  0.70,
			SummaryModel:    "openai/gpt-4o-mini",
			SummaryMaxToken: 200,
		},
		Budget: BudgetConfig{
			SessionLimitUSD: 0,
			WarnThreshold:   0.80,
		},
	}
}
