package ui

import (
	"code-snippets/configuration"
	"code-snippets/data"
	"code-snippets/debug"
	"code-snippets/ui/entrylist"
	"code-snippets/ui/horizontal"
	"code-snippets/ui/mdview"
	"code-snippets/ui/taginput"
	"code-snippets/ui/taglist"
	"code-snippets/ui/vertical"
	"code-snippets/util"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

type Model struct {
	repository           data.Repository
	screenWidth          int
	screenHeight         int
	compatibleTags       []string
	compatibleEntries    []*data.Entry
	renderedMarkdown     string
	partiallyInputtedTag string
	selectedEntry        *data.Entry
	root                 tea.Model
}

func New(repository data.Repository) tea.Model {
	debug.Milestone()

	tagListWidth := 20
	entryListHeight := 20

	pane := vertical.New()
	pane.Add(func(size util.Size) int { return entryListHeight }, entrylist.New())
	pane.Add(func(size util.Size) int { return size.Height - entryListHeight }, mdview.New())

	mainView := horizontal.New()
	mainView.Add(func(size util.Size) int { return tagListWidth }, taglist.New())
	mainView.Add(func(size util.Size) int { return size.Width - tagListWidth }, pane)

	root := vertical.New()
	root.Add(func(size util.Size) int { return size.Height - 1 }, mainView)
	root.Add(func(size util.Size) int { return 1 }, taginput.New())

	model := Model{
		repository:   repository,
		screenWidth:  0,
		screenHeight: 0,
		root:         root,
	}

	model.recomputeCompatibleTagsAndEntries([]string{})

	debug.Milestone()

	return model
}

func (model Model) Init() tea.Cmd {
	debug.Milestone()

	return tea.Batch(
		model.signalRefreshTagList(),
		model.signalRefreshEntryList(),
	)
}

func (model Model) View() string {
	return model.root.View()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "esc":
			return model, tea.Quit

		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
			return model, func() tea.Msg {
				return taginput.MsgAddCharacter{Character: message.String()}
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

		case "down":
			return model, func() tea.Msg {
				return entrylist.MsgSelectNext{}
			}

		case "up":
			return model, func() tea.Msg {
				return entrylist.MsgSelectPrevious{}
			}

		case "enter":
			command := model.rerenderMarkdownInBackground()
			return model, command

		case "ctrl+c":
			model.copyCodeblockToClipboard()
			return model, nil

		default:
			updatedRoot, command := model.root.Update(message)
			model.root = updatedRoot
			return model, command
		}

	case tea.WindowSizeMsg:
		model.screenWidth = message.Width
		model.screenHeight = message.Height

		updatedRoot, rootCommand := model.root.Update(message)
		model.root = updatedRoot
		markdownCommand := model.rerenderMarkdownInBackground()
		return model, tea.Batch(rootCommand, markdownCommand)

	case taginput.MsgSelectedTagsChanged:
		return model.onSelectedTagsChanged(message.SelectedTags)

	case taginput.MsgInputChanged:
		return model.onPartiallyInputtedTagUpdate(message.Input)

	case MsgMarkdownRendered:
		model.renderedMarkdown = message.renderedMarkdown
		return model, nil

	case entrylist.MsgEntrySelected:
		model.selectedEntry = message.Entry
		return model, nil

	default:
		updatedRoot, command := model.root.Update(message)
		model.root = updatedRoot
		return model, command
	}
}

func (model Model) onSelectedTagsChanged(updatedSelectedTags []string) (tea.Model, tea.Cmd) {
	model.recomputeCompatibleTagsAndEntries(updatedSelectedTags)

	return model, tea.Batch(
		model.signalRefreshTagList(),
		model.signalRefreshEntryList(),
	)
}

func (model Model) onPartiallyInputtedTagUpdate(partiallyInputtedTag string) (tea.Model, tea.Cmd) {
	model.partiallyInputtedTag = partiallyInputtedTag

	return model, tea.Batch(
		model.signalUpdateTagListFilter(),
		model.signalRefreshTagList(),
	)
}

func (model *Model) recomputeCompatibleTagsAndEntries(updatedSelectedTags []string) {
	selectedTags := util.NewSetFromSlice(updatedSelectedTags)
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

func (model *Model) signalRefreshTagList() tea.Cmd {
	// TODO model.compatibleTags should be copied into a separate array first
	return func() tea.Msg {
		slog.Debug("Sending MsgSetTags")

		return taglist.MsgSetTags{
			Tags: model.compatibleTags,
		}
	}
}

func (model *Model) signalRefreshEntryList() tea.Cmd {
	return func() tea.Msg {
		slog.Debug("Sending MsgSetEntries")

		return entrylist.MsgSetEntries{
			Entries: model.compatibleEntries,
		}
	}
}

func (model *Model) signalUpdateTagListFilter() tea.Cmd {
	return func() tea.Msg {
		return taglist.MsgSetFilter{
			Predicate: func(tag string) bool {
				return strings.Contains(tag, model.partiallyInputtedTag)
			},
		}
	}
}

func (model *Model) rerenderMarkdownInBackground() tea.Cmd {
	entry := model.selectedEntry

	if entry != nil {
		return func() tea.Msg {
			source, err := entry.GetSource()
			if err != nil {
				panic("failed to load markdown file")
			}

			return mdview.MsgSetSource{
				Source: source,
			}
		}
	} else {
		model.renderedMarkdown = ""
		return nil
	}
}

func (model *Model) copyCodeblockToClipboard() {
	// entry := model.entryList.GetSelectedEntry()
	// codeBlocks, err := entry.GetCodeBlocks()
	// if err != nil {
	// 	panic("failed to get code blocks from markdown file")
	// }

	// if len(codeBlocks) == 0 {
	// 	panic("no code block")
	// }

	// content := codeBlocks[0].Content
	// clipboard.Write(clipboard.FmtText, content)
}

type MsgMarkdownRendered struct {
	renderedMarkdown string
}

type MsgEntryLoaded struct {
	Source string
}

func Start(configuration *configuration.Configuration) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	if configuration.KeepLog {
		logFile, err := os.Create("ui.log")
		if err != nil {
			fmt.Println("Failed to create log")
		}
		defer logFile.Close()

		logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}

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
