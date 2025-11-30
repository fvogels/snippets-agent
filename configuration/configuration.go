package configuration

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	DataRoot string
	KeepLog  bool
}

func GetPath() (string, error) {
	parentDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filename := ".snippets.toml"
	fullPath := path.Join(parentDirectory, filename)

	return fullPath, nil
}

func Load(path string) (*Configuration, error) {
	var configuration Configuration
	_, err := toml.DecodeFile(path, &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, err
}
