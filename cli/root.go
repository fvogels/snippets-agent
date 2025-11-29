package cli

import (
	"code-snippets/cli/config"
	"code-snippets/cli/data"
	"code-snippets/cli/markdown"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var verbose bool

	rootCommand := cobra.Command{
		Use:   "snippets",
		Short: "Snippets CLI",
		Long:  `Snippets Command Line Interface`,
	}

	cobra.OnInitialize(func() {
		if verbose {
			slog.SetLogLoggerLevel(slog.LevelDebug)
			slog.Info("Verbose mode enabled")
		}
	})

	rootCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	rootCommand.AddCommand(config.NewConfigurationCommand())
	rootCommand.AddCommand(markdown.NewMarkdownCommand())
	rootCommand.AddCommand(data.NewDataCommand())
	rootCommand.AddCommand(NewStartUserInterfaceCommand())

	return &rootCommand
}

func Execute() {
	rootCommand := NewRootCommand()

	rootCommand.SilenceUsage = true
	// rootCommand.SilenceErrors = true

	if err := rootCommand.Execute(); err != nil {
		slog.Debug("An error occurred", "error", err.Error())
		os.Exit(1)
	}
}
