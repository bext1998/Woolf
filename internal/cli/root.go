package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"woolf/internal/config"
	"woolf/internal/session"
)

type App struct {
	configPath string
	verbose    bool
	noColor    bool
	loaded     config.Loaded
	store      session.Store
}

func NewRootCommand() *cobra.Command {
	app := &App{}
	cmd := &cobra.Command{
		Use:   "woolf",
		Short: "Woolf AI writing salon CLI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.CommandPath() == "woolf" {
				return nil
			}
			return app.load()
		},
	}
	cmd.PersistentFlags().StringVar(&app.configPath, "config", "", "config file path")
	cmd.PersistentFlags().BoolVar(&app.verbose, "verbose", false, "verbose output")
	cmd.PersistentFlags().BoolVar(&app.noColor, "no-color", false, "disable color output")

	cmd.AddCommand(
		newInitCommand(app),
		newStartCommand(app),
		newResumeCommand(app),
		newListCommand(app),
		newShowCommand(app),
		newExportCommand(app),
		newConfigCommand(app),
		newModelsCommand(app),
	)
	return cmd
}

func (app *App) load() error {
	loaded, err := config.Load(app.configPath)
	if err != nil {
		return err
	}
	if err := config.EnsureRuntimeDirs(loaded.Paths); err != nil {
		return err
	}
	app.loaded = loaded
	app.store = session.NewStore(loaded.Paths.SessionsDir)
	return nil
}

func (app *App) printPaths(cmd *cobra.Command) {
	fmt.Fprintf(cmd.OutOrStdout(), "config: %s\n", app.loaded.ConfigPath)
	fmt.Fprintf(cmd.OutOrStdout(), "sessions: %s\n", app.loaded.Paths.SessionsDir)
	fmt.Fprintf(cmd.OutOrStdout(), "agents: %s\n", app.loaded.Paths.AgentsDir)
	fmt.Fprintf(cmd.OutOrStdout(), "cache: %s\n", app.loaded.Paths.CacheDir)
}
