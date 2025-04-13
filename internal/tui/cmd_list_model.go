package tui

import (
	"strconv"
	"strings"

	"github.com/avearmin/shelly/internal/cmdstore"
	tea "github.com/charmbracelet/bubbletea"
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
	builder := strings.Builder{}

	unselectedBar := "│ "
	selectedBar := "┃ "

	for i, v := range m.filteredItems {
		bar := unselectedBar
		if m.index == i {
			bar = selectedBar
		}
		builder.Write([]byte(bar + v.Name + "\n"))
		builder.Write([]byte(bar + v.Description + "\n"))
		builder.Write([]byte(bar + v.Action + "\n"))
		builder.Write([]byte(bar + v.LastUsedInHumanTerms() + "\n"))
		builder.Write([]byte("\n"))
	}

	return builder.String()
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
