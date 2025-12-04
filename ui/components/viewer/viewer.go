package viewer

import (
	"code-snippets/ui/components/mdview"
	"code-snippets/util"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	document any
	size     util.Size
	viewer   tea.Model
}

func New() Model {
	model := Model{
		document: nil,
		viewer:   mdview.New(),
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

	case MsgSetDocument:
		return model.onSetDocument(message)

	default:
		updatedViewer, command := model.viewer.Update(message)
		model.viewer = updatedViewer
		return model, command
	}
}

func (model Model) View() string {
	slog.Debug("viewer size", "height", model.size.Height)
	style := lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).Width(model.size.Width - 2).Height(model.size.Height - 2)
	return style.Render(model.viewer.View())
}

func (model Model) onSetDocument(message MsgSetDocument) (tea.Model, tea.Cmd) {
	model.document = message.Document

	switch document := model.document.(type) {
	case Markdown:
		viewer := mdview.New()
		model.viewer = viewer
		command1 := viewer.Init()
		viewer2, command2 := viewer.Update(mdview.MsgSetSource{Source: document.Source})
		viewer3, command3 := viewer2.Update(tea.WindowSizeMsg{Width: model.size.Width - 2, Height: model.size.Height - 2})
		model.viewer = viewer3

		return model, tea.Sequence(command1, command2, command3)

	default:
		panic("unsupported document type")
	}
}

func (model Model) onResize(message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	updatedViewer, command := model.viewer.Update(tea.WindowSizeMsg{
		Width:  message.Width - 2,
		Height: message.Height - 2,
	})
	model.viewer = updatedViewer

	return model, command
}
