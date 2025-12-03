package mainview

import (
	tea "github.com/charmbracelet/bubbletea"
)

type mode interface {
	onKeyPressed(model Model, message tea.KeyMsg) (tea.Model, tea.Cmd)
}
