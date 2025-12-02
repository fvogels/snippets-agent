package entrylist

import (
	"code-snippets/data"
	"code-snippets/ui/stringlist"
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	entries    []*data.Entry
	stringList tea.Model
}

func New() Model {
	return Model{
		entries:    nil,
		stringList: stringlist.New(true),
	}
}

func (model Model) Init() tea.Cmd {
	return model.stringList.Init()
}

func (model Model) Update(message tea.Msg) (Model, tea.Cmd) {
	switch message := message.(type) {
	case MsgSetEntries:
		updatedEntries := message.Entries
		model.entries = message.Entries
		titles := util.Map(updatedEntries, func(entry *data.Entry) string { return entry.Title })

		updatedStringList, command := model.stringList.Update(stringlist.MsgSetItems{
			Items: titles,
		})
		model.stringList = updatedStringList
		return model, command

	case tea.WindowSizeMsg:
		updatedStringList, command := model.stringList.Update(message)
		model.stringList = updatedStringList
		return model, command

	default:
		updatedStringList, command := model.stringList.Update(message)
		model.stringList = updatedStringList
		return model, command
	}
}

func (model Model) View() string {
	return model.stringList.View()
}
