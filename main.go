package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	startPath := getStartPath()

	if err := validateDir(startPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nUsage: %s [directory]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  directory: Starting folder to browse\n")
		fmt.Fprintf(os.Stderr, "  (omit to start from home directory)\n")
		os.Exit(1)
	}

	// Start überzug++ daemon for image previews
	ueberzug, err := NewUeberzug()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not start ueberzug++: %v\n", err)
		fmt.Fprintf(os.Stderr, "Image previews will be disabled.\n")
		ueberzug = nil
	}

	// Ensure cleanup on exit
	defer func() {
		if ueberzug != nil {
			ueberzug.Close()
		}
	}()

	p := tea.NewProgram(
		NewModel(startPath, ueberzug),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getStartPath() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func validateDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access '%s': %v", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("'%s' is not a directory", path)
	}
	return nil
}
