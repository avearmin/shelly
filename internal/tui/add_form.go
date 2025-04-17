package tui

import (
	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	focusFormAlias    = "formalias"
	focusFormDesc     = "formdesc"
	focusFormAction   = "formaction"
	focusSubmitButton = "focussubmitbutton"
)

type submitMsg struct {
	alias       string
	description string
	action      string
}

type exitFormMsg struct{}

func saveCmd(msg submitMsg) tea.Cmd {
	return func() tea.Msg {
		config, err := configstore.Load()
		if err != nil {
			return tea.Quit
		}
		store, err := cmdstore.Load(config.CmdsPath)
		if err != nil {
			return tea.Quit
		}
		action := cmdstore.Command{
			Name:        msg.alias,
			Description: msg.description,
			Action:      msg.action,
		}

		store[msg.alias] = action
		cmdstore.Save(config.CmdsPath, store)
		return updateActionListMsg{action}
	}

}

func submitCmd(alias, description, action string) tea.Cmd {
	return func() tea.Msg {
		return submitMsg{
			alias:       alias,
			description: description,
			action:      action,
		}
	}
}

func exitFormCmd() tea.Cmd {
	return func() tea.Msg {	
		return exitFormMsg{}	
	}
}

type formModel struct {
	alias       input
	description input
	action      input
	focus       string
}

func (m formModel) Init() tea.Cmd {
	return nil
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.alias = input{"", 0}
			m.description = input{"", 0}
			m.action = input{"", 0}
			m.focus = focusFormAlias
			cmd = exitFormCmd()
		}
		switch m.focus {
		case focusFormAlias:
			switch msg.String() {
			case "up":
				m.focus = focusSubmitButton
			case "down":
				m.focus = focusFormDesc
			default:
				m.alias.handleKeys(msg)
			}
		case focusFormDesc:
			switch msg.String() {
			case "up":
				m.focus = focusFormAlias
			case "down":
				m.focus = focusFormAction
			default:
				m.description.handleKeys(msg)
			}
		case focusFormAction:
			switch msg.String() {
			case "up":
				m.focus = focusFormDesc
			case "down":
				m.focus = focusSubmitButton
			default:
				m.action.handleKeys(msg)
			}
		case focusSubmitButton:
			switch msg.String() {
			case "up":
				m.focus = focusFormAction
			case "down":
				m.focus = focusFormAlias
			case "enter":
				alias := m.alias.input
				description := m.description.input
				action := m.action.input

				m.alias = input{"", 0}
				m.description = input{"", 0}
				m.action = input{"", 0}
				m.focus = focusFormAlias

				cmd = tea.Batch(
					submitCmd(alias, description, action),
					exitFormCmd(),
				)
			}
		}
	}

	return m, cmd
}

func (m formModel) View() string {
	aliasStyle := unfocusedBoxStyle
	descriptionStyle := unfocusedBoxStyle
	actionStyle := unfocusedBoxStyle
	buttonStyle := unfocusedBoxStyle

	switch m.focus {
	case focusFormAlias:
		aliasStyle = focusedBoxStyle
	case focusFormDesc:
		descriptionStyle = focusedBoxStyle
	case focusFormAction:
		actionStyle = focusedBoxStyle
	case focusSubmitButton:
		buttonStyle = focusedBoxStyle
	}

	return aliasStyle.Render("Alias: ", m.alias.input) + "\n" +
		descriptionStyle.Render("Description: ", m.description.input) + "\n" +
		actionStyle.Render("Action: "+m.action.input) + "\n" +
		buttonStyle.Render("submit") + "\n"
}
