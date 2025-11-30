package ui

import (
	"code-snippets/configuration"
	"code-snippets/data"
	"code-snippets/ui/entrylist"
	"code-snippets/ui/taginput"
	"code-snippets/ui/taglist"
	"code-snippets/util"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	repository        data.Repository
	screenWidth       int
	screenHeight      int
	compatibleTags    []string
	compatibleEntries []*data.Entry

	tagList   taglist.Model
	entryList entrylist.Model
	tagInput  taginput.Model
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

	model.recomputeCompatibleTagsAndEntries()
	model.refreshTagList()
	model.refreshEntryList()

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

		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
			updatedTagInput, command := model.tagInput.Update(taginput.MsgAddCharacter{Character: message.String()})
			model.tagInput = updatedTagInput
			return model, command

		case "backspace":
			updatedTagInput, command := model.tagInput.Update(taginput.MsgClearSingle{})
			model.tagInput = updatedTagInput
			return model, command

		case "ctrl+w":
			updatedTagInput, command := model.tagInput.Update(taginput.MsgClearAll{})
			model.tagInput = updatedTagInput
			return model, command

		case " ":
			updatedTagInput, command := model.tagInput.Update(taginput.MsgAddTag{})
			model.tagInput = updatedTagInput
			return model, command

		case "down":
			updatedEntryList, command := model.entryList.Update(entrylist.MsgSelectNext{})
			model.entryList = updatedEntryList
			return model, command

		case "up":
			updatedEntryList, command := model.entryList.Update(entrylist.MsgSelectPrevious{})
			model.entryList = updatedEntryList
			return model, command

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

	case taginput.MsgSelectedTagsChanged:
		model.recomputeCompatibleTagsAndEntries()
		model.refreshTagList()
		model.refreshEntryList()
		return model, nil

	case taginput.MsgInputChanged:
		model.updateTagListFilter()
		model.refreshTagList()
		return model, nil
	}

	return model, nil
}

func (model *Model) recomputeCompatibleTagsAndEntries() {
	selectedTags := util.NewSetFromSlice(model.tagInput.GetTags())
	entries := []*data.Entry{}

	model.repository.EnumerateEntries(selectedTags, func(entry *data.Entry) error {
		entries = append(entries, entry)
		return nil
	})

	model.compatibleEntries = entries

	remainingTags := util.NewSet[string]()
	model.repository.EnumerateEntries(selectedTags, func(entry *data.Entry) error {
		remainingTags.Union(entry.Tags)
		return nil
	})

	sortedRemainingTags := remainingTags.ToSlice()
	slices.Sort(sortedRemainingTags)

	model.compatibleTags = sortedRemainingTags
}

func (model *Model) refreshTagList() {
	model.tagList.SetTags(model.compatibleTags)
}

func (model *Model) refreshEntryList() {
	model.entryList.SetEntries(model.compatibleEntries)
}

func (model *Model) updateTagListFilter() {
	model.tagList.SetFilter(func(tag string) bool {
		return strings.Contains(tag, model.tagInput.GetPartiallyInputtedTag())
	})
}

func Start(configuration *configuration.Configuration) error {
	logFile, err := os.Create("ui.log")
	if err != nil {
		fmt.Println("Failed to create log")
	}
	defer logFile.Close()

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

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
