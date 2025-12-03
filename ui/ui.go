package ui

import (
	"code-snippets/configuration"
	"code-snippets/data"
	"code-snippets/ui/mainview"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

func Start(configuration *configuration.Configuration) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	if configuration.KeepLog {
		logFile, err := os.Create("ui.log")
		if err != nil {
			fmt.Println("Failed to create log")
		}
		defer logFile.Close()

		logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}

	repository, err := data.LoadRepository(configuration.DataRoot)
	if err != nil {
		return err
	}
	model := mainview.New(repository)

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}
