package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"woolf/internal/session"
)

func newForkCommand(app *App) *cobra.Command {
	var draft string
	var title string
	cmd := &cobra.Command{
		Use:   "fork <session-id>",
		Short: "Fork a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" && draft != "" {
				title = strings.TrimSuffix(filepath.Base(draft), filepath.Ext(draft))
			}
			var source *session.Source
			if draft != "" {
				src, err := sourceFromDraft(draft)
				if err != nil {
					return err
				}
				source = &src
			}
			opts := session.ForkOptions{Title: title, Source: source}
			sess, path, err := app.store.Fork(args[0], opts)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "forked session: %s\npath: %s\n", sess.SessionID, path)
			return nil
		},
	}
	cmd.Flags().StringVar(&draft, "draft", "", "replacement draft file")
	cmd.Flags().StringVar(&title, "title", "", "fork title")
	return cmd
}
