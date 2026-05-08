package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	API      APIConfig      `toml:"api"`
	Defaults DefaultsConfig `toml:"defaults"`
	TUI      TUIConfig      `toml:"tui"`
	Context  ContextConfig  `toml:"context"`
	Budget   BudgetConfig   `toml:"budget"`
	Paths    PathsConfig    `toml:"paths"`
}

type APIConfig struct {
	OpenRouterKey string `toml:"openrouter_key"`
	BaseURL       string `toml:"base_url"`
	TimeoutSec    int    `toml:"timeout_seconds"`
	MaxRetries    int    `toml:"max_retries"`
}

type DefaultsConfig struct {
	MaxRounds         int    `toml:"max_rounds"`
	AutoSave          bool   `toml:"auto_save"`
	Language          string `toml:"language"`
	DefaultPreset     string `toml:"default_preset"`
	SummarizerEnabled bool   `toml:"summarizer_enabled"`
	SummarizerModel   string `toml:"summarizer_model"`
}

type TUIConfig struct {
	Theme          string `toml:"theme"`
	ShowTokenCount bool   `toml:"show_token_count"`
	ShowCost       bool   `toml:"show_cost"`
	StreamSpeed    string `toml:"stream_speed"`
	Editor         string `toml:"editor"`
}

type ContextConfig struct {
	MaxWindowRatio  float64 `toml:"max_window_ratio"`
	SummaryModel    string  `toml:"summary_model"`
	SummaryMaxToken int     `toml:"summary_max_tokens"`
}

type BudgetConfig struct {
	SessionLimitUSD float64 `toml:"session_limit_usd"`
	WarnThreshold   float64 `toml:"warn_threshold"`
}

type PathsConfig struct {
	SessionsDir string `toml:"sessions_dir"`
	AgentsDir   string `toml:"agents_dir"`
}

type Loaded struct {
	Config     Config
	ConfigPath string
	Paths      RuntimePaths
}

func Load(configPath string) (Loaded, error) {
	if configPath == "" {
		configPath = os.Getenv("WOOLF_CONFIG")
	}
	if configPath == "" {
		configPath = DefaultConfigPath()
	}

	cfg := Default()
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return Loaded{}, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return Loaded{}, err
	}

	applyEnv(&cfg)
	paths := ResolveRuntimePaths(cfg)
	paths.ConfigPath = configPath
	return Loaded{Config: cfg, ConfigPath: configPath, Paths: paths}, nil
}

func Save(path string, cfg Config) error {
	if path == "" {
		path = DefaultConfigPath()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func EnsureRuntimeDirs(paths RuntimePaths) error {
	for _, dir := range []string{paths.DataDir, paths.SessionsDir, paths.AgentsDir, paths.CacheDir, paths.LogsDir} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	return nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("OPENROUTER_API_KEY"); v != "" {
		cfg.API.OpenRouterKey = v
	}
	if v := os.Getenv("WOOLF_SESSIONS_DIR"); v != "" {
		cfg.Paths.SessionsDir = v
	}
}

func MaskSecret(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "****"
	}
	return value[:8] + "****"
}
