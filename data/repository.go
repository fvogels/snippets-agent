package data

type Repository interface {
	ListTags() []string
}

type MemoryRepository struct {
	tagTable map[string][]*Entry
}

func (repository *MemoryRepository) ListTags() []string {
	var tags []string

	for tag := range repository.tagTable {
		tags = append(tags, tag)
	}

	return tags
}

func LoadRepository(rootDirectory string) (Repository, error) {
	tagTable := make(map[string][]*Entry)

	err := ReadAllEntries(rootDirectory, func(entry *Entry) error {
		for _, tag := range entry.Tags {
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

	repository := MemoryRepository{
		tagTable: tagTable,
	}

	return &repository, nil
}
