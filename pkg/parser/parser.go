package parser

import (
	"encoding/json"
	"os"

	"github.com/Nostromos/cyoa/pkg/types"
)

func LoadStory(filePath string) (types.Storybook, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var storybook types.Storybook
	err = json.NewDecoder(file).Decode(&storybook)
	if err != nil {
		panic(err)
	}

	return storybook, nil
}
