package mainview

import (
	"code-snippets/ui/components/entrylist"
	"code-snippets/ui/components/target"

	tea "github.com/charmbracelet/bubbletea"
)

type GeneralMode struct{}

func (mode GeneralMode) onKeyPressed(model Model, message tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch message.String() {
	case "q":
		return model, tea.Quit

	case " ":
		model.mode = TagInputMode{}
		return model, func() tea.Msg {
			return target.MsgTargetted{
				Target:  model.tagInputIdentifier,
				Message: tea.FocusMsg{},
			}
		}

	case "down":
		return model, func() tea.Msg {
			return entrylist.MsgSelectNext{}
		}

	case "up":
		return model, func() tea.Msg {
			return entrylist.MsgSelectPrevious{}
		}

	case "c":
		model.copyCodeblockToClipboard()
		return model, nil

	default:
		return model, nil
	}
}
