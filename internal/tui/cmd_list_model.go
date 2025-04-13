package tui

import (
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
	selected       int
	viewPortLength int
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "up" || msg.String() == "k" {
			m.index = mod(m.index-1, len(m.filteredItems))
			return m, nil
		}
		if msg.String() == "down" || msg.String() == "j" {
			m.index = mod(m.index+1, len(m.filteredItems))
			return m, nil
		}
		if msg.String() == "enter" {
			m.selected = m.index
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m listModel) View() string {
	b := strings.Builder{}

	for i, v := range m.filteredItems {
		var bar lipgloss.Style
		if m.index == i {
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

	return b.String()
}

func Start(cmds []cmdstore.Command) (cmdstore.Command, error) {
	model := listModel{
		items:          cmds,
		filteredItems:  cmds,
		index:          0,
		selected:       0,
		viewPortLength: 10,
	}

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return cmdstore.Command{}, err
	}

	return finalModel.(listModel).filteredItems[finalModel.(listModel).selected], nil
}

// % is not a true mod, and wont work as expected with negative numbers
func mod(a, b int) int {
	return (a%b + b) % b
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
