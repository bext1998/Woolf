package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newModelsCommand(app *App) *cobra.Command {
	var pricing bool
	cmd := &cobra.Command{
		Use:   "models",
		Short: "Show model cache path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "models cache: %s\n", filepath.Join(app.loaded.Paths.CacheDir, "models.json"))
			if pricing {
				fmt.Fprintln(cmd.OutOrStdout(), "pricing: OpenRouter pricing cache will be wired later")
			}
		},
	}
	cmd.Flags().BoolVar(&pricing, "pricing", false, "show pricing")
	return cmd
}
