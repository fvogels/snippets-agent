package horizontal

import (
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size     util.Size
	children []child
}

type child struct {
	widthFunction func(size util.Size) int
	model         tea.Model
}

func New() Model {
	return Model{
		children: nil,
		size:     util.Size{Width: 0, Height: 0},
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.size.Width = message.Width
		model.size.Height = message.Height

		commands := []tea.Cmd{}

		for index, child := range model.children {
			updatedChild, command := child.model.Update(tea.WindowSizeMsg{
				Width:  child.widthFunction(model.size),
				Height: message.Height,
			})
			model.children[index].model = updatedChild
			commands = append(commands, command)
		}

		return model, tea.Batch(commands...)

	default:
		commands := []tea.Cmd{}

		for index, child := range model.children {
			updatedChild, command := child.model.Update(message)
			model.children[index].model = updatedChild
			commands = append(commands, command)
		}

		return model, tea.Batch(commands...)
	}
}

func (model Model) View() string {
	parts := []string{}

	for _, child := range model.children {
		style := lipgloss.NewStyle().Width(child.widthFunction(model.size)).Height(model.size.Height)
		part := style.Render(child.model.View())
		parts = append(parts, part)
	}

	style := lipgloss.NewStyle().Width(model.size.Width).Height(model.size.Height)
	return style.Render(lipgloss.JoinHorizontal(0, parts...))
}

func (model *Model) Add(widthFunction func(size util.Size) int, childModel tea.Model) {
	child := child{
		widthFunction: widthFunction,
		model:         childModel,
	}

	model.children = append(model.children, child)
}
