package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"woolf/internal/agents"
)

func newAgentsCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Manage agent roles and presets",
	}
	preset := &cobra.Command{
		Use:   "preset",
		Short: "Manage agent presets",
	}
	preset.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List presets",
			RunE: func(cmd *cobra.Command, args []string) error {
				registry, err := agents.NewRegistry(app.loaded.Paths.AgentsDir)
				if err != nil {
					return err
				}
				for _, item := range registry.ListPresets() {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", item.Name, item.DisplayName, strings.Join(item.Roles, ","))
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "show <name>",
			Short: "Show preset",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				registry, err := agents.NewRegistry(app.loaded.Paths.AgentsDir)
				if err != nil {
					return err
				}
				preset, ok := registry.Preset(args[0])
				if !ok {
					return fmt.Errorf("CFG-003: preset %s not found", args[0])
				}
				fmt.Fprintf(cmd.OutOrStdout(), "name: %s\n", preset.Name)
				fmt.Fprintf(cmd.OutOrStdout(), "display_name: %s\n", preset.DisplayName)
				fmt.Fprintf(cmd.OutOrStdout(), "roles: %s\n", strings.Join(preset.Roles, ", "))
				return nil
			},
		},
	)
	cmd.AddCommand(preset)
	return cmd
}
