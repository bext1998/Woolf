package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
				return fmt.Errorf("only md export is wired in the Phase 1 skeleton")
			}
			sess, _, err := app.store.Load(args[0])
			if err != nil {
				return err
			}
			content := fmt.Sprintf("# %s\n\n- session_id: %s\n- status: %s\n- rounds_completed: %d\n", sess.Title, sess.SessionID, sess.Status, sess.Totals.RoundsCompleted)
			if output == "" {
				fmt.Fprint(cmd.OutOrStdout(), content)
				return nil
			}
			return os.WriteFile(output, []byte(content), 0o600)
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "export format: md|pdf")
	cmd.Flags().StringVar(&output, "output", "", "output path")
	return cmd
}
