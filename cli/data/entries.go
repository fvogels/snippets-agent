package data

import (
	"code-snippets/cli/common"
	"code-snippets/data"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type listEntriesCommand struct {
	common.Command
}

func NewListEntriesCommand() *cobra.Command {
	var command *listEntriesCommand

	command = &listEntriesCommand{
		Command: common.Command{
			CobraCommand: cobra.Command{
				Use:   "entries",
				Short: "List all entries",
				Long:  `List all entries`,
				Args:  cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return command.execute()
				},
			},
		},
	}

	return command.AsCobraCommand()
}

func (c *listEntriesCommand) execute() error {
	c.LoadConfiguration()

	var entries []*data.Entry
	err := data.FindFiles(c.Configuration.DataRoot, func(path string) error {
		entry, err := data.ReadEntry(path)
		if err != nil {
			return err
		}

		fmt.Printf("title: %s\n", entry.Title)
		entries = append(entries, entry)
		return nil
	})
	if err != nil {
		return err
	}

	buffer, err := yaml.Marshal(entries)
	if err != nil {
		return err
	}

	fmt.Println(string(buffer))

	return nil
}
