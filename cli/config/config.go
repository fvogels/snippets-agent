package config

import (
	"github.com/spf13/cobra"
)

func NewConfigurationCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  `Manage configuration`,
	}

	command.AddCommand(NewShowConfigurationPathCommand())

	return &command
}
