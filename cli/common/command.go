package common

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Command struct {
	CobraCommand  cobra.Command
	Configuration *Configuration
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
	configurationPath, err := GetConfigurationFilePath()
	if err != nil {
		return fmt.Errorf("failed to determine configuration file path: %w", err)
	}

	configuration, err := LoadConfiguration(configurationPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration file: %w", err)
	}

	command.Configuration = configuration
	return nil
}
