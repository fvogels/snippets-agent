package markdown

import (
	"code-snippets/cli/common"
	"code-snippets/markdown"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type extractMarkdownFencedCodeBlocksCommand struct {
	common.Command
}

func NewExtractMarkdownFencedCodeBlocksCommand() *cobra.Command {
	var command *extractMarkdownFencedCodeBlocksCommand

	command = &extractMarkdownFencedCodeBlocksCommand{
		Command: common.Command{
			CobraCommand: cobra.Command{
				Use:   "blocks",
				Short: "Code blocks extraction",
				Long:  `Prints out the code blocks of a markdown file.`,
				Args:  cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					return command.execute(args[0])
				},
			},
		},
	}

	return command.AsCobraCommand()
}

func (c *extractMarkdownFencedCodeBlocksCommand) execute(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ast, _ := markdown.Parse(source)
	blocks := markdown.ExtractCodeBlocks(source, ast)

	for index, block := range blocks {
		fmt.Printf("Block %d: language %s\n", index, string(block.Language))
		fmt.Println(string(block.Content))
		fmt.Println()
	}

	return nil
}
