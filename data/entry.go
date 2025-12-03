package data

import (
	"code-snippets/markdown"
	"code-snippets/util"
	"os"
	"strings"
)

type Entry struct {
	Identifier int              `json:"identifier"`
	Path       string           `json:"path"`
	Title      string           `json:"title"`
	Tags       util.Set[string] `json:"tags"`
}

func (entry *Entry) GetSource() (string, error) {
	source, err := entry.loadSource()
	if err != nil {
		return "", err
	}

	contents := string(source)
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

	return strings.Join(contentsWithoutMetadata, ""), nil
}

func (entry *Entry) loadSource() ([]byte, error) {
	data, err := os.ReadFile(entry.Path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (entry *Entry) GetCodeBlocks() ([]markdown.CodeBlock, error) {
	source, err := entry.loadSource()
	if err != nil {
		return nil, err
	}

	ast, _ := markdown.Parse(source)
	codeBlocks := markdown.ExtractCodeBlocks(source, ast)

	return codeBlocks, nil
}
