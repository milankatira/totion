package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("205")).
			Padding(0, 1)

	welcomeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("205")).
			PaddingRight(2).
			PaddingLeft(2)
)

const helpText = "Ctrl+N: New file   Ctrl+L: List   Esc: Back   Ctrl+S: Save   Ctrl+D: Delete   Q: Quit"

func (m model) View() tea.View {
	welcome := welcomeStyle.Render("Welcome to Totion 🧠")

	// Priority: list > editor > filename input (mirrors the dispatch in Update).
	view := ""
	switch {
	case m.showListVisible:
		view = m.list.View()
	case m.currentFile != nil:
		view = m.noteTextArea.View()
	case m.createFileInputVisible:
		view = m.newFileInput.View()
	}

	return tea.NewView(fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, helpText, view))
}
