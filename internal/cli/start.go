package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newStartCommand(app *App) *cobra.Command {
	var draft string
	var preset string
	var agents string
	var rounds int
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Create a new session",
		RunE: func(cmd *cobra.Command, args []string) error {
			if preset == "" {
				preset = app.loaded.Config.Defaults.DefaultPreset
			}
			if rounds == 0 {
				rounds = app.loaded.Config.Defaults.MaxRounds
			}
			title := "untitled"
			if draft != "" {
				title = strings.TrimSuffix(filepath.Base(draft), filepath.Ext(draft))
			}
			sess, path, err := app.store.Create(title, draft)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "created session: %s\n", sess.SessionID)
			fmt.Fprintf(cmd.OutOrStdout(), "path: %s\n", path)
			fmt.Fprintf(cmd.OutOrStdout(), "preset: %s\n", preset)
			fmt.Fprintf(cmd.OutOrStdout(), "rounds: %d\n", rounds)
			if agents != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "agents: %s\n", agents)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&draft, "draft", "", "draft file")
	cmd.Flags().StringVar(&preset, "preset", "", "agent preset")
	cmd.Flags().StringVar(&agents, "agents", "", "comma-separated agent names")
	cmd.Flags().IntVar(&rounds, "rounds", 0, "discussion rounds")
	return cmd
}
