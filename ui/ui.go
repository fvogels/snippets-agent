package ui

import (
	"code-snippets/configuration"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) View() string {
	return "tik"
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "q":
			return model, tea.Quit
		}
	}

	return model, nil
}

func Start(configuration *configuration.Configuration) error {
	model := Model{}

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}
