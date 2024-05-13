package tui

import (
	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table        table.Model
	input        textinput.Model
	originalRows []table.Row
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
				m.input.Focus()
			} else {
				m.input.Blur()
				m.table.Focus()
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.input.Focused() {
				m.table.SetRows(filterRows(m.originalRows, m.input.Value()))
			} else {
				cmdParts := strings.Fields(m.table.SelectedRow()[2])
			
				action := exec.Command(cmdParts[0], cmdParts[1:]...)

				action.Stdin = os.Stdin
				action.Stdout = os.Stdout
				action.Stderr = os.Stderr

				if err := action.Run(); err != nil {
					tea.Println(err)
					return m, tea.Quit
				}

				return m, tea.Quit
			}
		}
	}	

	updateInput, inputCmd := m.input.Update(msg)
	updateTable, tableCmd := m.table.Update(msg)

	m.input = updateInput
	m.table = updateTable

	return m, tea.Batch(inputCmd, tableCmd)
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.input.View() + "\nctrl+c : quit | esc : switch view | enter (table) : run cmd | enter (input): run search\n"
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

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0] 
	})

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

	m := model{table: t, input: textinput.New(), originalRows: rows}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}
	return nil
}

func filterRows(rows []table.Row, s string) []table.Row {
	filterRows := []table.Row{}

	for _, v := range rows {
		if strings.Contains(v[0], s) {
			filterRows = append(filterRows, v)
		} else if strings.Contains(v[1], s) {
			filterRows = append(filterRows, v)
		} else if strings.Contains(v[2], s) {
			filterRows = append(filterRows, v)
		}
	}

	return filterRows
}
