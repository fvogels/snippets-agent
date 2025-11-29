package data

import (
	"code-snippets/markdown"
	"fmt"
	"io/fs"
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
	Path  string
	Title string
	Tags  []string
}

func ReadEntry(path string) (*Entry, error) {
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

	var tags []string

	if tagsList, ok := tagsObject.([]string); ok {
		tags = tagsList
	} else if tagsList, ok := tagsObject.([]any); ok {
		for _, tagObject := range tagsList {
			tag, ok := tagObject.(string)
			if !ok {
				return nil, fmt.Errorf("tags %v have invalid type %s in %s", tagsObject, path, reflect.TypeOf(tagsObject))
			}
			tags = append(tags, tag)
		}
	} else if tag, ok := tagsObject.(string); ok {
		tags = []string{tag}
	} else {
		return nil, fmt.Errorf("tags %v have invalid type %s in %s", tagsObject, path, reflect.TypeOf(tagsObject))
	}

	entry := Entry{
		Path:  path,
		Title: title,
		Tags:  tags,
	}

	return &entry, nil
}

func ReadAllEntries(rootDirectory string, callback func(*Entry) error) error {
	return FindFiles(rootDirectory, func(path string) error {
		entry, err := ReadEntry(path)
		if err != nil {
			return err
		}

		if err := callback(entry); err != nil {
			return err
		}

		return nil
	})
}
