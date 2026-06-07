package main

import (
	"log"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-5)

		// Editor sits inside a bordered panel with a title line: subtract the
		// panel frame plus the header (title + blank line) and footer chrome.
		ph, pv := panelStyle.GetFrameSize()
		m.noteTextArea.SetWidth(max(20, msg.Width-h-ph))
		m.noteTextArea.SetHeight(max(3, msg.Height-v-pv-9))

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			if m.createFileInputVisible {
				m.createFileInputVisible = false
			}

			if m.currentFile != nil {
				m.noteTextArea.SetValue("")
				m.currentFile = nil
			}

			if m.showListVisible {
				if m.list.FilterState() == list.Filtering {
					// Fall through to the dispatch below so the list
					// receives esc and cancels its own filter.
					break
				}
				m.showListVisible = false
			}

			return m, nil

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			// Only quit from the home screen — "q" must stay typeable
			// inside the editor, the filename input, and the list filter.
			if !m.createFileInputVisible && m.currentFile == nil && !m.showListVisible {
				return m, tea.Quit
			}

		case "ctrl+s":
			return m.handleSave()

		case "ctrl+l":
			m.list.SetItems(listFiles())
			m.showListVisible = true
			return m, nil

		case "ctrl+n":
			m.createFileInputVisible = true
			m.newFileInput.Focus()
			return m, nil

		case "enter":
			if m.createFileInputVisible {
				return m.handleCreateFile()
			}
			if m.showListVisible {
				return m.handleOpenSelected()
			}

		case "ctrl+d":
			if m.showListVisible {
				return m.handleDeleteSelected()
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

	return m, cmd
}

// handleSave writes the editor content to the open note and closes it.
func (m model) handleSave() (tea.Model, tea.Cmd) {
	if m.currentFile == nil {
		return m, nil
	}

	if err := saveNote(m.currentFile, m.noteTextArea.Value()); err != nil {
		log.Fatalf("can not save the file 😢: %v", err)
	}

	m.currentFile = nil
	m.noteTextArea.SetValue("")

	return m, nil
}

// handleCreateFile creates a new note from the filename input and opens the editor.
func (m model) handleCreateFile() (tea.Model, tea.Cmd) {
	fileName := m.newFileInput.Value()
	if fileName == "" {
		return m, nil
	}

	f, err := createNote(fileName)
	if err != nil {
		log.Fatalf("Error in opening/creating file: %v", err)
	}

	m.currentFile = f
	m.createFileInputVisible = false
	m.newFileInput.SetValue("")
	m.noteTextArea.Focus()

	return m, nil
}

// handleOpenSelected loads the highlighted note into the editor.
func (m model) handleOpenSelected() (tea.Model, tea.Cmd) {
	if selected, ok := m.list.SelectedItem().(item); ok {
		content, f, err := openNote(selected.title)
		if err != nil {
			log.Fatalf("Error in opening file: %v", err)
		}
		m.noteTextArea.SetValue(content)
		m.currentFile = f
		m.noteTextArea.Focus()
	}

	m.showListVisible = false
	return m, nil
}

// handleDeleteSelected removes the highlighted note and refreshes the list.
func (m model) handleDeleteSelected() (tea.Model, tea.Cmd) {
	if selected, ok := m.list.SelectedItem().(item); ok {
		if err := deleteNote(selected.title); err != nil {
			log.Fatalf("Error in deleting file: %v", err)
		}
	}

	m.list.SetItems(listFiles())
	return m, nil
}
