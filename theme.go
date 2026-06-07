package main

import (
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

// Snazzy color palette — https://github.com/sindresorhus/hyper-snazzy
var (
	snazzyBg      = lipgloss.Color("#282a36")
	snazzySurface = lipgloss.Color("#3a3d4d")
	snazzyFg      = lipgloss.Color("#eff0eb")
	snazzyRed     = lipgloss.Color("#ff5c57")
	snazzyGreen   = lipgloss.Color("#5af78e")
	snazzyYellow  = lipgloss.Color("#f3f99d")
	snazzyBlue    = lipgloss.Color("#57c7ff")
	snazzyPink    = lipgloss.Color("#ff6ac1")
	snazzyCyan    = lipgloss.Color("#9aedfe")
	snazzyMuted   = lipgloss.Color("#686868")
)

var (
	docStyle = lipgloss.NewStyle().Padding(1, 2)

	// Header
	logoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(snazzyBg).
			Background(snazzyPink).
			Padding(0, 1)

	taglineStyle = lipgloss.NewStyle().
			Foreground(snazzyMuted).
			Italic(true)

	modePillStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(snazzyBg).
			Padding(0, 1)

	// Content panels
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(snazzySurface).
			Padding(1, 2)

	panelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(snazzyCyan)

	// List title pill (kept as titleStyle — referenced from model.go)
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(snazzyBg).
			Background(snazzyPink).
			Padding(0, 1)

	// Help bar
	helpKeyStyle  = lipgloss.NewStyle().Bold(true).Foreground(snazzyPink)
	helpDescStyle = lipgloss.NewStyle().Foreground(snazzyMuted)
	helpSepStyle  = lipgloss.NewStyle().Foreground(snazzySurface)
)

// helpBar renders key/description pairs as a dotted, color-coded help line.
func helpBar(items ...[2]string) string {
	sep := helpSepStyle.Render(" · ")
	parts := make([]string, 0, len(items))
	for _, it := range items {
		parts = append(parts, helpKeyStyle.Render(it[0])+" "+helpDescStyle.Render(it[1]))
	}
	return strings.Join(parts, sep)
}

// newSnazzyDelegate returns a list delegate styled with the Snazzy palette.
func newSnazzyDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(snazzyFg)
	d.Styles.NormalDesc = d.Styles.NormalDesc.Foreground(snazzyMuted)

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(snazzyPink).
		BorderLeftForeground(snazzyPink)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(snazzyCyan).
		BorderLeftForeground(snazzyPink)

	d.Styles.DimmedTitle = d.Styles.DimmedTitle.Foreground(snazzyMuted)
	d.Styles.DimmedDesc = d.Styles.DimmedDesc.Foreground(snazzySurface)

	d.Styles.FilterMatch = d.Styles.FilterMatch.
		Foreground(snazzyYellow).
		Underline(true)

	return d
}

// newSnazzyInputStyles themes the filename text input.
func newSnazzyInputStyles() textinput.Styles {
	s := textinput.DefaultDarkStyles()

	s.Focused.Prompt = s.Focused.Prompt.Foreground(snazzyPink)
	s.Focused.Text = s.Focused.Text.Foreground(snazzyFg)
	s.Focused.Placeholder = s.Focused.Placeholder.Foreground(snazzyMuted)
	s.Blurred.Placeholder = s.Blurred.Placeholder.Foreground(snazzyMuted)

	s.Cursor.Color = snazzyPink
	s.Cursor.Blink = true

	return s
}

// newSnazzyTextareaStyles themes the note editor.
func newSnazzyTextareaStyles() textarea.Styles {
	s := textarea.DefaultDarkStyles()

	s.Focused.Text = s.Focused.Text.Foreground(snazzyFg)
	s.Focused.Placeholder = s.Focused.Placeholder.Foreground(snazzyMuted)
	s.Focused.Prompt = s.Focused.Prompt.Foreground(snazzyPink)
	s.Focused.CursorLine = s.Focused.CursorLine.Background(snazzySurface)
	s.Focused.EndOfBuffer = s.Focused.EndOfBuffer.Foreground(snazzySurface)

	s.Cursor.Color = snazzyPink
	s.Cursor.Blink = true

	return s
}

// applySnazzyListStyles themes the note list chrome (title, status, pagination).
func applySnazzyListStyles(l *list.Model) {
	l.Styles.Title = titleStyle
	l.Styles.StatusBar = l.Styles.StatusBar.Foreground(snazzyMuted)
	l.Styles.StatusBarActiveFilter = l.Styles.StatusBarActiveFilter.Foreground(snazzyCyan)
	l.Styles.StatusBarFilterCount = l.Styles.StatusBarFilterCount.Foreground(snazzyMuted)
	l.Styles.NoItems = lipgloss.NewStyle().Foreground(snazzyMuted)
	l.Styles.ActivePaginationDot = l.Styles.ActivePaginationDot.Foreground(snazzyPink)
	l.Styles.InactivePaginationDot = l.Styles.InactivePaginationDot.Foreground(snazzySurface)
	l.Styles.DividerDot = l.Styles.DividerDot.Foreground(snazzySurface)
	l.Styles.Filter = newSnazzyInputStyles()
}
