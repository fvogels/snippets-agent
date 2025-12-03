package data

import (
	"code-snippets/markdown"
	"code-snippets/util"
	"os"
	"strings"

	"github.com/yuin/goldmark/ast"
)

type Entry struct {
	Identifier int              `json:"identifier"`
	Path       string           `json:"path"`
	Title      string           `json:"title"`
	Tags       util.Set[string] `json:"tags"`
}

type EntryData struct {
	source     []byte
	ast        ast.Node
	codeBlocks []markdown.CodeBlock
}

// Contents returns the markdown file, excluding the metadata section.
func (data *EntryData) Contents() string {
	contents := string(data.source)
	lines := strings.Lines(contents)

	contentsWithoutMetadata := []string{}
	metadataLinesFound := 0
	lines(func(line string) bool {
		if metadataLinesFound < 2 {
			if strings.TrimSpace(line) == "---" {
				metadataLinesFound++
			}
		} else {
			contentsWithoutMetadata = append(contentsWithoutMetadata, line)
		}
		return true
	})

	return strings.Join(contentsWithoutMetadata, "")
}

func (entry *Entry) LoadData() (EntryData, error) {
	source, err := entry.loadSource()
	if err != nil {
		return EntryData{}, err
	}

	return EntryData{
		source: source,
	}, nil
}

func (entry *Entry) loadSource() ([]byte, error) {
	data, err := os.ReadFile(entry.Path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (entry *EntryData) ParseAST() ast.Node {
	if entry.ast == nil {
		ast, _ := markdown.Parse(entry.source)
		entry.ast = ast
	}

	return entry.ast
}

func (entry *EntryData) GetCodeBlocks() []markdown.CodeBlock {
	if entry.codeBlocks == nil {
		source := entry.source
		ast := entry.ParseAST()
		entry.codeBlocks = markdown.ExtractCodeBlocks(source, ast)
	}

	return entry.codeBlocks
}
