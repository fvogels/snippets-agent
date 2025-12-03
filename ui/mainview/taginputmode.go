package mainview

import (
	"code-snippets/ui/components/target"

	tea "github.com/charmbracelet/bubbletea"
)

type TagInputMode struct{}

func (mode TagInputMode) onKeyPressed(model Model, message tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Ensure all key commands only reach the tag input
	updatedRoot, command := model.root.Update(target.MsgTargetted{
		Target:  model.tagInputIdentifier,
		Message: message,
	})

	model.root = updatedRoot
	return model, command
}
