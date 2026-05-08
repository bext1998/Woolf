package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newResumeCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "resume <session-id>",
		Short: "Resume a paused session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sess, path, err := app.store.Resume(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "session: %s\nstatus: %s\npath: %s\n", sess.SessionID, sess.Status, path)
			return nil
		},
	}
}
