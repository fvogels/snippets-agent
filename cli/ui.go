package cli

import (
	"code-snippets/cli/common"
	"code-snippets/ui"

	"github.com/spf13/cobra"
)

type startUserInterfaceCommand struct {
	common.Command
}

func NewStartUserInterfaceCommand() *cobra.Command {
	var command *startUserInterfaceCommand

	command = &startUserInterfaceCommand{
		Command: common.Command{
			CobraCommand: cobra.Command{
				Use:   "ui",
				Short: "Start ui",
				Long:  `Start user interface`,
				Args:  cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return command.execute()
				},
			},
		},
	}

	return command.AsCobraCommand()
}

func (c *startUserInterfaceCommand) execute() error {
	if err := c.LoadConfiguration(); err != nil {
		return err
	}

	return ui.Start(c.Configuration)
}
