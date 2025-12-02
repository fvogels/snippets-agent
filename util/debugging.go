package util

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

func DebugShowMessage(message tea.Msg) {
	_, file, _, _ := runtime.Caller(1)

	slog.Debug("message received", slog.String("file", file), slog.String("message", DebugMessageToString(message)))
}

func DebugMessageToString(message tea.Msg) string {
	switch message := message.(type) {
	case tea.KeyMsg:
		return fmt.Sprintf("KeyMsg[%s]", message.String())

	case tea.WindowSizeMsg:
		return fmt.Sprintf("WindowSizeMsg[%d x %d]", message.Width, message.Height)

	default:
		return reflect.TypeOf(message).String()
	}
}
