package data

import (
	"code-snippets/cli/common"
	"code-snippets/data"
	"fmt"

	"github.com/spf13/cobra"
)

type listDataFilesCommand struct {
	common.Command
}

func NewListDataFilesCommand() *cobra.Command {
	var command *listDataFilesCommand

	command = &listDataFilesCommand{
		Command: common.Command{
			CobraCommand: cobra.Command{
				Use:   "ls",
				Short: "List all data files",
				Long:  `List all data files`,
				Args:  cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return command.execute()
				},
			},
		},
	}

	return command.AsCobraCommand()
}

func (c *listDataFilesCommand) execute() error {
	c.LoadConfiguration()

	data.FindFiles(c.Configuration.DataRoot, func(path string) error {
		fmt.Println(path)
		return nil
	})

	return nil
}
