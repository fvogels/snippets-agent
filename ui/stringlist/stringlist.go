package stringlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	strings           []string
	firstVisibleIndex int
	selectedIndex     int
	width             int
	maximumHeight     int
}

func New(allowSelection bool) Model {
	model := Model{
		strings:           nil,
		firstVisibleIndex: 0,
		selectedIndex:     0,
		width:             0,
		maximumHeight:     0,
	}

	if !allowSelection {
		model.selectedIndex = -1
	}

	return model
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (Model, tea.Cmd) {
	switch message.(type) {
	case MsgSelectPrevious:
		if model.selectedIndex > 0 {
			model.selectedIndex--
		}
		model.ensureSelectedIsVisible()
		return model, nil

	case MsgSelectNext:
		if model.selectedIndex != -1 && model.selectedIndex+1 < len(model.strings) {
			model.selectedIndex++
		}
		model.ensureSelectedIsVisible()
		return model, nil
	}

	return model, nil
}

func (model Model) View() string {
	if len(model.strings) == 0 {
		return ""
	}

	itemStyle := lipgloss.NewStyle().Width(model.width)
	selectedItemStyle := itemStyle.Background(lipgloss.Color("#AAAAAA"))
	rowIndex := 0
	var lines []string

	for rowIndex < model.maximumHeight && model.firstVisibleIndex+rowIndex < len(model.strings) {
		itemIndex := model.firstVisibleIndex + rowIndex
		item := model.strings[itemIndex]

		var line string
		if itemIndex == model.selectedIndex {
			line = selectedItemStyle.Render(item)
		} else {
			line = itemStyle.Render(item)
		}
		lines = append(lines, line)

		rowIndex++
	}

	return lipgloss.JoinVertical(0, lines...)
}

func (model *Model) GetStrings() []string {
	return model.strings
}

func (model *Model) SetStrings(strings []string) {
	model.strings = strings
	model.firstVisibleIndex = 0

	if model.selectedIndex >= 0 {
		model.selectedIndex = 0
	}
}

func (model *Model) SetWidth(width int) {
	model.width = width
}

func (model *Model) SetMaximumHeight(height int) {
	model.maximumHeight = height
}

func (model *Model) ensureSelectedIsVisible() {
	if model.firstVisibleIndex > model.selectedIndex {
		model.firstVisibleIndex = model.selectedIndex
	} else if model.firstVisibleIndex+model.maximumHeight < model.selectedIndex {
		model.firstVisibleIndex = model.selectedIndex - model.maximumHeight + 1
	}
}
