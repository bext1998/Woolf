package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type RuntimePaths struct {
	ConfigPath  string
	DataDir     string
	SessionsDir string
	AgentsDir   string
	CacheDir    string
	LogsDir     string
}

func DefaultConfigPath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(firstEnv("APPDATA", homeDir()), "woolf", "config.toml")
	case "darwin":
		return filepath.Join(homeDir(), "Library", "Application Support", "woolf", "config.toml")
	default:
		base := firstEnv("XDG_CONFIG_HOME", filepath.Join(homeDir(), ".config"))
		return filepath.Join(base, "woolf", "config.toml")
	}
}

func DefaultDataDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(firstEnv("LOCALAPPDATA", firstEnv("APPDATA", homeDir())), "woolf")
	case "darwin":
		return filepath.Join(homeDir(), "Library", "Application Support", "woolf")
	default:
		base := firstEnv("XDG_DATA_HOME", filepath.Join(homeDir(), ".local", "share"))
		return filepath.Join(base, "woolf")
	}
}

func ResolveRuntimePaths(cfg Config) RuntimePaths {
	dataDir := DefaultDataDir()
	sessionsDir := cfg.Paths.SessionsDir
	if sessionsDir == "" {
		sessionsDir = filepath.Join(dataDir, "sessions")
	}
	agentsDir := cfg.Paths.AgentsDir
	if agentsDir == "" {
		agentsDir = filepath.Join(dataDir, "agents")
	}
	return RuntimePaths{
		ConfigPath:  DefaultConfigPath(),
		DataDir:     dataDir,
		SessionsDir: sessionsDir,
		AgentsDir:   agentsDir,
		CacheDir:    filepath.Join(dataDir, "cache"),
		LogsDir:     filepath.Join(dataDir, "logs"),
	}
}

func firstEnv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func homeDir() string {
	if dir, err := os.UserHomeDir(); err == nil {
		return dir
	}
	return "."
}
