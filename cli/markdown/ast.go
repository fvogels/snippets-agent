package markdown

import (
	"code-snippets/cli/common"
	"code-snippets/markdown"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark/ast"
	yaml "gopkg.in/yaml.v2"
)

type extractMarkdownAstCommand struct {
	common.Command
}

func NewExtractAstCommand() *cobra.Command {
	var command *extractMarkdownAstCommand

	command = &extractMarkdownAstCommand{
		Command: common.Command{
			CobraCommand: cobra.Command{
				Use:   "ast",
				Short: "Markdown AST dump",
				Long:  `Prints out the AST of a markdown file.`,
				Args:  cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					return command.execute(args[0])
				},
			},
		},
	}

	return command.AsCobraCommand()
}

func (c *extractMarkdownAstCommand) execute(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	document, _ := markdown.Parse(source)
	tree := convert(document, source)

	buffer, err := yaml.Marshal(tree)
	if err != nil {
		return err
	}

	fmt.Println(string(buffer))

	return nil
}

type AstNode struct {
	TypeName string
	Children []any `yaml:"children,omitempty"`
}

type DocumentNode struct {
	AstNode
}

type TextNode struct {
	AstNode `yaml:",inline"`
	Content string
}

func convert(node ast.Node, markdown []byte) any {
	children := convertChildren(node, markdown)
	typeName := reflect.TypeOf(node).String()

	switch node := node.(type) {
	case *ast.Document:
		return &DocumentNode{
			AstNode: AstNode{
				TypeName: typeName,
				Children: children,
			},
		}

	case *ast.Text:
		return &struct {
			AstNode `yaml:",inline"`
			Content string
		}{
			AstNode: AstNode{
				TypeName: typeName,
				Children: children,
			},
			Content: string(node.Segment.Value(markdown)),
		}

	case *ast.Heading:
		return &struct {
			AstNode `yaml:",inline"`
			Level   int
		}{
			AstNode: AstNode{
				TypeName: typeName,
				Children: children,
			},
			Level: node.Level,
		}

	case *ast.FencedCodeBlock:
		return &struct {
			AstNode  `yaml:",inline"`
			Info     any
			Language string
			Lines    string
		}{
			AstNode: AstNode{
				TypeName: typeName,
				Children: children,
			},
			Info:     convert(node.Info, markdown),
			Language: string(node.Language(markdown)),
			Lines:    string(node.Lines().Value(markdown)),
		}

	default:
		return &AstNode{
			TypeName: typeName,
			Children: children,
		}
	}
}

func convertChildren(node ast.Node, markdown []byte) []any {
	if node.HasChildren() {
		result := []any{}

		child := node.FirstChild()

		for child != nil {
			result = append(result, convert(child, markdown))
			child = child.NextSibling()
		}

		return result
	} else {
		return nil
	}
}
