package taginput

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	completedTags     []string
	inProgress        string
	completedTagStyle lipgloss.Style
	inProgressStyle   lipgloss.Style
}

func New() Model {
	return Model{
		completedTags:     nil,
		inProgress:        "",
		completedTagStyle: lipgloss.NewStyle().Background(lipgloss.Color("#AAFFAA")),
		inProgressStyle:   lipgloss.NewStyle(),
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
			model.inProgress += message.String()
			return model, nil

		case "backspace":
			if len(model.inProgress) > 0 {
				model.removeLastCharacterFromInProgress()
			} else {
				model.dropLastCompletedTag()
			}
			return model, nil

		case "ctrl+w":
			if len(model.inProgress) > 0 {
				model.clearInProgress()
			} else {
				model.clearCompletedTags()
			}

		case " ":
			model.completedTags = append(model.completedTags, model.inProgress)
			model.inProgress = ""
		}
	}

	return model, nil
}

func (model Model) View() string {
	var parts []string

	for _, completedTag := range model.completedTags {
		styledTag := model.completedTagStyle.Render(completedTag)
		parts = append(parts, styledTag, " ")
	}

	styledInProgress := model.inProgressStyle.Render(model.inProgress)
	parts = append(parts, styledInProgress)

	return lipgloss.JoinHorizontal(0, parts...)
}

func (model *Model) removeLastCharacterFromInProgress() {
	if len(model.inProgress) > 0 {
		model.inProgress = model.inProgress[:len(model.inProgress)-1]
	}
}

func (model *Model) clearInProgress() {
	model.inProgress = ""
}

func (model *Model) dropLastCompletedTag() {
	if len(model.completedTags) > 0 {
		model.completedTags = model.completedTags[:len(model.completedTags)-1]
	}
}

func (model *Model) clearCompletedTags() {
	model.completedTags = nil
}
