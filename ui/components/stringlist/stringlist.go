package stringlist

import (
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	allowSelection     bool
	items              []string
	filteredItems      []string
	firstVisibleIndex  int
	selectedIndex      int
	width              int
	maximumHeight      int
	filter             func(item string) bool
	emptyListMessage   string
	messageTransformer func(tea.Msg) tea.Msg
}

func New(allowSelection bool) Model {
	model := Model{
		allowSelection:     allowSelection,
		items:              nil,
		firstVisibleIndex:  0,
		selectedIndex:      0,
		width:              0,
		maximumHeight:      0,
		filter:             func(item string) bool { return true },
		messageTransformer: func(m tea.Msg) tea.Msg { return m },
	}

	return model
}

func (model Model) Init() tea.Cmd {
	return model.signalItemSelected()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case MsgSelectPrevious:
		return model.onSelectPrevious()

	case MsgSelectNext:
		return model.onSelectNext()

	case tea.WindowSizeMsg:
		return model.onResize(message)

	case MsgSetFilter:
		return model.onSetFilter(message)

	case MsgSetItems:
		return model.onSetItems(message)
	}

	return model, nil
}

func (model Model) onSetFilter(message MsgSetFilter) (tea.Model, tea.Cmd) {
	model.filter = message.Predicate
	model.refresh()
	return model, model.signalItemSelected()
}

func (model Model) onSetItems(message MsgSetItems) (tea.Model, tea.Cmd) {
	model.items = message.Items
	model.refresh()
	return model, model.signalItemSelected()
}

func (model *Model) signalItemSelected() tea.Cmd {
	if len(model.items) > 0 {
		selectedIndex := model.selectedIndex
		selectedItem := model.items[selectedIndex]

		return func() tea.Msg {
			msg := MsgItemSelected{
				Index: selectedIndex,
				Item:  selectedItem,
			}

			return model.messageTransformer(msg)
		}
	} else {
		return nil
	}
}

func (model Model) onSelectPrevious() (tea.Model, tea.Cmd) {
	if model.allowSelection {
		if model.selectedIndex > 0 {
			model.selectedIndex--
		}
		model.ensureSelectedIsVisible()

		return model, model.signalItemSelected()
	} else {
		return model, nil
	}
}

func (model Model) onSelectNext() (tea.Model, tea.Cmd) {
	if model.allowSelection {
		if model.selectedIndex != -1 && model.selectedIndex+1 < len(model.items) {
			model.selectedIndex++
		}
		model.ensureSelectedIsVisible()

		return model, model.signalItemSelected()
	} else {
		return model, nil
	}
}

func (model Model) onResize(message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	model.width = message.Width
	model.maximumHeight = message.Height
	return model, nil
}

func (model Model) View() string {
	itemsToBeShown := model.filteredItems

	if len(itemsToBeShown) == 0 {
		return model.emptyListMessage
	}

	itemStyle := lipgloss.NewStyle().Width(model.width)
	selectedItemStyle := itemStyle.Background(lipgloss.Color("#AAAAAA"))
	rowIndex := 0
	var lines []string

	for rowIndex < model.maximumHeight && model.firstVisibleIndex+rowIndex < len(itemsToBeShown) {
		itemIndex := model.firstVisibleIndex + rowIndex
		item := itemsToBeShown[itemIndex]

		var line string
		if model.allowSelection && itemIndex == model.selectedIndex {
			line = selectedItemStyle.Render(item)
		} else {
			line = itemStyle.Render(item)
		}
		lines = append(lines, line)

		rowIndex++
	}

	return lipgloss.JoinVertical(0, lines...)
}

func (model *Model) refresh() {
	model.filteredItems = util.Filter(model.items, model.filter)

	model.firstVisibleIndex = 0
	if model.selectedIndex >= 0 {
		model.selectedIndex = 0
	}
}

func (model *Model) ensureSelectedIsVisible() {
	if model.allowSelection {
		if model.firstVisibleIndex > model.selectedIndex {
			model.firstVisibleIndex = model.selectedIndex
		} else if model.firstVisibleIndex+model.maximumHeight < model.selectedIndex {
			model.firstVisibleIndex = model.selectedIndex - model.maximumHeight + 1
			if model.firstVisibleIndex < 0 {
				model.firstVisibleIndex = 0
			}
		}
	}
}

func (model *Model) SetEmptyListMessage(message string) {
	model.emptyListMessage = message
}

func (model *Model) GetSelectedItem() string {
	return model.filteredItems[model.GetSelectedIndex()]
}

func (model *Model) GetSelectedIndex() int {
	if model.selectedIndex == -1 {
		panic("selecting is not enabled")
	}

	return model.selectedIndex
}

func (model *Model) SetMessageTransformer(transformer func(tea.Msg) tea.Msg) {
	model.messageTransformer = transformer
}
