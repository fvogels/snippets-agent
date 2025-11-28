package config

import (
	"code-snippets/cli/common"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

type showConfigurationPathCommand struct {
	common.Command
}

func NewShowConfigurationPathCommand() *cobra.Command {
	var command *showConfigurationPathCommand

	command = &showConfigurationPathCommand{
		Command: common.Command{
			CobraCommand: cobra.Command{
				Use:   "path",
				Short: "Show configuration file path",
				Long:  `Prints out the path of the configuration file.`,
				Args:  cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return command.execute()
				},
			},
		},
	}

	return command.AsCobraCommand()
}

func (c *showConfigurationPathCommand) execute() error {
	parentDirectory, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filename := ".snippets.toml"
	fullPath := path.Join(parentDirectory, filename)
	fmt.Println(fullPath)
	return nil
}
