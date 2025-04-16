package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var cursorStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("212"))

type searchForInputMsg string
type resendSearchMsg struct{}

func searchCmd(searchFor string) tea.Cmd {
	return func() tea.Msg {
		return searchForInputMsg(searchFor)
	}
}

type searchModel struct {
	input     input
	isFocused bool
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case resendSearchMsg:
		break // do nothing special
	case tea.KeyMsg:
		m.input.handleKeys(msg)
	}

	return m, searchCmd(m.input.input)
}

func (m searchModel) View() string {
	b := strings.Builder{}
	b.WriteString("Search: ")

	if len(m.input.input) == 0 {
		if m.isFocused {
			b.WriteString(cursorStyle.Render(" "))
		}
		return b.String()
	}

	before := m.input.input[:m.input.cursor]
	after := m.input.input[m.input.cursor:]

	var cursorChar string
	if len(after) != 0 {
		if m.isFocused {
			cursorChar = cursorStyle.Render(string(after[0]))
		}
	} else {
		if m.isFocused {
			cursorChar = cursorStyle.Render(" ")
		}
	}

	b.WriteString(before)
	b.WriteString(cursorChar)
	b.WriteString(after)

	return b.String()
}
