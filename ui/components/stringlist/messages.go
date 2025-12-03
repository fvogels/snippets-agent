package stringlist

type MsgSetItems struct {
	Items []string
}

type MsgSelectPrevious struct{}

type MsgSelectNext struct{}

type MsgSetFilter struct {
	Predicate func(tag string) bool
}

type MsgItemSelected struct {
	Index int
	Item  string
}
