package mdview

import (
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
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
	util.DebugShowMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.size = util.Size{
			Width:  message.Width,
			Height: message.Height,
		}
		return model, nil

	case MsgSetSource:
		width := model.size.Width

		command := func() tea.Msg {
			renderer, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(width),
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

	case msgRenderingDone:
		model.renderedMarkdown = message.renderedMarkdown
		return model, nil
	}

	return model, nil
}

func (model Model) View() string {
	return model.renderedMarkdown
}
