package entrylist

import (
	"code-snippets/data"
	"code-snippets/ui/stringlist"
)

type MsgSelectPrevious = stringlist.MsgSelectPrevious
type MsgSelectNext = stringlist.MsgSelectNext

type MsgSetEntries struct {
	Entries []*data.Entry
}
