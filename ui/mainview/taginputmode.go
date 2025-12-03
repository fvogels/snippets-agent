package mainview

import (
	"code-snippets/ui/components/taginput"
	"code-snippets/ui/components/target"

	tea "github.com/charmbracelet/bubbletea"
)

type TagInputMode struct{}

func (mode TagInputMode) onKeyPressed(model Model, message tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch message.String() {
	case "esc":
		model.mode = GeneralMode{}
		return model, func() tea.Msg {
			return taginput.MsgSetFocus{Focused: false}
		}

	default:
		// Ensure all key commands only reach the tag input
		updatedRoot, command := model.root.Update(target.MsgTargetted{
			Target:  model.tagInputIdentifier,
			Message: message,
		})

		model.root = updatedRoot
		return model, command
	}
}
