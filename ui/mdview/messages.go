package mdview

type MsgSetSource struct {
	Source string
}

type msgRenderingDone struct {
	renderedMarkdown string
}
