package main

import (
	"image/color"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	var body, help string

	// Priority: list > editor > filename input > home (mirrors the dispatch in Update).
	switch {
	case m.showListVisible:
		body = m.list.View()
		danger := helpKeyStyle.Foreground(snazzyRed).Render("ctrl+d") + " " + helpDescStyle.Render("delete")
		help = helpBar(
			[2]string{"enter", "open"},
			[2]string{"/", "filter"},
		) + helpSepStyle.Render(" · ") + danger + helpSepStyle.Render(" · ") + helpBar(
			[2]string{"esc", "back"},
		)

	case m.currentFile != nil:
		body = m.editorView()
		help = helpBar(
			[2]string{"ctrl+s", "save & close"},
			[2]string{"esc", "discard"},
		)

	case m.createFileInputVisible:
		body = m.newNotePromptView()
		help = helpBar(
			[2]string{"enter", "create"},
			[2]string{"esc", "cancel"},
		)

	default:
		body = m.homeView()
		help = helpBar(
			[2]string{"ctrl+n", "new note"},
			[2]string{"ctrl+l", "browse notes"},
			[2]string{"q", "quit"},
		)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerView(),
		"",
		body,
		"",
		help,
	)

	return tea.NewView(docStyle.Render(content))
}

// headerView renders the logo, tagline, and a right-aligned mode pill.
func (m model) headerView() string {
	left := logoStyle.Render("✦ TOTION") + " " + taglineStyle.Render("your second brain, in the terminal")

	mode, color := m.mode()
	pill := modePillStyle.Background(color).Render(mode)

	gap := m.contentWidth() - lipgloss.Width(left) - lipgloss.Width(pill)
	if gap < 1 {
		gap = 1
	}

	return left + strings.Repeat(" ", gap) + pill
}

// mode names the active screen and picks its accent color.
func (m model) mode() (string, color.Color) {
	switch {
	case m.showListVisible:
		return "NOTES", snazzyCyan
	case m.currentFile != nil:
		return "EDIT", snazzyBlue
	case m.createFileInputVisible:
		return "NEW", snazzyYellow
	default:
		return "HOME", snazzyGreen
	}
}

// homeView renders the welcome card with quick-start hints.
func (m model) homeView() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(snazzyPink).Render("✦ T O T I O N ✦")
	subtitle := taglineStyle.Render("capture thoughts at the speed of your terminal")

	hints := strings.Join([]string{
		helpKeyStyle.Render("ctrl+n") + helpDescStyle.Render("  capture a new note"),
		helpKeyStyle.Render("ctrl+l") + helpDescStyle.Render("  browse your notes"),
		helpKeyStyle.Render("q     ") + helpDescStyle.Render("  quit"),
	}, "\n")

	card := lipgloss.JoinVertical(lipgloss.Center, title, subtitle, "", hints)
	return panelStyle.Render(card)
}

// newNotePromptView renders the filename prompt in a pink-framed panel.
func (m model) newNotePromptView() string {
	prompt := lipgloss.JoinVertical(
		lipgloss.Left,
		panelTitleStyle.Foreground(snazzyPink).Render("✎ New note"),
		"",
		m.newFileInput.View(),
	)
	return panelStyle.BorderForeground(snazzyPink).Render(prompt)
}

// editorView renders the note editor in a blue-framed panel with the open
// note's name as the panel title.
func (m model) editorView() string {
	name := strings.TrimSuffix(filepath.Base(m.currentFile.Name()), ".md")

	editor := lipgloss.JoinVertical(
		lipgloss.Left,
		panelTitleStyle.Render("● "+name),
		"",
		m.noteTextArea.View(),
	)
	return panelStyle.BorderForeground(snazzyBlue).Render(editor)
}

// contentWidth is the usable width inside the document padding.
func (m model) contentWidth() int {
	if m.width == 0 {
		return 76 // sensible default before the first WindowSizeMsg
	}
	h, _ := docStyle.GetFrameSize()
	return m.width - h
}
