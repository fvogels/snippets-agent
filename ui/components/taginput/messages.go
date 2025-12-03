package taginput

type MsgSetFocus struct {
	Focused bool
}

// MsgSelectedTagsChanged is an outgoing message
type MsgSelectedTagsChanged struct {
	SelectedTags []string
}

// MsgInputChanged is an outgoing message
type MsgInputChanged struct {
	Input string
}

type MsgRequestBlur struct{}
