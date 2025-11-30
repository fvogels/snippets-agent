package data

import (
	"code-snippets/markdown"
	"code-snippets/util"
	"fmt"
	"io/fs"
	"os"
	pathlib "path"
	"path/filepath"
	"reflect"
)

func FindFiles(rootDirectory string, callback func(path string) error) error {
	walker := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Necessary on Windows
		path = filepath.ToSlash(path)

		if info.IsDir() {
			base := pathlib.Base(path)
			if base == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		if err := callback(path); err != nil {
			return err
		}

		return nil
	}

	return filepath.Walk(rootDirectory, walker)
}

type Entry struct {
	Identifier int              `json:"identifier"`
	Path       string           `json:"path"`
	Title      string           `json:"title"`
	Tags       util.Set[string] `json:"tags"`
}

func (entry *Entry) LoadSource() (string, error) {
	data, err := os.ReadFile(entry.Path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ReadEntry(path string, identifier int) (*Entry, error) {
	markdownFile, err := markdown.ParseFile(path)
	if err != nil {
		return nil, err
	}

	titleObject, ok := markdownFile.Metadata["title"]
	if !ok {
		return nil, fmt.Errorf("title missing in %s", path)
	}
	title, ok := titleObject.(string)
	if !ok {
		return nil, fmt.Errorf("title %v is not a string in %s", title, path)
	}

	tagsObject, ok := markdownFile.Metadata["tags"]
	if !ok {
		return nil, fmt.Errorf("tags missing in %s", path)
	}

	var tags util.Set[string]
	if tagsList, ok := tagsObject.([]string); ok {
		tags = util.NewSetFromSlice(tagsList)
	} else if tagsList, ok := tagsObject.([]any); ok {
		tags = util.NewSet[string]()
		for _, tagObject := range tagsList {
			tag, ok := tagObject.(string)
			if !ok {
				return nil, fmt.Errorf("tags %v have invalid type %s in %s", tagsObject, path, reflect.TypeOf(tagsObject))
			}
			tags.Add(tag)
		}
	} else if tag, ok := tagsObject.(string); ok {
		tags = util.NewSetFromSlice([]string{tag})
	} else {
		return nil, fmt.Errorf("tags %v have invalid type %s in %s", tagsObject, path, reflect.TypeOf(tagsObject))
	}

	entry := Entry{
		Identifier: identifier,
		Path:       path,
		Title:      title,
		Tags:       tags,
	}

	return &entry, nil
}

func ReadAllEntries(rootDirectory string, callback func(*Entry) error) error {
	counter := 0

	return FindFiles(rootDirectory, func(path string) error {
		entry, err := ReadEntry(path, counter)
		counter++

		if err != nil {
			return err
		}

		if err := callback(entry); err != nil {
			return err
		}

		return nil
	})
}
