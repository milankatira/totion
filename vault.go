package main

import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/list"
)

var vaultDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error in getting home dir: %v", err)
	}
	vaultDir = fmt.Sprintf("%s/.totion", homeDir)
}

// ensureVaultDir creates the vault directory if it does not exist yet.
func ensureVaultDir() error {
	return os.MkdirAll(vaultDir, 0755)
}

// notePath builds the absolute path for a note file inside the vault.
func notePath(name string) string {
	return fmt.Sprintf("%s/%s", vaultDir, name)
}

// createNote opens (or creates) a note with the given name and returns its handle.
func createNote(name string) (*os.File, error) {
	f, err := os.OpenFile(notePath(name+".md"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create note %q: %w", name, err)
	}
	return f, nil
}

// openNote reads an existing note's content and returns it along with an open handle.
func openNote(name string) (string, *os.File, error) {
	path := notePath(name)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read note %q: %w", name, err)
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open note %q: %w", name, err)
	}

	return string(content), f, nil
}

// saveNote overwrites the open note file with content and closes the handle.
func saveNote(f *os.File, content string) error {
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate note: %w", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek note: %w", err)
	}
	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write note: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close note: %w", err)
	}
	return nil
}

// deleteNote removes a note file from the vault.
func deleteNote(name string) error {
	if err := os.Remove(notePath(name)); err != nil {
		return fmt.Errorf("failed to delete note %q: %w", name, err)
	}
	return nil
}

// listFiles returns every note in the vault as a list item.
func listFiles() []list.Item {
	items := make([]list.Item, 0)
	files, err := os.ReadDir(vaultDir)
	if err != nil {
		log.Fatalf("Error in reading Notes 😢")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		modTime := info.ModTime().Format("2006-01-02 15:04")
		items = append(items, item{
			title: file.Name(),
			desc:  fmt.Sprintf("Modified: %s", modTime),
		})
	}

	return items
}
