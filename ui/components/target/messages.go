package target

import tea "github.com/charmbracelet/bubbletea"

type MsgTargetted struct {
	Target  Identifier
	Message tea.Msg
}
