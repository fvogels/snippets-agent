package debug

import (
	"code-snippets/ui/bundle"
	"code-snippets/util"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func ShowBubbleTeaMessage(message tea.Msg) {
	_, file, line, _ := runtime.Caller(1)

	slog.Debug(
		"message received",
		slog.String("location", fmt.Sprintf("%s:%d", file, line)),
		slog.String("message", BubbleTeaMessageToString(message)),
	)
}

func BubbleTeaMessageToString(message tea.Msg) string {
	switch message := message.(type) {
	case tea.KeyMsg:
		return fmt.Sprintf("KeyMsg[%s]", message.String())

	case tea.WindowSizeMsg:
		return fmt.Sprintf("WindowSizeMsg[%d x %d]", message.Width, message.Height)

	case bundle.MessageBundle:
		return fmt.Sprintf("Bundle[%s]", strings.Join(util.Map(message.Messages, BubbleTeaMessageToString), ", "))

	default:
		return reflect.TypeOf(message).String()
	}
}

func Milestone() {
	_, file, line, _ := runtime.Caller(1)

	slog.Debug(
		"milestone reached",
		slog.String("location", fmt.Sprintf("%s:%d", file, line)),
	)
}
