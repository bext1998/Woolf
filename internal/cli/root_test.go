package cli

import "testing"

func TestNewRootCommandRegistersExpectedSubcommands(t *testing.T) {
	cmd := NewRootCommand()

	if cmd.Use != "woolf" {
		t.Fatalf("root command use = %q, want %q", cmd.Use, "woolf")
	}

	expected := []string{
		"init",
		"start",
		"resume",
		"list",
		"show",
		"export",
		"config",
		"models",
	}
	for _, name := range expected {
		if found, _, err := cmd.Find([]string{name}); err != nil || found == nil || found.Name() != name {
			t.Fatalf("expected subcommand %q to be registered", name)
		}
	}
}
