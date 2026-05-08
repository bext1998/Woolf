package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"woolf/internal/session"
)

func newListCommand(app *App) *cobra.Command {
	var limit int
	var since string
	var status string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := app.store.List(session.ListOptions{
				Limit:  limit,
				Status: session.Status(status),
			})
			if err != nil {
				return err
			}
			if len(items) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no sessions found. sessions: %s\n", app.store.Dir())
				return nil
			}
			for i, item := range items {
				fmt.Fprintf(cmd.OutOrStdout(), "[%d] %s %s (%s, %d rounds)\n", i+1, item.SessionID, item.Title, item.Status, item.Rounds)
			}
			if since != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "since filter is reserved: %s\n", since)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum rows")
	cmd.Flags().StringVar(&since, "since", "", "start date")
	cmd.Flags().StringVar(&status, "status", "", "session status")
	return cmd
}
