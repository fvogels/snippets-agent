package mainview

import (
	"code-snippets/data"
	"code-snippets/debug"
	"code-snippets/ui/components/entrylist"
	"code-snippets/ui/components/horizontal"
	"code-snippets/ui/components/mdview"
	"code-snippets/ui/components/taginput"
	"code-snippets/ui/components/taglist"
	"code-snippets/ui/components/target"
	"code-snippets/ui/components/vertical"
	"code-snippets/util"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

type Model struct {
	repository           data.Repository
	screenWidth          int
	screenHeight         int
	selectedTags         util.Set[string]
	compatibleTags       []string
	compatibleEntries    []*data.Entry
	partiallyInputtedTag string
	selectedEntry        SelectedEntry
	root                 tea.Model
	tagInputIdentifier   target.Identifier
	tagListIdentifier    target.Identifier
	entryListIdentifier  target.Identifier
	mode                 mode
}

type SelectedEntry struct {
	entry            *data.Entry
	renderedMarkdown string
}

func New(repository data.Repository) tea.Model {
	debug.Milestone()

	tagListWidth := 20
	entryListHeight := 20
	var tagListIdentifier target.Identifier
	var entryListIdentifier target.Identifier
	var tagInputIdentifier target.Identifier

	pane := vertical.New()
	pane.Add(func(size util.Size) int { return entryListHeight }, target.New(entrylist.New(), &entryListIdentifier))
	pane.Add(func(size util.Size) int { return size.Height - entryListHeight }, mdview.New())

	mainView := horizontal.New()
	mainView.Add(func(size util.Size) int { return tagListWidth }, target.New(taglist.New(), &tagListIdentifier))
	mainView.Add(func(size util.Size) int { return size.Width - tagListWidth }, pane)

	root := vertical.New()
	root.Add(func(size util.Size) int { return size.Height - 1 }, mainView)
	root.Add(func(size util.Size) int { return 1 }, target.New(taginput.New(), &tagInputIdentifier))

	model := Model{
		repository:           repository,
		screenWidth:          0,
		screenHeight:         0,
		selectedTags:         util.NewSet[string](),
		compatibleTags:       nil,
		compatibleEntries:    nil,
		partiallyInputtedTag: "",
		root:                 root,
		tagInputIdentifier:   tagInputIdentifier,
		tagListIdentifier:    tagListIdentifier,
		entryListIdentifier:  entryListIdentifier,
		mode:                 GeneralMode{},
		selectedEntry: SelectedEntry{
			entry:            nil,
			renderedMarkdown: "",
		},
	}

	model.recomputeCompatibleTagsAndEntries([]string{})

	return model
}

func (model Model) Init() tea.Cmd {
	debug.Milestone()

	return tea.Batch(
		model.root.Init(),
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
		return model.onKeyPressed(message)

	case tea.WindowSizeMsg:
		return model.onResize(message)

	case taginput.MsgSelectedTagsChanged:
		return model.onSelectedTagsChanged(message.SelectedTags)

	case taginput.MsgInputChanged:
		return model.onPartiallyInputtedTagUpdate(message.Input)

	case MsgMarkdownRendered:
		model.selectedEntry.renderedMarkdown = message.renderedMarkdown
		return model, nil

	case entrylist.MsgEntrySelected:
		model.selectedEntry = SelectedEntry{entry: message.Entry}
		return model, model.rerenderMarkdownInBackground()

	case taginput.MsgReleaseFocus:
		model.mode = GeneralMode{}
		return model, nil

	default:
		updatedRoot, command := model.root.Update(message)
		model.root = updatedRoot
		return model, command
	}
}

func (model Model) onResize(message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	model.screenWidth = message.Width
	model.screenHeight = message.Height

	updatedRoot, rootCommand := model.root.Update(message)
	model.root = updatedRoot
	markdownCommand := model.rerenderMarkdownInBackground()
	return model, tea.Batch(rootCommand, markdownCommand)
}

func (model Model) onKeyPressed(message tea.KeyMsg) (tea.Model, tea.Cmd) {
	return model.mode.onKeyPressed(model, message)
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
		model.signalTagInputCandidates(),
	)
}

func (model *Model) recomputeCompatibleTagsAndEntries(updatedSelectedTags []string) {
	model.selectedTags = util.NewSetFromSlice(updatedSelectedTags)
	entries := []*data.Entry{}

	model.repository.EnumerateEntries(model.selectedTags, func(entry *data.Entry) error {
		entries = append(entries, entry)
		return nil
	})

	model.compatibleEntries = entries

	remainingTags := util.NewSet[string]()
	model.repository.EnumerateEntries(model.selectedTags, func(entry *data.Entry) error {
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
		return target.MsgTargetted{
			Target: model.tagListIdentifier,
			Message: taglist.MsgSetTags{
				Tags: model.compatibleTags,
			},
		}
	}
}

func (model *Model) signalRefreshEntryList() tea.Cmd {
	return func() tea.Msg {
		return target.MsgTargetted{
			Target: model.entryListIdentifier,
			Message: entrylist.MsgSetEntries{
				Entries: model.compatibleEntries,
			},
		}
	}
}

func (model *Model) signalTagInputCandidates() tea.Cmd {
	return func() tea.Msg {
		return target.MsgTargetted{
			Target: model.tagInputIdentifier,
			Message: taginput.MsgSetCandidates{
				Candidates: model.compatibleTags,
			},
		}
	}
}

func (model *Model) signalUpdateTagListFilter() tea.Cmd {
	selectedTags := model.selectedTags.Copy()

	return func() tea.Msg {
		return target.MsgTargetted{
			Target: model.tagListIdentifier,
			Message: taglist.MsgSetFilter{
				Predicate: func(tag string) bool {
					if selectedTags.Contains(tag) {
						return false
					}

					return strings.Contains(tag, model.partiallyInputtedTag)
				},
			},
		}
	}
}

func (model *Model) rerenderMarkdownInBackground() tea.Cmd {
	entry := model.selectedEntry.entry

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
		model.selectedEntry.renderedMarkdown = ""
		return nil
	}
}

func (model *Model) copyCodeblockToClipboard() {
	entry := model.selectedEntry.entry
	codeBlocks, err := entry.GetCodeBlocks()
	if err != nil {
		panic("failed to get code blocks from markdown file")
	}

	if len(codeBlocks) == 0 {
		panic("no code block")
	}

	content := codeBlocks[0].Content
	clipboard.Write(clipboard.FmtText, content)
}

type MsgMarkdownRendered struct {
	renderedMarkdown string
}

type MsgEntryLoaded struct {
	Source string
}
