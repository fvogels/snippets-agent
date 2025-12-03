package taginput

import (
	"code-snippets/debug"

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
		inProgressStyle:   lipgloss.NewStyle().Background(lipgloss.Color("#AAAAFF")),
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case MsgAddCharacter:
		return model.onAddCharacter(message)

	case MsgClearSingle:
		return model.onClearSingle()

	case MsgClearAll:
		return model.onClearAll()

	case MsgAddTag:
		return model.onAddTag()
	}

	return model, nil
}

func (model Model) onAddCharacter(message MsgAddCharacter) (tea.Model, tea.Cmd) {
	model.inProgress += message.Character
	return model, model.signalInputChanged()
}

func (model Model) onClearSingle() (tea.Model, tea.Cmd) {
	if len(model.inProgress) > 0 {
		command := model.removeLastCharacterFromInProgress()
		return model, command
	} else {
		command := model.dropLastCompletedTag()
		return model, command
	}
}

func (model Model) onClearAll() (tea.Model, tea.Cmd) {
	if len(model.inProgress) > 0 {
		command := model.clearInProgress()
		return model, command
	} else {
		command := model.clearCompletedTags()
		return model, command
	}
}

func (model Model) onAddTag() (tea.Model, tea.Cmd) {
	model.completedTags = append(model.completedTags, model.inProgress)
	model.inProgress = ""
	return model, tea.Batch(
		model.signalSelectedTagsChanged(),
		model.signalInputChanged(),
	)
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

func (model *Model) removeLastCharacterFromInProgress() tea.Cmd {
	if len(model.inProgress) > 0 {
		model.inProgress = model.inProgress[:len(model.inProgress)-1]
	}

	return model.signalInputChanged()
}

func (model *Model) clearInProgress() tea.Cmd {
	model.inProgress = ""
	return model.signalInputChanged()
}

func (model *Model) dropLastCompletedTag() tea.Cmd {
	if len(model.completedTags) > 0 {
		model.completedTags = model.completedTags[:len(model.completedTags)-1]
		return model.signalSelectedTagsChanged()
	}

	return nil
}

func (model *Model) clearCompletedTags() tea.Cmd {
	if len(model.completedTags) > 0 {
		model.completedTags = nil
		return model.signalSelectedTagsChanged()
	}

	return nil
}

func signal(message tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return message
	}
}

func (model *Model) signalSelectedTagsChanged() tea.Cmd {
	return signal(MsgSelectedTagsChanged{
		SelectedTags: model.completedTags,
	})
}

func (model *Model) signalInputChanged() tea.Cmd {
	return signal(MsgInputChanged{
		Input: model.inProgress,
	})
}

func (model *Model) GetTags() []string {
	return model.completedTags
}

func (model *Model) GetPartiallyInputtedTag() string {
	return model.inProgress
}
