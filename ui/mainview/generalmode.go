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

	case "c":
		model.copyCodeblockToClipboard()
		return model, nil

	case "esc", "`":
		return model.unselectCodeBlock()

	case "1":
		return model.selectCodeblock(0)

	case "2":
		return model.selectCodeblock(1)

	case "3":
		return model.selectCodeblock(2)

	case "4":
		return model.selectCodeblock(3)

	case "5":
		return model.selectCodeblock(4)

	case "6":
		return model.selectCodeblock(5)

	case "7":
		return model.selectCodeblock(6)

	case "8":
		return model.selectCodeblock(7)

	case "9":
		return model.selectCodeblock(8)

	case "0":
		return model.selectCodeblock(9)

	default:
		return model, nil
	}
}
