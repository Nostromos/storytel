package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/Nostromos/cyoa/pkg/parser"
	"github.com/Nostromos/cyoa/pkg/types"
)

type Mode int

const (
	ModeSelecting Mode = iota
	ModePlaying
)

type StoryFile struct {
	Name  string
	Path  string
	Title string
}

type Model struct {
	Mode        Mode
	StoryFiles  []StoryFile
	CurrentFile string
	Storybook   types.Storybook
	Chapter     types.Chapter
	Title       string
	Story       []string
	Choices     []types.Option
	Cursor      int
	Selected    map[int]struct{}
}

func InitialModel() Model {
	// Get all story files from the stories directory
	storyFiles := []StoryFile{}

	files, err := os.ReadDir("./stories")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			path := filepath.Join("./stories", file.Name())
			// Load the story briefly to get its title
			story, err := parser.LoadStory(path)
			if err != nil {
				continue
			}

			// Try to find the intro chapter to get the title
			var title string
			if intro, ok := story["intro"]; ok {
				title = intro.Title
			} else if intro, ok := story["timetravel-intro"]; ok {
				title = intro.Title
			} else if intro, ok := story["mystery-intro"]; ok {
				title = intro.Title
			} else if intro, ok := story["bakeoff-intro"]; ok {
				title = intro.Title
			} else {
				// Fallback to filename without extension
				title = strings.TrimSuffix(file.Name(), ".json")
			}

			storyFiles = append(storyFiles, StoryFile{
				Name:  file.Name(),
				Path:  path,
				Title: title,
			})
		}
	}

	return Model{
		Mode:       ModeSelecting,
		StoryFiles: storyFiles,
		Cursor:     0,
		Selected:   make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Key Press?
	case tea.KeyMsg:

		// which key?
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			switch m.Mode {
			case ModeSelecting:
				if m.Cursor > 0 {
					m.Cursor--
				}
			case ModePlaying:
				if m.Cursor > 0 {
					m.Cursor--
				}
			}

		case "down", "j":
			switch m.Mode {
			case ModeSelecting:
				if m.Cursor < len(m.StoryFiles)-1 {
					m.Cursor++
				}
			case ModePlaying:
				if m.Cursor < len(m.Choices)-1 {
					m.Cursor++
				}
			}

		case "enter", " ":
			if m.Mode == ModeSelecting {
				// Load the selected story
				selectedFile := m.StoryFiles[m.Cursor]
				story, err := parser.LoadStory(selectedFile.Path)
				if err != nil {
					return m, tea.Quit
				}

				// Find the intro chapter
				var introKey string
				for key := range story {
					if strings.Contains(key, "intro") {
						introKey = key
						break
					}
				}
				if introKey == "" {
					// Fallback to first key
					for key := range story {
						introKey = key
						break
					}
				}

				initialChapter := story[introKey]

				// Switch to playing mode
				m.Mode = ModePlaying
				m.CurrentFile = selectedFile.Name
				m.Storybook = story
				m.Chapter = initialChapter
				m.Title = initialChapter.Title
				m.Story = initialChapter.Story
				m.Choices = initialChapter.Options
				m.Cursor = 0
				m.Selected = make(map[int]struct{})

			} else if m.Mode == ModePlaying {
				selectedOption := m.Choices[m.Cursor]
				nextArc := selectedOption.Arc

				nextChapter, exists := m.Storybook[nextArc]
				if exists {
					m.Chapter = nextChapter
					m.Title = nextChapter.Title
					m.Story = nextChapter.Story
					m.Choices = nextChapter.Options
					m.Cursor = 0                        // reset cursor
					m.Selected = make(map[int]struct{}) // clear selection
				} else {
					return m, tea.Quit
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.Mode == ModeSelecting {
		// Story selection view
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Width(118).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Foreground(lipgloss.Color("#FAFAFA")).
			PaddingTop(2).
			PaddingBottom(2)

		title := "Choose Your Adventure"
		StyledTitle := titleStyle.Render(title)

		// Build story list
		var storyList string = "Select a story:\n\n"
		for i, story := range m.StoryFiles {
			cursor := " "
			if m.Cursor == i {
				cursor = ">"
			}
			storyList += fmt.Sprintf("%s %s\n", cursor, story.Title)
		}

		listStyle := lipgloss.NewStyle().
			Width(120).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#9F56F4")).
			PaddingTop(2).
			PaddingBottom(2).
			PaddingLeft(4).
			PaddingRight(4)

		StyledList := listStyle.Render(storyList)

		help := "\n↑/↓ - navigate | enter - select | q - quit\n"
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D3D3D3")).
			PaddingLeft(2)
		StyledHelp := helpStyle.Render(help)

		return "\n" + StyledTitle + "\n\n" + StyledList + "\n" + StyledHelp

	} else {
		// Story playing view
		rawTitle := m.Title
		rawStory := m.Story
		joinedStory := strings.Join(rawStory, "\n\n")

		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Width(118).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Foreground(lipgloss.Color("#FAFAFA")).
			PaddingTop(2).
			PaddingBottom(2)

		storyStyle := lipgloss.NewStyle().
			Bold(true).
			Width(120).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#9F56F4")).
			PaddingTop(2).
			PaddingBottom(2).
			PaddingLeft(4).
			PaddingRight(4)

		var StyledTitle = titleStyle.Render(rawTitle)
		var StyledStory = storyStyle.Render(joinedStory)

		var opts string
		// the header
		if len(m.Choices) > 0 {
			opts = "Options:\n\n"

			for i, choice := range m.Choices {

				// is cursor on this choice?
				cursor := " "
				if m.Cursor == i {
					cursor = ">" // <-- our cursor
				}

				// render the row
				opts += fmt.Sprintf("%s %s\n", cursor, choice.Text)
			}
		}

		// footer
		help := "\nenter - select option | q - quit\n"
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D3D3D3")).
			PaddingLeft(2)
		StyledHelp := helpStyle.Render(help)

		var optionsStyle = lipgloss.NewStyle().
			Bold(true).
			Width(120).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(2).
			PaddingLeft(4).
			PaddingRight(4)

		var StyledOptions = optionsStyle.Render(opts)

		var text = "\n" + StyledTitle + "\n\n" + StyledStory + "\n\n\n" + StyledOptions + "\n" + StyledHelp
		return text
	}
}
