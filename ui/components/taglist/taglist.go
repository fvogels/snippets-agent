package taglist

import (
	"code-snippets/debug"
	"code-snippets/ui/components/stringlist"
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	stringList tea.Model
	size       util.Size
}

func New() Model {
	stringList := stringlist.New(false)
	model := Model{
		stringList: stringList,
	}

	emptyListMessageStyle := lipgloss.NewStyle().Italic(true)
	emptyListMessage := emptyListMessageStyle.Render("no tags found")
	stringList.SetEmptyListMessage(emptyListMessage)

	return model
}

func (model Model) Init() tea.Cmd {
	return model.stringList.Init()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResize(message)

	case MsgSetTags:
		return model.onSetTags(message)

	case MsgSetFilter:
		return model.onSetFilter(message)

	default:
		updatedStringList, command := model.stringList.Update(message)
		model.stringList = updatedStringList
		return model, command
	}
}

func (model Model) onSetFilter(message MsgSetFilter) (tea.Model, tea.Cmd) {
	updatedStringList, command := model.stringList.Update(message)
	model.stringList = updatedStringList
	return model, command
}

func (model Model) onSetTags(message MsgSetTags) (tea.Model, tea.Cmd) {
	updatedStringList, command := model.stringList.Update(stringlist.MsgSetItems{
		Items: message.Tags,
	})
	model.stringList = updatedStringList
	return model, command
}

func (model Model) onResize(message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	model.size = util.Size{Width: message.Width, Height: message.Height}
	updatedStringList, command := model.stringList.Update(tea.WindowSizeMsg{
		Width:  message.Width - 2,
		Height: message.Height - 2,
	})
	model.stringList = updatedStringList
	return model, command
}

func (model Model) View() string {
	// Note that the border is drawn outside the given width and height, so we need to decrease them by 2 to compensate
	style := lipgloss.NewStyle().Width(model.size.Width-2).Height(model.size.Height-2).Border(lipgloss.DoubleBorder(), true).Padding(0)
	return style.Render(model.stringList.View())
}
