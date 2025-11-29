package taglist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tags                 []string
	firstVisibleTagIndex int
	width                int
	maximumHeight        int
}

func New() Model {
	return Model{
		tags:                 nil,
		firstVisibleTagIndex: 0,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (Model, tea.Cmd) {
	return model, nil
}

func (model Model) View() string {
	tagStyle := lipgloss.NewStyle().Width(model.width)
	rowIndex := 0
	var lines []string

	for rowIndex < model.maximumHeight && model.firstVisibleTagIndex+rowIndex < len(model.tags) {
		tagIndex := model.firstVisibleTagIndex + rowIndex
		tag := model.tags[tagIndex]
		line := tagStyle.Render(tag)
		lines = append(lines, line)

		rowIndex++
	}

	return lipgloss.JoinVertical(0, lines...)
}

func (model *Model) SetTags(tags []string) {
	model.tags = tags
}

func (model *Model) SetWidth(width int) {
	model.width = width
}

func (model *Model) SetMaximumHeight(height int) {
	model.maximumHeight = height
}
