package main

import (
	"os"
	"testing"
)

// useTempVault points vaultDir at a temp directory for the duration of a test.
func useTempVault(t *testing.T) {
	t.Helper()
	original := vaultDir
	vaultDir = t.TempDir()
	t.Cleanup(func() { vaultDir = original })
}

func TestCreateSaveAndOpenNote(t *testing.T) {
	useTempVault(t)

	f, err := createNote("hello")
	if err != nil {
		t.Fatalf("createNote: %v", err)
	}

	if err := saveNote(f, "some content"); err != nil {
		t.Fatalf("saveNote: %v", err)
	}

	content, f, err := openNote("hello.md")
	if err != nil {
		t.Fatalf("openNote: %v", err)
	}
	defer f.Close()

	if content != "some content" {
		t.Errorf("openNote content = %q, want %q", content, "some content")
	}
}

func TestSaveNoteOverwritesLongerContent(t *testing.T) {
	useTempVault(t)

	f, err := createNote("note")
	if err != nil {
		t.Fatalf("createNote: %v", err)
	}
	if err := saveNote(f, "a much longer first draft"); err != nil {
		t.Fatalf("saveNote: %v", err)
	}

	_, f, err = openNote("note.md")
	if err != nil {
		t.Fatalf("openNote: %v", err)
	}
	if err := saveNote(f, "short"); err != nil {
		t.Fatalf("saveNote: %v", err)
	}

	content, f, err := openNote("note.md")
	if err != nil {
		t.Fatalf("openNote: %v", err)
	}
	defer f.Close()

	if content != "short" {
		t.Errorf("stale bytes left after overwrite: got %q, want %q", content, "short")
	}
}

func TestDeleteNote(t *testing.T) {
	useTempVault(t)

	f, err := createNote("doomed")
	if err != nil {
		t.Fatalf("createNote: %v", err)
	}
	f.Close()

	if err := deleteNote("doomed.md"); err != nil {
		t.Fatalf("deleteNote: %v", err)
	}

	if _, err := os.Stat(notePath("doomed.md")); !os.IsNotExist(err) {
		t.Errorf("note still exists after delete (stat err: %v)", err)
	}
}

func TestListFilesSkipsDirectories(t *testing.T) {
	useTempVault(t)

	f, err := createNote("visible")
	if err != nil {
		t.Fatalf("createNote: %v", err)
	}
	f.Close()

	if err := os.Mkdir(notePath("subdir"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	items := listFiles()
	if len(items) != 1 {
		t.Fatalf("listFiles returned %d items, want 1", len(items))
	}

	note, ok := items[0].(item)
	if !ok {
		t.Fatalf("listFiles item has unexpected type %T", items[0])
	}
	if note.title != "visible.md" {
		t.Errorf("listFiles title = %q, want %q", note.title, "visible.md")
	}
}
