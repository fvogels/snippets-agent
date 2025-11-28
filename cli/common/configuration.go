package common

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	DataDirectory string
}

func GetConfigurationFilePath() (string, error) {
	parentDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filename := ".snippets.toml"
	fullPath := path.Join(parentDirectory, filename)

	return fullPath, nil
}

func LoadConfiguration(path string) (*Configuration, error) {
	var configuration Configuration
	_, err := toml.DecodeFile(path, &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, err
}
