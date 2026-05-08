package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newDeleteCommand(app *App) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "delete <session-id>",
		Short: "Delete a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				fmt.Fprintf(cmd.OutOrStdout(), "delete session %s? type yes to confirm: ", args[0])
				scanner := bufio.NewScanner(cmd.InOrStdin())
				if !scanner.Scan() {
					return scanner.Err()
				}
				if strings.TrimSpace(scanner.Text()) != "yes" {
					fmt.Fprintln(cmd.OutOrStdout(), "delete cancelled")
					return nil
				}
			}
			path, err := app.store.Delete(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted session: %s\n", path)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "delete without confirmation")
	return cmd
}
