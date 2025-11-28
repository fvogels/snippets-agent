package markdown

import (
	"code-snippets/cli/common"
	"code-snippets/markdown"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type extractMarkdownMetadataCommand struct {
	common.Command
}

func NewExtractMarkdownMetadataCommand() *cobra.Command {
	var command *extractMarkdownMetadataCommand

	command = &extractMarkdownMetadataCommand{
		Command: common.Command{
			CobraCommand: cobra.Command{
				Use:   "meta",
				Short: "Markdown meta extraction",
				Long:  `Prints out the metadata of a markdown file.`,
				Args:  cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					return command.execute(args[0])
				},
			},
		},
	}

	return command.AsCobraCommand()
}

func (c *extractMarkdownMetadataCommand) execute(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, metadata := markdown.Parse(source)

	for key, value := range metadata {
		fmt.Printf("%s: %v\n", key, value)
	}

	return nil
}
