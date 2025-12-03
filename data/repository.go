package data

import (
	"code-snippets/util"
	"sort"
)

type Repository interface {
	ListTags() []string
	EnumerateTags(callback func(string) error) error
	EnumerateEntries(selectedTags util.Set[string], callback func(*Entry) error) error
}

type MemoryRepository struct {
	tagTable map[string][]*Entry
	entries  []*Entry
}

func (repository *MemoryRepository) ListTags() []string {
	var tags []string

	for tag := range repository.tagTable {
		tags = append(tags, tag)
	}

	return tags
}

func (repository *MemoryRepository) EnumerateTags(callback func(string) error) error {
	for tag := range repository.tagTable {
		if err := callback(tag); err != nil {
			return err
		}
	}

	return nil
}

func (repository *MemoryRepository) EnumerateEntries(selectedTags util.Set[string], callback func(*Entry) error) error {
	for _, entry := range repository.entries {
		if selectedTags.IsSubsetOf(entry.Tags) {
			if err := callback(entry); err != nil {
				return err
			}
		}
	}

	return nil
}

func LoadRepository(rootDirectory string) (Repository, error) {
	tagTable := make(map[string][]*Entry)
	var entries []*Entry

	err := ReadAllEntries(rootDirectory, func(entry *Entry) error {
		entries = append(entries, entry)

		for _, tag := range entry.Tags.ToSlice() {
			entriesWithTag, ok := tagTable[tag]
			if !ok {
				entriesWithTag = nil
			}

			updatedEntriesWithTag := append(entriesWithTag, entry)
			tagTable[tag] = updatedEntriesWithTag
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Title < entries[j].Title
	})

	repository := MemoryRepository{
		tagTable: tagTable,
		entries:  entries,
	}

	return &repository, nil
}
