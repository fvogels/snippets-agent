package taginput

// MsgSelectedTagsChanged is an outgoing message
type MsgSelectedTagsChanged struct {
	SelectedTags []string
}

// MsgInputChanged is an outgoing message
type MsgInputChanged struct {
	Input string
}

type MsgReleaseFocus struct{}
