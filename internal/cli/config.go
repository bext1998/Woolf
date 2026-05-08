package cli

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"

	woolfconfig "woolf/internal/config"
)

func newConfigCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage config",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "show",
			Short: "Show current config",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg := app.loaded.Config
				cfg.API.OpenRouterKey = woolfconfig.MaskSecret(cfg.API.OpenRouterKey)
				data, err := toml.Marshal(cfg)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "path: %s\n%s", app.loaded.ConfigPath, data)
				return nil
			},
		},
		&cobra.Command{
			Use:   "reset",
			Short: "Reset config to defaults",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := woolfconfig.Save(app.loaded.ConfigPath, woolfconfig.Default()); err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "config reset: %s\n", app.loaded.ConfigPath)
				return nil
			},
		},
		&cobra.Command{
			Use:   "edit",
			Short: "Print config path",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Fprintln(cmd.OutOrStdout(), app.loaded.ConfigPath)
			},
		},
	)
	return cmd
}
