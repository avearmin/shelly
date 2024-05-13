package tui

import (
	"fmt"
	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			cmdParts := strings.Fields(m.table.SelectedRow()[2])

			action := exec.Command(cmdParts[0], cmdParts[1:]...)

			action.Stdin = os.Stdin
			action.Stdout = os.Stdout
			action.Stderr = os.Stderr

			if err := action.Run(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return m, tea.Quit
			}

			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func Start() error {
	columns := []table.Column{
		{Title: "Alias", Width: 10},
		{Title: "Description", Width: 50},
		{Title: "Command", Width: 30},
	}

	rows := []table.Row{}

	cmdsPath, err := configstore.GetCmdsPath()
	if err != nil {
		return err
	}

	cmds, err := cmdstore.Load(cmdsPath)
	if err != nil {
		return err
	}

	for _, v := range cmds {
		rows = append(rows, table.Row{v.Name, v.Description, v.Action})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}
	return nil
}
