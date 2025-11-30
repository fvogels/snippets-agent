package stringlist

import (
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	items             []string
	filteredItems     []string
	firstVisibleIndex int
	selectedIndex     int
	width             int
	maximumHeight     int
	filter            func(item string) bool
}

func New(allowSelection bool) Model {
	model := Model{
		items:             nil,
		firstVisibleIndex: 0,
		selectedIndex:     0,
		width:             0,
		maximumHeight:     0,
		filter:            func(item string) bool { return true },
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
		if model.selectedIndex != -1 && model.selectedIndex+1 < len(model.items) {
			model.selectedIndex++
		}
		model.ensureSelectedIsVisible()
		return model, nil
	}

	return model, nil
}

func (model Model) View() string {
	itemsToBeShown := model.filteredItems

	if len(itemsToBeShown) == 0 {
		return ""
	}

	itemStyle := lipgloss.NewStyle().Width(model.width)
	selectedItemStyle := itemStyle.Background(lipgloss.Color("#AAAAAA"))
	rowIndex := 0
	var lines []string

	for rowIndex < model.maximumHeight && model.firstVisibleIndex+rowIndex < len(itemsToBeShown) {
		itemIndex := model.firstVisibleIndex + rowIndex
		item := itemsToBeShown[itemIndex]

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
	return model.items
}

func (model *Model) SetStrings(strings []string) {
	model.items = strings
	model.refresh()
}

func (model *Model) SetFilter(filter func(item string) bool) {
	model.filter = filter
	model.refresh()
}

func (model *Model) refresh() {
	model.filteredItems = util.Filter(model.items, model.filter)

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
