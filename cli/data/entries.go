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

	var entries []*ExportableEntry
	err := data.ReadAllEntries(c.Configuration.DataRoot, func(entry *data.Entry) error {
		entries = append(entries, convert(entry))
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

func convert(entry *data.Entry) *ExportableEntry {
	return &ExportableEntry{
		Title: entry.Title,
		Path:  entry.Path,
		Tags:  entry.Tags.ToSlice(),
	}
}

type ExportableEntry struct {
	Title string
	Path  string
	Tags  []string
}
