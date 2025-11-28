package markdown

import (
	"github.com/spf13/cobra"
)

func NewMarkdownCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "markdown",
		Short: "Markdown functionality",
		Long:  `Markdown functionality`,
	}

	command.AddCommand(NewDumpMarkdownCommand())

	return &command
}
