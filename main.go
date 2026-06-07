package main

import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	vaultDir string
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

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
	list                   list.Model
	showListVisible        bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-5)

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "esc":
			if m.createFileInputVisible {
				m.createFileInputVisible = false
			}

			if m.currentFile != nil {
				m.currentFile = nil
			}

			if m.showListVisible {
				if m.list.FilterState() == list.Filtering {
					break
				}
				m.showListVisible = false
			}

			return m, nil

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+s":
			// text area value ->write in that file decsrptior and close it
			content := m.noteTextArea.Value()
			if m.currentFile == nil {
				break
			}

			if err := m.currentFile.Truncate(0); err != nil {
				log.Fatalf("can not save the file 😢")
				return m, nil
			}
			if _, err := m.currentFile.Seek(0, 0); err != nil {
				log.Fatalf("can not save the file 😢")
				return m, nil
			}

			if _, err := m.currentFile.WriteString(content); err != nil {
				log.Fatalf("can not save the file 😢")
				return m, nil
			}

			if err := m.currentFile.Close(); err != nil {
				log.Fatalf("can not save the file 😢")
			}

			m.currentFile = nil
			m.noteTextArea.SetValue("")

			return m, nil

		case "ctrl+l":
			noteList := listFiles()
			m.list.SetItems(noteList)
			m.showListVisible = true
			return m, nil
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

			if m.showListVisible {
				item, ok := m.list.SelectedItem().(item)

				if ok {
					filepath := fmt.Sprintf("%s/%s", vaultDir, item.title)
					content, err := os.ReadFile(filepath)
					if err != nil {
						log.Fatalf("Error in reading file: %v", err)
						return m, nil
					}
					m.noteTextArea.SetValue(string(content))
					m.currentFile, err = os.OpenFile(filepath, os.O_RDWR, 0644)
					if err != nil {
						log.Fatalf("Error in opening file: %v", err)
					}
					m.noteTextArea.Focus()
				}
				m.showListVisible = false
				return m, nil
			}

		}

	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	if m.showListVisible {
		m.list, cmd = m.list.Update(msg)
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

	help := "Ctrl+N: New file   Ctrl+L: List   Esc: Back   Ctrl+S: Save   Ctrl+Q: Quit"
	welcome = style.Render(welcome)
	view := ""

	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}

	// open a text area on enter
	if m.currentFile != nil {
		view = m.noteTextArea.View()
	}

	if m.showListVisible {
		view = m.list.View()
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

	// list
	noteList := listFiles()
	finalList := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "All Notes 📕"
	finalList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("205")).Padding(0, 1)

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
		list:                   finalList,
	}
}

func main() {
	p := tea.NewProgram(initializedModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there`s been an error: %v", err)
		os.Exit(1)
	}

}

func listFiles() []list.Item {
	items := make([]list.Item, 0)
	files, err := os.ReadDir(vaultDir)
	if err != nil {
		log.Fatalf("Error in reading Notes 😢")
	}

	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				continue
			}
			modTime := info.ModTime().Format("2026-11-02 15:04")
			items = append(items, item{
				title: file.Name(),
				desc:  fmt.Sprintf("Modified: %s", modTime),
			})

		}
	}

	return items
}
