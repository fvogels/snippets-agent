package mainview

import (
	"code-snippets/ui/components/taginput"
	"code-snippets/util"

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

	case "backspace":
		return model, func() tea.Msg {
			return taginput.MsgClearSingle{}
		}

	case "ctrl+w":
		return model, func() tea.Msg {
			return taginput.MsgClearAll{}
		}

	case " ":
		return model, func() tea.Msg {
			return taginput.MsgAddTag{}
		}

	default:
		if len(message.String()) == 1 {
			char := message.String()[0]

			if util.IsLowercaseLetter(char) || util.IsDigit(char) || char == '-' {
				return model, func() tea.Msg {
					return taginput.MsgAddCharacter{Character: message.String()}
				}
			}
		}

		return model, nil
	}
}
