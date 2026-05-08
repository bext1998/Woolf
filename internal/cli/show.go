package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newShowCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "show <session-id>",
		Short: "Show session JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sess, path, err := app.store.Find(args[0])
			if err != nil {
				return err
			}
			data, err := json.MarshalIndent(sess, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "path: %s\n%s\n", path, data)
			return nil
		},
	}
}
