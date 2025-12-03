package entrylist

import (
	"code-snippets/data"
	"code-snippets/ui/components/stringlist"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgSelectPrevious = stringlist.MsgSelectPrevious
type MsgSelectNext = stringlist.MsgSelectNext

type MsgSetEntries struct {
	Entries []*data.Entry
}

type msgStringListMessageWrapper struct {
	message tea.Msg
}

type MsgEntrySelected struct {
	Index int
	Entry *data.Entry
}
