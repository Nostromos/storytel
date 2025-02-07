package parser

import (
	"encoding/json"
	"os"
)

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Storybook map[string]Arc

func LoadStory(filePath string) (Storybook, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var story Storybook
	err = json.NewDecoder(file).Decode(&story)
	if err != nil {
		panic(err)
	}

	return story, nil
}
