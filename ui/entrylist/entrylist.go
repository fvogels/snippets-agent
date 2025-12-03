package entrylist

import (
	"code-snippets/data"
	"code-snippets/debug"
	"code-snippets/ui/stringlist"
	"code-snippets/util"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	entries    []*data.Entry
	stringList tea.Model
	size       util.Size
}

func New() Model {
	stringList := stringlist.New(true)
	stringList.SetMessageTransformer(func(message tea.Msg) tea.Msg {
		slog.Debug("transforming message")
		return msgStringListMessageWrapper{
			message: message,
		}
	})

	return Model{
		entries:    nil,
		stringList: stringList,
	}
}

func (model Model) Init() tea.Cmd {
	return model.stringList.Init()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case MsgSetEntries:
		return model.onSetEntries(message)

	case tea.WindowSizeMsg:
		model.size = util.Size{Width: message.Width, Height: message.Height}
		updatedStringList, command := model.stringList.Update(tea.WindowSizeMsg{
			Width:  message.Width - 2,
			Height: message.Height - 2,
		})
		model.stringList = updatedStringList
		return model, command

	case msgStringListMessageWrapper:
		switch message := message.message.(type) {
		case stringlist.MsgItemSelected:
			selectedEntry := model.entries[message.Index]
			return model, func() tea.Msg {
				return MsgEntrySelected{
					Index: message.Index,
					Entry: selectedEntry,
				}
			}

		default:
			return model, nil
		}

	default:
		updatedStringList, command := model.stringList.Update(message)
		model.stringList = updatedStringList
		return model, command
	}
}

func (model Model) onSetEntries(message MsgSetEntries) (tea.Model, tea.Cmd) {
	updatedEntries := message.Entries
	model.entries = message.Entries
	titles := util.Map(updatedEntries, func(entry *data.Entry) string { return entry.Title })

	updatedStringList, command := model.stringList.Update(stringlist.MsgSetItems{
		Items: titles,
	})
	model.stringList = updatedStringList
	return model, command
}

func (model Model) View() string {
	// Note that the border is drawn outside the given width and height, so we need to decrease them by 2 to compensate
	style := lipgloss.NewStyle().Width(model.size.Width-2).Height(model.size.Height-2).Border(lipgloss.DoubleBorder(), true).Padding(0)
	return style.Render(model.stringList.View())
}
