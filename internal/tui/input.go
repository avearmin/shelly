package tui

import tea "github.com/charmbracelet/bubbletea"

type input struct {
	input  string
	cursor int
}

func (i *input) handleKeys(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyRight:
		if i.cursor != len(i.input) {
			i.cursor++
		}
	case tea.KeyLeft:
		if i.cursor != 0 {
			i.cursor--
		}
	case tea.KeyBackspace:
		if i.cursor != 0 {
			i.input = i.input[:i.cursor-1] + i.input[i.cursor:]
			i.cursor--
		}
	case tea.KeySpace:
		i.input = i.input[:i.cursor] + " " + i.input[i.cursor:]
		i.cursor++
	case tea.KeyRunes:
		for _, rune := range msg.Runes {
			i.input = i.input[:i.cursor] + string(rune) + i.input[i.cursor:]
			i.cursor++
		}
	}
}
