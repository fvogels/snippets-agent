package common

import (
	"code-snippets/configuration"
	"fmt"

	"github.com/spf13/cobra"
)

type Command struct {
	CobraCommand  cobra.Command
	Configuration *configuration.Configuration
}

func (command *Command) PrintErrorf(formatString string, args ...any) {
	fmt.Fprintf(command.CobraCommand.ErrOrStderr(), formatString, args...)
}

func (command *Command) Printf(formatString string, args ...any) {
	fmt.Fprintf(command.CobraCommand.OutOrStdout(), formatString, args...)
}

func (command *Command) AsCobraCommand() *cobra.Command {
	return &command.CobraCommand
}

func (command *Command) LoadConfiguration() error {
	configurationPath, err := configuration.GetPath()
	if err != nil {
		return fmt.Errorf("failed to determine configuration file path: %w", err)
	}

	configuration, err := configuration.Load(configurationPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration file: %w", err)
	}

	command.Configuration = configuration
	return nil
}
