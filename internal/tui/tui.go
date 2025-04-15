package tui

import (
	"github.com/avearmin/shelly/internal/cmdstore"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	focusSearch = iota
	focusCommandList
)

type focusOtherMsg string

type model struct {
	search      searchModel
	commandList listModel
	focus       int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var result tea.Model

	switch msg := msg.(type) {
	case searchForInputMsg:
		result, cmd = m.commandList.Update(msg)
		m.commandList = result.(listModel)
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return model{}, tea.Quit
		}
		if msg.String() == "esc" {
			if m.focus == focusSearch {
				m.focus = focusCommandList
				m.search.isFocused = false
			} else if m.focus == focusCommandList {
				m.focus = focusSearch
				m.search.isFocused = true
			}
		}

		switch m.focus {
		case focusSearch:
			result, cmd = m.search.Update(msg)
			m.search = result.(searchModel)
		case focusCommandList:
			result, cmd = m.commandList.Update(msg)
			m.commandList = result.(listModel)
		}
	}
	return m, cmd
}

func (m model) View() string {
	return m.search.View() + "\n\n" + m.commandList.View()
}

func Start(cmds []cmdstore.Command) (cmdstore.Command, error) {
	m := model{
		search: searchModel{
			input:  "",
			cursor: 0,
		},
		commandList: listModel{
			items:          cmds,
			filteredItems:  cmds,
			index:          0,
			cursor:         0,
			selected:       cmdstore.Command{},
			viewPortLength: 5,
		},
		focus: focusCommandList,
	}

	finalModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return cmdstore.Command{}, err
	}

	return finalModel.(model).commandList.selected, nil
}
