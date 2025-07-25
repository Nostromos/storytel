package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/Nostromos/cyoa/pkg/parser"
	"github.com/Nostromos/cyoa/pkg/types"
)

type Model struct {
	Storybook types.Storybook
	Chapter   types.Chapter
	Title     string
	Story     []string
	Choices   []types.Option
	Cursor    int
	Selected  map[int]struct{}
}

func InitialModel() Model {

	story, err := parser.LoadStory("./gopher.json")
	if err != nil {
		panic(err)
	}

	initialChapter := story["intro"]
	initialStory := initialChapter.Story
	initialOptions := initialChapter.Options

	return Model{
		Storybook: story,
		Chapter:   initialChapter,
		Title:     initialChapter.Title,
		Story:     initialStory,
		Choices:   initialOptions,
		// cursor is nil
		Selected: make(map[int]struct{}),
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
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}

		case "enter", " ":
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

	return m, nil
}

func (m Model) View() string {
	// the story
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

			// is the choice selected?
			// checked := " "
			// if _, ok := m.Selected[i]; ok {
			// 	checked = "x"
			// }

			// render the row
			opts += fmt.Sprintf("%s %s\n", cursor, choice.Text)
		}
	}

	// footer
	opts += "\nPress q to quit.\n"

	var optionsStyle = lipgloss.NewStyle().
		Bold(true).
		Width(120).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(2).
		PaddingLeft(4).
		PaddingRight(4)

	var styledOptions = optionsStyle.Render(opts)

	var text = "\n" + StyledTitle + "\n\n" + StyledStory + "\n\n\n" + styledOptions
	return text
}
