package taginput

// MsgAddCharacter is an incoming message
type MsgAddCharacter struct {
	Character string
}

// MsgAddTag is an incoming message
type MsgAddTag struct{}

// MsgClearSingle is an incoming message
type MsgClearSingle struct{}

// MsgClearAll is an incoming message
type MsgClearAll struct{}

// MsgSelectedTagsChanged is an outgoing message
type MsgSelectedTagsChanged struct {
	SelectedTags []string
}

// MsgInputChanged is an outgoing message
type MsgInputChanged struct {
	Input string
}
