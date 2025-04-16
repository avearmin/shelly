package tui

import (
	"github.com/avearmin/shelly/internal/cmdstore"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	focusSearch = iota
	focusCommandList
	focusForm
)

var (
	boxStyle          = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	focusedBoxStyle   = boxStyle.BorderForeground(lipgloss.Color("205")) // magenta
	unfocusedBoxStyle = boxStyle.BorderForeground(lipgloss.Color("240")) // gray
)

type focusOtherMsg string

type model struct {
	search     searchModel
	actionList listModel
	form       formModel
	focus      int
	width      int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var result tea.Model

	switch msg := msg.(type) {
	case searchForInputMsg, updateActionListMsg:
		result, cmd = m.actionList.Update(msg)
		m.actionList = result.(listModel)
	case resendSearchMsg:
		result, cmd = m.search.Update(msg)
		m.search = result.(searchModel)
	case submitMsg:
		cmd = saveCmd(msg)
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
			if msg.String() == "a" {
				m.focus = focusForm
			} else {
				result, cmd = m.actionList.Update(msg)
				m.actionList = result.(listModel)
			}
		case focusForm:
			result, cmd = m.form.Update(msg)
			m.form = result.(formModel)
		}
	}
	return m, cmd
}

func (m model) View() string {
	searchStyle := unfocusedBoxStyle
	listStyle := unfocusedBoxStyle

	if m.focus == focusSearch {
		searchStyle = focusedBoxStyle
	} else {
		listStyle = focusedBoxStyle
	}

	searchBox := searchStyle.
		Width(m.width).
		Render(m.search.View())

	commandBox := listStyle.
		Width(m.width).
		Render(m.actionList.View())

	var currentView string
	if m.focus == focusForm {
		currentView = m.form.View()
	} else {
		currentView = lipgloss.JoinVertical(lipgloss.Left, searchBox, commandBox)

	}

	return currentView
}

func Start(cmds []cmdstore.Command) (cmdstore.Command, error) {
	// in the future we'll add these values to the config
	appWidth := 80
	appViewPortLength := 5

	m := model{
		search: searchModel{
			input:     input{"", 0},
			isFocused: false,
		},
		actionList: listModel{
			items:          cmds,
			filteredItems:  cmds,
			index:          0,
			cursor:         0,
			selected:       cmdstore.Command{},
			viewPortLength: appViewPortLength,
		},
		form: formModel{
			alias: input{"", 0},
			description: input{"", 0},
			action: input{"", 0},
			focus: focusFormAlias,
		},
		focus: focusCommandList,
		width: appWidth,
	}

	finalModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return cmdstore.Command{}, err
	}

	return finalModel.(model).actionList.selected, nil
}
