package detector

import (
	"encoding/json"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

func readJSONFile[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func readTOMLFile(path string) (tree map[string]any, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err = toml.Decode(string(data), &tree); err != nil {
		return nil, err
	}

	return tree, nil
}

func readYAMLFile[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var t T
	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	return &t, nil
}
