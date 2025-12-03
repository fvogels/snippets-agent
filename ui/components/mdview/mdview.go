package mdview

import (
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	source           []byte
	renderedMarkdown string
	size             util.Size
}

func New() Model {
	model := Model{
		source: nil,
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

	case msgRenderingDone:
		model.renderedMarkdown = message.renderedMarkdown
		return model, nil
	}

	return model, nil
}

func (model Model) View() string {
	borderStyle := lipgloss.NewStyle().Width(model.size.Width).Height(model.size.Height)
	innerStyle := lipgloss.NewStyle().MaxWidth(model.size.Width - 2).MaxHeight(model.size.Height - 2)
	return borderStyle.Render(innerStyle.Render(model.renderedMarkdown))
}

func (model Model) onSetSource(message MsgSetSource) (tea.Model, tea.Cmd) {
	width := model.size.Width

	command := func() tea.Msg {
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(width-2),
		)
		if err != nil {
			panic("failed to create markdown renderer")
		}
		renderedMarkdown, err := renderer.Render(message.Source)
		if err != nil {
			panic("failed to render markdown file")
		}
		return msgRenderingDone{
			renderedMarkdown: renderedMarkdown,
		}
	}

	return model, command
}

func (model Model) onResize(message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}
	return model, nil
}
