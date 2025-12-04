package mainview

import (
	"code-snippets/data"
	"code-snippets/debug"
	"code-snippets/ui/components/entrylist"
	"code-snippets/ui/components/horizontal"
	"code-snippets/ui/components/taginput"
	"code-snippets/ui/components/taglist"
	"code-snippets/ui/components/target"
	"code-snippets/ui/components/vertical"
	"code-snippets/ui/components/viewer"
	"code-snippets/util"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

type Model struct {
	repository              data.Repository
	screenSize              util.Size
	selectedTags            util.Set[string]
	compatibleTags          []string
	entriesWithSelectedTags []*data.Entry
	partiallyInputtedTag    string
	selectedEntry           *SelectedEntry
	root                    tea.Model
	targets                 Targets
	mode                    mode
}

type SelectedEntry struct {
	entry                  data.Entry
	data                   *data.EntryData
	selectedCodeblockIndex *int
}

type Targets struct {
	tagInput  target.Identifier
	tagList   target.Identifier
	entryList target.Identifier
	viewer    target.Identifier
}

func New(repository data.Repository) tea.Model {
	debug.Milestone()

	tagListWidth := 20
	entryListHeight := 20
	targets := Targets{}

	pane := vertical.New()
	pane.Add(func(size util.Size) int { return entryListHeight }, target.New(entrylist.New(), &targets.entryList))
	pane.Add(func(size util.Size) int { return size.Height - entryListHeight }, target.New(viewer.New(), &targets.viewer))

	mainView := horizontal.New()
	mainView.Add(func(size util.Size) int { return tagListWidth }, target.New(taglist.New(), &targets.tagList))
	mainView.Add(func(size util.Size) int { return size.Width - tagListWidth }, pane)

	root := vertical.New()
	root.Add(func(size util.Size) int { return size.Height - 1 }, mainView)
	root.Add(func(size util.Size) int { return 1 }, target.New(taginput.New(), &targets.tagInput))

	model := Model{
		repository:              repository,
		screenSize:              util.Size{Width: 0, Height: 0},
		selectedTags:            util.NewSet[string](),
		compatibleTags:          nil,
		entriesWithSelectedTags: nil,
		partiallyInputtedTag:    "",
		root:                    root,
		targets:                 targets,
		mode:                    GeneralMode{},
		selectedEntry:           nil,
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

	case msgEntryLoaded:
		model.selectedEntry.data = &message.Data
		return model, model.signalUpdateViewer()

	case entrylist.MsgEntrySelected:
		selectedEntry := message.Entry
		model.selectedEntry = &SelectedEntry{
			entry: *selectedEntry,
		}
		return model, model.signalLoadSelectedEntry(selectedEntry)

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
	model.screenSize = util.Size{Width: message.Width, Height: message.Height}

	updatedRoot, command := model.root.Update(message)
	model.root = updatedRoot
	return model, command
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

	model.entriesWithSelectedTags = entries

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
			Target: model.targets.tagList,
			Message: taglist.MsgSetTags{
				Tags: model.compatibleTags,
			},
		}
	}
}

func (model *Model) signalRefreshEntryList() tea.Cmd {
	return func() tea.Msg {
		return target.MsgTargetted{
			Target: model.targets.entryList,
			Message: entrylist.MsgSetEntries{
				Entries: model.entriesWithSelectedTags,
			},
		}
	}
}

func (model *Model) signalTagInputCandidates() tea.Cmd {
	return func() tea.Msg {
		return target.MsgTargetted{
			Target: model.targets.tagInput,
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
			Target: model.targets.tagList,
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

func (model *Model) signalLoadSelectedEntry(entry *data.Entry) tea.Cmd {
	if entry != nil {
		return func() tea.Msg {
			data, err := entry.LoadData()
			if err != nil {
				panic("failed to load entry data")
			}

			return msgEntryLoaded{
				Data: data,
			}
		}
	} else {
		model.selectedEntry = nil
		return nil
	}
}

func (model *Model) signalUpdateViewer() tea.Cmd {
	var source string
	if model.selectedEntry != nil && model.selectedEntry.data != nil {
		if model.selectedEntry.selectedCodeblockIndex == nil {
			// No particular code block selected, so show entire markdown contents
			source = model.selectedEntry.data.Contents()
		} else {
			// Specific code block selected
			codeBlock := model.selectedEntry.data.GetCodeBlock(*model.selectedEntry.selectedCodeblockIndex)
			source = string(codeBlock.Content)
		}
	} else {
		source = ""
	}

	return func() tea.Msg {
		return target.MsgTargetted{
			Target: model.targets.viewer,
			Message: viewer.MsgSetDocument{
				Document: viewer.Markdown{
					Source: source,
				},
			},
		}
	}
}

func (model *Model) copyCodeblockToClipboard() {
	if model.selectedEntry != nil {
		// Determine code block index to use
		var selectedCodeBlockIndex int
		if model.selectedEntry.selectedCodeblockIndex != nil {
			// A code block has been explicitly selected
			selectedCodeBlockIndex = *model.selectedEntry.selectedCodeblockIndex
		} else {
			// No code block has been explicitly selected; default to the first code block
			selectedCodeBlockIndex = 0
		}

		// Ensure that the code block with the given index exists
		if selectedCodeBlockIndex < model.selectedEntry.data.GetCodeBlockCount() {
			codeBlock := model.selectedEntry.data.GetCodeBlock(selectedCodeBlockIndex)
			content := codeBlock.Content
			clipboard.Write(clipboard.FmtText, content)
		}
	}
}

func (model Model) selectCodeblock(index int) (tea.Model, tea.Cmd) {
	if model.selectedEntry != nil && model.selectedEntry.data != nil && index < model.selectedEntry.data.GetCodeBlockCount() {
		model.selectedEntry.selectedCodeblockIndex = &index
		return model, model.signalUpdateViewer()
	} else {
		return model, nil
	}
}

func (model Model) unselectCodeBlock() (tea.Model, tea.Cmd) {
	if model.selectedEntry != nil && model.selectedEntry.selectedCodeblockIndex != nil {
		model.selectedEntry.selectedCodeblockIndex = nil
		return model, model.signalUpdateViewer()
	} else {
		return model, nil
	}
}

func (model Model) selectNextCodeBlock() (tea.Model, tea.Cmd) {
	if model.selectedEntry == nil || model.selectedEntry.data == nil {
		return model, nil
	}

	if model.selectedEntry.selectedCodeblockIndex == nil {
		// No code block was selected, so select first
		index := 0
		model.selectedEntry.selectedCodeblockIndex = &index
	} else {
		index := *model.selectedEntry.selectedCodeblockIndex
		index = (index + 1) % model.selectedEntry.data.GetCodeBlockCount()
		model.selectedEntry.selectedCodeblockIndex = &index
	}

	return model, model.signalUpdateViewer()
}

func (model Model) selectPreviousCodeBlock() (tea.Model, tea.Cmd) {
	if model.selectedEntry == nil || model.selectedEntry.data == nil {
		return model, nil
	}

	if model.selectedEntry.selectedCodeblockIndex == nil {
		// No code block was selected, so select last
		index := model.selectedEntry.data.GetCodeBlockCount() - 1
		model.selectedEntry.selectedCodeblockIndex = &index
	} else {
		index := *model.selectedEntry.selectedCodeblockIndex
		index = (index + model.selectedEntry.data.GetCodeBlockCount() - 1) % model.selectedEntry.data.GetCodeBlockCount()
		model.selectedEntry.selectedCodeblockIndex = &index
	}

	return model, model.signalUpdateViewer()
}
