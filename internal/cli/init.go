package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"woolf/internal/config"
)

func newInitCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize Woolf config and data directories",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := app.configPath
			if path == "" {
				path = os.Getenv("WOOLF_CONFIG")
			}
			if path == "" {
				path = config.DefaultConfigPath()
			}
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if err := config.Save(path, config.Default()); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			if err := app.load(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Woolf initialized")
			app.printPaths(cmd)
			return nil
		},
	}
}
