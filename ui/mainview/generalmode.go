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
				Target:  model.targets.tagInput,
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

	case "1":
		model.copyCodeblockToClipboard(0)
		return model, nil

	case "2":
		model.copyCodeblockToClipboard(1)
		return model, nil

	case "3":
		model.copyCodeblockToClipboard(2)
		return model, nil

	case "4":
		model.copyCodeblockToClipboard(3)
		return model, nil

	case "5":
		model.copyCodeblockToClipboard(4)
		return model, nil

	case "6":
		model.copyCodeblockToClipboard(5)
		return model, nil

	case "7":
		model.copyCodeblockToClipboard(6)
		return model, nil

	case "8":
		model.copyCodeblockToClipboard(7)
		return model, nil

	case "9":
		model.copyCodeblockToClipboard(8)
		return model, nil

	case "0":
		model.copyCodeblockToClipboard(9)
		return model, nil

	default:
		return model, nil
	}
}
