package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"woolf/internal/exporter"
)

func newExportCommand(app *App) *cobra.Command {
	var format string
	var output string
	cmd := &cobra.Command{
		Use:   "export <session-id>",
		Short: "Export session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if format == "" {
				return fmt.Errorf("export requires --format md|pdf")
			}
			if format != "md" {
				return fmt.Errorf("only md export is supported in this Phase 1 slice")
			}
			sess, _, err := app.store.Load(args[0])
			if err != nil {
				return err
			}
			data, err := exporter.MarkdownExporter{}.Export(sess)
			if err != nil {
				return err
			}
			if output == "" {
				fmt.Fprint(cmd.OutOrStdout(), string(data))
				return nil
			}
			return os.WriteFile(output, data, 0o600)
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "export format: md|pdf")
	cmd.Flags().StringVar(&output, "output", "", "output path")
	return cmd
}
