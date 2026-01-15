package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Entry struct {
	Name    string
	IsDir   bool
	IsImage bool
}

type Model struct {
	currentPath   string
	entries       []Entry
	cursor        int
	width         int
	height        int
	statusMessage string
	quitting      bool
	ueberzug      *Ueberzug
}

type entriesLoadedMsg struct {
	entries []Entry
}

type wallpaperSetMsg struct {
	success bool
	message string
}

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".gif": true, ".bmp": true, ".webp": true,
}

func NewModel(startPath string, ueberzug *Ueberzug) Model {
	return Model{
		currentPath: startPath,
		cursor:      0,
		ueberzug:    ueberzug,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadEntries
}

func (m Model) loadEntries() tea.Msg {
	var entries []Entry

	dirEntries, err := os.ReadDir(m.currentPath)
	if err != nil {
		return entriesLoadedMsg{entries: nil}
	}

	var folders []Entry
	var files []Entry

	for _, entry := range dirEntries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		if entry.IsDir() {
			folders = append(folders, Entry{
				Name:  name,
				IsDir: true,
			})
		} else {
			ext := strings.ToLower(filepath.Ext(name))
			if imageExts[ext] {
				files = append(files, Entry{
					Name:    name,
					IsImage: true,
				})
			}
		}
	}

	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Name < folders[j].Name
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	entries = append(folders, files...)
	return entriesLoadedMsg{entries: entries}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			if m.ueberzug != nil {
				m.ueberzug.Close()
			}
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}

		case "enter", "l", "right":
			return m.handleEnter()

		case "backspace", "h", "left":
			return m.goToParent()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case entriesLoadedMsg:
		m.entries = msg.entries

	case wallpaperSetMsg:
		m.statusMessage = msg.message
	}

	return m, nil
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if len(m.entries) == 0 || m.cursor >= len(m.entries) {
		return m, nil
	}

	entry := m.entries[m.cursor]

	if entry.IsDir {
		m.currentPath = filepath.Join(m.currentPath, entry.Name)
		m.cursor = 0
		m.statusMessage = ""
		return m, m.loadEntries
	}

	if entry.IsImage {
		wallpaperPath := filepath.Join(m.currentPath, entry.Name)
		return m, func() tea.Msg {
			success, message := SetWallpaper(wallpaperPath)
			return wallpaperSetMsg{success: success, message: message}
		}
	}

	return m, nil
}

func (m Model) goToParent() (tea.Model, tea.Cmd) {
	parent := filepath.Dir(m.currentPath)
	if parent == m.currentPath {
		return m, nil
	}
	m.currentPath = parent
	m.cursor = 0
	m.statusMessage = ""
	return m, m.loadEntries
}

func (m Model) CurrentEntry() *Entry {
	if len(m.entries) == 0 || m.cursor >= len(m.entries) {
		return nil
	}
	return &m.entries[m.cursor]
}
