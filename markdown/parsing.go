package markdown

import (
	"os"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type MarkdownFile struct {
	Path     string
	AST      ast.Node
	Metadata map[string]any
}

func Parse(source []byte) (ast.Node, map[string]any) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	context := parser.NewContext()

	reader := text.NewReader(source)
	ast := markdown.Parser().Parse(reader, parser.WithContext(context))
	metadata := meta.Get(context)

	return ast, metadata
}

func ParseFile(path string) (*MarkdownFile, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ast, metadata := Parse(source)
	markdownFile := MarkdownFile{
		Path:     path,
		AST:      ast,
		Metadata: metadata,
	}

	return &markdownFile, nil
}

type CodeBlock struct {
	Language string
	Content  string
}

func ExtractCodeBlocks(source []byte, node ast.Node) []CodeBlock {
	result := []CodeBlock{}

	findFencedCodeBlocks(node, func(fencedCodeBlock *ast.FencedCodeBlock) {
		language := string(fencedCodeBlock.Language(source))
		content := string(fencedCodeBlock.Lines().Value(source))
		codeblock := CodeBlock{
			Language: language,
			Content:  content,
		}

		result = append(result, codeblock)
	})

	return result
}

func findFencedCodeBlocks(node ast.Node, callback func(*ast.FencedCodeBlock)) {
	switch node := node.(type) {
	case *ast.FencedCodeBlock:
		callback(node)

	default:
		child := node.FirstChild()
		for child != nil {
			findFencedCodeBlocks(child, callback)
			child = child.NextSibling()
		}
	}
}
