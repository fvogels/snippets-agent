package markdown

import (
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func Parse(source []byte) (ast.Node, map[string]any) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	context := parser.NewContext()

	parser := markdown.Parser()
	reader := text.NewReader(source)
	document := parser.Parse(reader)
	metadata := meta.Get(context)

	return document, metadata
}
