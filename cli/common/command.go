package common

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Command struct {
	CobraCommand cobra.Command
}

func (command *Command) PrintErrorf(formatString string, args ...any) {
	fmt.Fprintf(command.CobraCommand.ErrOrStderr(), formatString, args...)
}

func (command *Command) Printf(formatString string, args ...any) {
	fmt.Fprintf(command.CobraCommand.OutOrStdout(), formatString, args...)
}

func (command *Command) AsCobraCommand() *cobra.Command {
	return &command.CobraCommand
}
