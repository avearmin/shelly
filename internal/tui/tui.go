package tui

import (
	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
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

// This func should only be called from inside the TUI.
// The TUI should not be accessible if the config/cmdstore
// does not exist. If we get here without them a panic is needed.
func deleteFromStoreCmd(alias string) tea.Cmd {
	return func() tea.Msg {
		config, err := configstore.Load()
		if err != nil {
			panic(err)
		}

		actions, err := cmdstore.Load(config.CmdsPath)
		if err != nil {
			panic(err)
		}

		delete(actions, alias)

		if err := cmdstore.Save(config.CmdsPath, actions); err != nil {
			panic(err)
		}

		return nil
	}
}

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
	case editMsg:
		m.focus = focusForm
		m.form.previous = msg.payload.Name
		m.form.alias.input = msg.payload.Name
		m.form.description.input = msg.payload.Description
		m.form.action.input = msg.payload.Action
		m.form.focus = focusFormAlias
	case delFromItems:
		result, cmd = m.actionList.Update(msg)
		m.actionList = result.(listModel)
	case delFromStoreMsg:
		cmd = deleteFromStoreCmd(msg.alias)
	case searchForInputMsg, updateActionListMsg:
		result, cmd = m.actionList.Update(msg)
		m.actionList = result.(listModel)
	case resendSearchMsg:
		result, cmd = m.search.Update(msg)
		m.search = result.(searchModel)
	case submitMsg:
		cmd = saveCmd(msg)
	case exitFormMsg:
		m.focus = focusCommandList
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
			items:             cmds,
			filteredItems:     filter("", cmds),
			filterIndex:       0,
			cursor:            0,
			selected:          cmdstore.Command{},
			viewPortLengthMax: appViewPortLength,
			viewPortLengthCur: min(len(cmds), appViewPortLength),
		},
		form: formModel{
			previous:    "",
			alias:       input{"", 0},
			description: input{"", 0},
			action:      input{"", 0},
			focus:       focusFormAlias,
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
