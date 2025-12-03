package mainview

import (
	"code-snippets/ui/components/entrylist"

	tea "github.com/charmbracelet/bubbletea"
)

type GeneralMode struct{}

func (mode GeneralMode) onKeyPressed(model Model, message tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch message.String() {
	case "q":
		return model, tea.Quit

	case "t":
		model.mode = TagInputMode{}
		return model, nil

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
