package target

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Identifier int

var NextId Identifier

type Model struct {
	identifier Identifier
	child      tea.Model
}

func New(child tea.Model, idReceiver *Identifier) Model {
	*idReceiver = NextId
	NextId++

	return Model{
		identifier: *idReceiver,
		child:      child,
	}
}

func (model Model) Init() tea.Cmd {
	return model.child.Init()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case MsgTargetted:
		if message.Target == model.identifier {
			updatedChild, command := model.child.Update(message.Message)
			model.child = updatedChild
			return model, command
		} else {
			updatedChild, command := model.child.Update(message)
			model.child = updatedChild
			return model, command
		}

	default:
		updatedChild, command := model.child.Update(message)
		model.child = updatedChild
		return model, command
	}
}

func (model Model) View() string {
	return model.child.View()
}
