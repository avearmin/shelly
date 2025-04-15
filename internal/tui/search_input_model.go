package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var cursorStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("212"))

type searchForInputMsg string

type searchModel struct {
	input     string
	cursor    int
	isFocused bool
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRight:
			if m.cursor != len(m.input) {
				m.cursor++
			}
		case tea.KeyLeft:
			if m.cursor != 0 {
				m.cursor--
			}
		case tea.KeyBackspace:
			if m.cursor != 0 {
				m.input = m.input[:m.cursor-1] + m.input[m.cursor:]
				m.cursor--
			}
		case tea.KeySpace:
			m.input = m.input[:m.cursor] + " " + m.input[m.cursor:]
			m.cursor++
		case tea.KeyRunes:
			for _, rune := range msg.Runes {
				m.input = m.input[:m.cursor] + string(rune) + m.input[m.cursor:]
				m.cursor++
			}
		}
	}

	cmd := func() tea.Msg {
		return searchForInputMsg(m.input)
	}
	return m, cmd
}

func (m searchModel) View() string {
	b := strings.Builder{}
	b.WriteString("Search: ")

	if len(m.input) == 0 {
		if m.isFocused {
			b.WriteString(cursorStyle.Render(" "))
		}
		return b.String()
	}

	before := m.input[:m.cursor]
	after := m.input[m.cursor:]

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
