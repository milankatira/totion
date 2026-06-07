package main

import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	vaultDir string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error in getting home dir: %v", err)
	}
	vaultDir = fmt.Sprintf("%s/.totion", homeDir)
}

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
	currentFile            *os.File
	noteTextArea           textarea.Model
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
		case "enter":
			if m.createFileInputVisible {
				fileName := m.newFileInput.Value()
				if fileName != "" {
					filepath := fmt.Sprintf("%s/%s.md", vaultDir, fileName)

					f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
					if err != nil {
						log.Fatalf("Error in opening/creating file: %v", err)
					}

					m.currentFile = f
					m.createFileInputVisible = false
					m.newFileInput.SetValue("")
					m.noteTextArea.Focus()
				}
				return m, nil
			}
		}

	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
		return m, cmd
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
		return m, cmd
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

	// open a text area on enter
	if m.currentFile != nil {
		view = m.noteTextArea.View()
	}

	return tea.NewView(fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, help, view))
}

func initializedModel() model {

	// initialised home dir
	err := os.MkdirAll(vaultDir, 0755)
	if err != nil {
		log.Fatalf("Error in creating vault dir: %v", err)
	}

	// initialised new file input
	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(40)

	// initialised text area
	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.Placeholder = "Write your note here..."
	ta.Focus()

	// Configure cursor style
	s := ta.Styles()
	s.Cursor.Color = lipgloss.Color("205")
	s.Cursor.Blink = true
	ta.SetStyles(s)

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
	}
}

func main() {
	p := tea.NewProgram(initializedModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there`s been an error: %v", err)
		os.Exit(1)
	}

}
