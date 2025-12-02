package taglist

import "code-snippets/ui/stringlist"

type MsgSetTags struct {
	Tags []string
}

type MsgSetFilter = stringlist.MsgSetFilter
