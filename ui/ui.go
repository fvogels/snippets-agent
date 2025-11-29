package ui

import (
	"code-snippets/configuration"
	"code-snippets/data"
	"code-snippets/ui/entrylist"
	"code-snippets/ui/taginput"
	"code-snippets/ui/taglist"
	"code-snippets/util"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	repository   data.Repository
	screenWidth  int
	screenHeight int
	tagList      taglist.Model
	entryList    entrylist.Model
	tagInput     taginput.Model
}

func New(repository data.Repository) tea.Model {
	model := Model{
		repository:   repository,
		screenWidth:  0,
		screenHeight: 0,
		tagList:      taglist.New(),
		entryList:    entrylist.New(),
		tagInput:     taginput.New(),
	}

	// entries := []*data.Entry{}
	// model.repository.EnumerateEntries(nil, func(entry *data.Entry) error {
	// 	entries = append(entires, entry)
	// 	return nil
	// })
	// model.entryList.SetEntries(entries)
	// model.tagList.SetTags(model.repository.ListTags())
	model.refreshLists()

	return model
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) View() string {
	return lipgloss.JoinVertical(
		0,
		lipgloss.JoinHorizontal(
			0,
			lipgloss.NewStyle().Width(20).Height(model.screenHeight-1).Render(model.tagList.View()),
			lipgloss.NewStyle().Height(model.screenHeight-1).Render(model.entryList.View()),
		),
		model.tagInput.View(),
	)
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "ctrl+c":
			return model, tea.Quit

		default:
			updatedTagInput, command := model.tagInput.Update(message)
			model.tagInput = updatedTagInput
			return model, command
		}

	case tea.WindowSizeMsg:
		model.screenWidth = message.Width
		model.screenHeight = message.Height
		model.tagList.SetWidth(20)
		model.tagList.SetMaximumHeight(model.screenHeight - 1)
		model.entryList.SetWidth(model.screenWidth - 20)
		model.entryList.SetMaximumHeight(model.screenHeight - 1)
		return model, nil

	case taginput.SelectedTagsChangedMessage:
		model.refreshLists()
		return model, nil
	}

	return model, nil
}

func (model *Model) refreshLists() {
	selectedTags := util.NewSetFromSlice(model.tagInput.GetTags())
	entries := []*data.Entry{}

	model.repository.EnumerateEntries(selectedTags, func(entry *data.Entry) error {
		entries = append(entries, entry)
		return nil
	})

	model.entryList.SetEntries(entries)

	remainingTags := util.NewSet[string]()
	model.repository.EnumerateEntries(selectedTags, func(entry *data.Entry) error {
		remainingTags.Union(entry.Tags)
		return nil
	})

	sortedRemainingTags := remainingTags.ToSlice()
	slices.Sort(sortedRemainingTags)

	model.tagList.SetTags(sortedRemainingTags)
}

func Start(configuration *configuration.Configuration) error {
	repository, err := data.LoadRepository(configuration.DataRoot)
	if err != nil {
		return err
	}
	model := New(repository)

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}
