package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

// version is overridden at release time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v", "version":
			fmt.Printf("totion %s\n", version)
			return
		}
	}

	p := tea.NewProgram(initializedModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there`s been an error: %v", err)
		os.Exit(1)
	}
}
