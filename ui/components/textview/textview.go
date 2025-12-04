package textview

import (
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size   util.Size
	source string
}

func New() Model {
	model := Model{
		source: "",
	}

	return model
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResize(message)

	case MsgSetSource:
		return model.onSetSource(message)
	}

	return model, nil
}

func (model Model) View() string {
	style := lipgloss.NewStyle().MaxWidth(model.size.Width).MaxHeight(model.size.Height)
	return style.Render(model.source)
}

func (model Model) onSetSource(message MsgSetSource) (tea.Model, tea.Cmd) {
	model.source = message.Source
	return model, nil
}

func (model Model) onResize(message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}
	return model, nil
}
