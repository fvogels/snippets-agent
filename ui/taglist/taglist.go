package taglist

import (
	"code-snippets/ui/stringlist"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	stringList stringlist.Model
}

func New() Model {
	return Model{
		stringList: stringlist.New(false),
	}
}

func (model Model) Init() tea.Cmd {
	return model.stringList.Init()
}

func (model Model) Update(message tea.Msg) (Model, tea.Cmd) {
	updatedStringList, command := model.stringList.Update(message)
	model.stringList = updatedStringList
	return model, command
}

func (model Model) View() string {
	return model.stringList.View()
}

func (model *Model) SetTags(tags []string) {
	model.stringList.SetStrings(tags)
}

func (model *Model) SetWidth(width int) {
	model.stringList.SetWidth(width)
}

func (model *Model) SetMaximumHeight(height int) {
	model.stringList.SetMaximumHeight(height)
}
