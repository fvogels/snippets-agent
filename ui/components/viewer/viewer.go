package viewer

import (
	"code-snippets/ui/components/mdview"
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
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
	return model.viewer.View()
}

func (model Model) onSetDocument(message MsgSetDocument) (tea.Model, tea.Cmd) {
	model.document = message.Document

	switch document := model.document.(type) {
	case Markdown:
		viewer := mdview.New()
		model.viewer = viewer
		command1 := viewer.Init()
		updatedViewer, command2 := viewer.Update(mdview.MsgSetSource{Source: document.Source})
		model.viewer = updatedViewer

		return model, tea.Sequence(command1, command2)

	default:
		panic("unsupported document type")
	}
}

func (model Model) onResize(message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	updatedViewer, command := model.viewer.Update(message)
	model.viewer = updatedViewer

	return model, command
}
