package stringlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	strings           []string
	firstVisibleIndex int
	width             int
	maximumHeight     int
}

func New() Model {
	return Model{
		strings:           nil,
		firstVisibleIndex: 0,
		width:             0,
		maximumHeight:     0,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (Model, tea.Cmd) {
	return model, nil
}

func (model Model) View() string {
	itemStyle := lipgloss.NewStyle().Width(model.width)
	rowIndex := 0
	var lines []string

	for rowIndex < model.maximumHeight && model.firstVisibleIndex+rowIndex < len(model.strings) {
		itemIndex := model.firstVisibleIndex + rowIndex
		item := model.strings[itemIndex]
		line := itemStyle.Render(item)
		lines = append(lines, line)

		rowIndex++
	}

	return lipgloss.JoinVertical(0, lines...)
}

func (model *Model) SetStrings(strings []string) {
	model.strings = strings
}

func (model *Model) SetWidth(width int) {
	model.width = width
}

func (model *Model) SetMaximumHeight(height int) {
	model.maximumHeight = height
}
