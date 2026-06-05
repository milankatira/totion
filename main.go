package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	msg string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			fmt.Println("user clicked", msg.String())
			return m, tea.Quit

			// // The "up" and "k" keys move the cursor up
			// case "up", "k":
			// 	if m.cursor > 0 {
			// 		m.cursor--
			// 	}

			// // The "down" and "j" keys move the cursor down
			// case "down", "j":
			// 	if m.cursor < len(m.choices)-1 {
			// 		m.cursor++
			// 	}

			// // The enter key and the space bar toggle the selected state for the
			// // item that the cursor is pointing at.
			// case "enter", "space":
			// 	_, ok := m.selected[m.cursor]
			// 	if ok {
			// 		delete(m.selected, m.cursor)
			// 	} else {
			// 		m.selected[m.cursor] = struct{}{}
			// 	}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
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
	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, help, view)
}

func initializedModel() model {
	return model{
		msg: "Hii",
	}
}

func main() {
	p := tea.NewProgram(initializedModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there`s been an error: %v", err)
		os.Exit(1)
	}

}
