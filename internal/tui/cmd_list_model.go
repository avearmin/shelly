package tui

import (
	"fmt"
	"strings"

	"github.com/avearmin/shelly/internal/cmdstore"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		SetString("┃")

	unselectedBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		SetString("│")

	aliasStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229")) // e.g., bright yellow

	descriptionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")) // dim gray

	commandStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")) // cyan

	lastUsedNeverStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")) // gray

	lastUsedDaysAgoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("34")) // bright green

	lastUsedWeeksAgoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("70")) // duller green

	lastUsedMonthsAgoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("228")) // yellow

	lastUsedYearsAgoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")) // red

)

type listModel struct {
	items          []cmdstore.Command
	filteredItems  []cmdstore.Command
	index          int
	cursor         int
	selected       cmdstore.Command
	viewPortLength int
	viewPortStart  int
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "up" || msg.String() == "k" {
			if m.viewPortStart == m.cursor {
				if m.index == 0 {
					m.cursor = m.viewPortLength - 1
					m.viewPortStart = len(m.filteredItems) - m.viewPortLength
				} else {
					m.viewPortStart--
				}
			} else {
				m.cursor--
			}

			m.index = mod(m.index-1, len(m.filteredItems))

			return m, nil
		}
		if msg.String() == "down" || msg.String() == "j" {
			if m.viewPortLength - 1 == m.cursor {
				if m.index == len(m.filteredItems) - 1 {
					m.cursor = 0
					m.viewPortStart = 0
				} else {
					m.viewPortStart++
				}
			} else {
				m.cursor++
			}

			m.index = mod(m.index+1, len(m.filteredItems))
			return m, nil
		}
		// here we return an empty model to clear the tui
		// to reduce screen clutter when the user exits
		if msg.String() == "enter" {
			m.selected = m.filteredItems[m.index]
			return listModel{selected: m.selected}, tea.Quit
		}
		if msg.String() == "ctrl+c" {
			return listModel{}, tea.Quit
		}
	}

	return m, nil
}

func (m listModel) View() string {
	b := strings.Builder{}

	for i, v := range m.viewportItems(){
		var bar lipgloss.Style
		if m.cursor == i {
			bar = selectedBar
		} else {
			bar = unselectedBar
		}

		b.WriteString(bar.Render(" ") + aliasStyle.Render(v.Name) + "\n")
		b.WriteString(bar.Render(" ") + descriptionStyle.Render(v.Description) + "\n")
		b.WriteString(bar.Render(" ") + commandStyle.Render(v.Action) + "\n")
		b.WriteString(bar.Render(" ") + lastUsedStyleFor(v.LastUsedInHumanTerms()) + "\n")
		b.WriteString("\n")
	}
	b.WriteString(fmt.Sprintf("%d/%d", m.index + 1, len(m.filteredItems)))

	return b.String()
}

func (m listModel) viewportItems() []cmdstore.Command {
	viewportEnd := min(len(m.filteredItems), m.viewPortStart+m.viewPortLength)
	return m.filteredItems[m.viewPortStart:viewportEnd] 
}

func Start(cmds []cmdstore.Command) (cmdstore.Command, error) {
	model := listModel{
		items:          cmds,
		filteredItems:  cmds,
		index:          0,
		cursor:         0,
		selected:       cmdstore.Command{},
		viewPortLength: 5,
	}

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return cmdstore.Command{}, err
	}

	return finalModel.(listModel).selected, nil
}

// % is not a true mod, and wont work as expected with negative numbers
func mod(a, b int) int {
	return (a%b + b) % b
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func lastUsedStyleFor(s string) string {
	switch {
	case s == "Never":
		return lastUsedNeverStyle.Render(s)
	case strings.HasSuffix(s, "day ago") || strings.HasSuffix(s, "days ago") || s == "Today" || s == "Yesterday":
		return lastUsedDaysAgoStyle.Render(s)
	case strings.HasSuffix(s, "week ago") || strings.HasSuffix(s, "weeks ago"):
		return lastUsedWeeksAgoStyle.Render(s)
	case strings.HasSuffix(s, "month ago") || strings.HasSuffix(s, "months ago"):
		return lastUsedMonthsAgoStyle.Render(s)
	case strings.HasSuffix(s, "year ago") || strings.HasSuffix(s, "years ago"):
		return lastUsedYearsAgoStyle.Render(s)
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render(s) // fallback gray
	}
}
