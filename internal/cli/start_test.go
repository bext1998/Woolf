package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"woolf/internal/openrouter"
)

type fakeStartClient struct{}

func (fakeStartClient) StreamChat(ctx context.Context, req openrouter.ChatRequest) (<-chan openrouter.StreamEvent, error) {
	ch := make(chan openrouter.StreamEvent, 1)
	ch <- openrouter.StreamEvent{Content: req.Model, Usage: &openrouter.Usage{PromptTokens: 1, CompletionTokens: 1, TotalTokens: 2}}
	close(ch)
	return ch, nil
}

func TestStartCommandRunsPipelineWithFakeClient(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	sessionsDir := filepath.Join(dir, "sessions")
	draftPath := filepath.Join(dir, "draft.md")
	if err := os.WriteFile(draftPath, []byte("# Draft\n\nHello"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := NewRootCommandWithClient(fakeStartClient{})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{
		"--config", configPath,
		"start",
		"--draft", draftPath,
		"--preset", "editorial",
		"--rounds", "1",
	})
	t.Setenv("WOOLF_SESSIONS_DIR", sessionsDir)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v\n%s", err, out.String())
	}
	matches, err := filepath.Glob(filepath.Join(sessionsDir, "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("session files = %d, want 1; output:\n%s", len(matches), out.String())
	}
}
