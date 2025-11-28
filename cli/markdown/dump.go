package markdown

import (
	"code-snippets/cli/common"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	yaml "gopkg.in/yaml.v2"
)

type dumpMarkdownCommand struct {
	common.Command
}

func NewDumpMarkdownCommand() *cobra.Command {
	var command *dumpMarkdownCommand

	command = &dumpMarkdownCommand{
		Command: common.Command{
			CobraCommand: cobra.Command{
				Use:   "dump",
				Short: "Markdown AST dump",
				Long:  `Prints out the AST of a markdown file.`,
				Args:  cobra.NoArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					return command.execute()
				},
			},
		},
	}

	return command.AsCobraCommand()
}

func (c *dumpMarkdownCommand) execute() error {
	markdown, err := os.ReadFile("f:/repos/code-snippets/data/test/ssh-config.md")
	if err != nil {
		return err
	}

	parser := goldmark.DefaultParser()
	reader := text.NewReader(markdown)
	document := parser.Parse(reader)

	tree := convert(document, markdown)

	// buffer, err := json.MarshalIndent(tree, "", "  ")
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
			AstNode `yaml:",inline"`
			Info    any
			Lines   string
		}{
			AstNode: AstNode{
				TypeName: typeName,
				Children: children,
			},
			Info:  convert(node.Info, markdown),
			Lines: string(node.Lines().Value(markdown)),
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
