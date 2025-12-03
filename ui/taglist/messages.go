package taglist

import "code-snippets/ui/components/stringlist"

type MsgSetTags struct {
	Tags []string
}

type MsgSetFilter = stringlist.MsgSetFilter
