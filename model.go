package main

import (
	"log"
	"os"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// item is a single note entry in the list view.
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
	currentFile            *os.File
	noteTextArea           textarea.Model
	list                   list.Model
	showListVisible        bool
	width, height          int
}

func (m model) Init() tea.Cmd {
	return nil
}

func initializedModel() model {
	if err := ensureVaultDir(); err != nil {
		log.Fatalf("Error in creating vault dir: %v", err)
	}

	// new file name input
	ti := textinput.New()
	ti.Placeholder = "What would you like to call it?"
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(40)
	ti.SetStyles(newSnazzyInputStyles())

	// note editor
	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.Placeholder = "Write your note here..."
	ta.Focus()
	ta.SetStyles(newSnazzyTextareaStyles())

	// note list
	finalList := list.New(listFiles(), newSnazzyDelegate(), 0, 0)
	finalList.Title = "✦ All Notes"
	applySnazzyListStyles(&finalList)

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
		list:                   finalList,
	}
}
