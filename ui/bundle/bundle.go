package bundle

import tea "github.com/charmbracelet/bubbletea"

type MessageBundle struct {
	Messages []tea.Msg
}

func (bundle MessageBundle) UpdateAll(model tea.Model) (tea.Model, tea.Cmd) {
	commands := []tea.Cmd{}

	for _, message := range bundle.Messages {
		updatedModel, command := model.Update(message)
		model = updatedModel
		commands = append(commands, command)
	}

	return model, tea.Batch(commands...)
}

func BundleCommands(commands ...tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		messages := []tea.Msg{}

		for _, command := range commands {
			message := command()
			messages = append(messages, message)
		}

		return MessageBundle{
			Messages: messages,
		}
	}
}
