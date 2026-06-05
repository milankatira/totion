package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)


type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+n":
			m.createFileInputVisible = true
			m.newFileInput.Focus()
			return m, nil
		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	return m, cmd
}

func (m model) View() tea.View {
	welcome := "Welcome to Totion 🧠"

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingRight(2).
		PaddingLeft(2)

	help := "Ctrl+N: New file   Ctrl+L: List   Esc: Back/Save   Ctrl+S: Save   Ctrl+Q: Quit"
	welcome = style.Render(welcome)
	view := ""

	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}
	return tea.NewView(fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, help, view))
}

func initializedModel() model {
	// initialised new file input
	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(40)

	// Configure cursor style
	s := ti.Styles()
	s.Cursor.Color = lipgloss.Color("205")
	s.Cursor.Blink = true
	ti.SetStyles(s)

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
	}
}

func main() {
	p := tea.NewProgram(initializedModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there`s been an error: %v", err)
		os.Exit(1)
	}

}
