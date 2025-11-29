package data

import (
	"github.com/spf13/cobra"
)

func NewDataCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "data",
		Short: "Manage data",
		Long:  `Manage data`,
	}

	command.AddCommand(NewListDataFilesCommand())
	command.AddCommand(NewListEntriesCommand())

	return &command
}
