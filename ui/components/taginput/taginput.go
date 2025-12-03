package taginput

import (
	"code-snippets/debug"
	"code-snippets/ui/bundle"
	"code-snippets/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	completedTags     []string
	inProgress        string
	completedTagStyle lipgloss.Style
	inProgressStyle   lipgloss.Style
	focused           bool
	size              util.Size
}

func New() Model {
	return Model{
		completedTags:     nil,
		inProgress:        "",
		completedTagStyle: lipgloss.NewStyle().Background(lipgloss.Color("#AAFFAA")),
		inProgressStyle:   lipgloss.NewStyle(),
		focused:           false,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.KeyMsg:
		return model.onKeyPressed(message)

	case bundle.MessageBundle:
		return message.UpdateAll(model)

	case tea.WindowSizeMsg:
		return model.onResize(message)

	case tea.FocusMsg:
		return model.onFocus()

	case tea.BlurMsg:
		return model.onBlur()
	}

	return model, nil
}

func (model Model) onKeyPressed(message tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch message.String() {
	case "esc":
		model.focused = false

		commands := []tea.Cmd{model.signalReleaseFocus()}
		if len(model.inProgress) > 0 {
			model.inProgress = ""
			commands = append(commands, model.signalInputChanged())
		}

		return model, tea.Batch(commands...)

	case "backspace":
		return model.onClearSingle()

	case "ctrl+w":
		return model.onClearAll()

	case "enter":
		model.focused = false
		updatedModel, command := model.onAddTag()
		return updatedModel, tea.Batch(command, model.signalReleaseFocus())

	case " ":
		return model.onAddTag()

	default:
		if len(message.String()) == 1 {
			char := message.String()[0]

			if util.IsLowercaseLetter(char) || util.IsDigit(char) || char == '-' {
				return model.onAddCharacter(char)
			}
		}

		return model, nil
	}
}

func (model Model) signalReleaseFocus() tea.Cmd {
	return func() tea.Msg {
		return MsgReleaseFocus{}
	}
}

func (model Model) onResize(message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	model.size.Width = message.Width
	model.size.Height = message.Height

	return model, nil
}

func (model Model) onFocus() (tea.Model, tea.Cmd) {
	model.focused = true
	return model, nil
}

func (model Model) onBlur() (tea.Model, tea.Cmd) {
	model.focused = false
	return model, nil
}

func (model Model) onAddCharacter(char byte) (tea.Model, tea.Cmd) {
	model.inProgress += string(char)
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
	if len(model.inProgress) > 0 {
		model.completedTags = append(model.completedTags, model.inProgress)
		model.inProgress = ""
		return model, bundle.BundleCommands(
			model.signalSelectedTagsChanged(),
			model.signalInputChanged(),
		)
	} else {
		return model, nil
	}
}

func (model Model) View() string {
	renderedCompletedTags := model.renderSelectedTags()
	renderedInProgress := model.renderInProgressTag(model.size.Width - lipgloss.Width(renderedCompletedTags))
	return lipgloss.JoinHorizontal(0, renderedCompletedTags, renderedInProgress)
}

func (model *Model) renderSelectedTags() string {
	var completedParts []string

	for _, completedTag := range model.completedTags {
		styledTag := model.completedTagStyle.Render(completedTag)
		completedParts = append(completedParts, styledTag, " ")
	}

	return lipgloss.JoinHorizontal(0, completedParts...)
}

func (model *Model) renderInProgressTag(width int) string {
	style := model.inProgressStyle.Width(width)
	if model.focused {
		style = style.Background(lipgloss.Color("#AAA"))
	}
	return style.Render(model.inProgress)
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
