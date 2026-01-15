package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))

	folderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")).
			Padding(0, 1)

	folderSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("24")).
				Bold(true).
				Padding(0, 1)

	imageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Padding(0, 1)

	imageSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(true).
				Padding(0, 1)

	listPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1)

	previewPaneStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.width == 0 {
		return "Loading..."
	}

	listWidth := m.width * 35 / 100
	previewWidth := m.width - listWidth - 4

	listContent := m.renderList(listWidth - 4)
	listPane := listPaneStyle.
		Width(listWidth).
		Height(m.height - 4).
		Render(listContent)

	previewContent := m.renderPreview(listWidth, previewWidth-4, m.height-6)
	previewPane := previewPaneStyle.
		Width(previewWidth).
		Height(m.height - 4).
		Render(previewContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, listPane, previewPane)
}

func (m Model) renderList(width int) string {
	var b strings.Builder

	title := titleStyle.Render("Wallpaper Selector")
	b.WriteString(title)
	b.WriteString("\n")

	path := pathStyle.Render(m.currentPath)
	b.WriteString(path)
	b.WriteString("\n\n")

	if len(m.entries) == 0 {
		b.WriteString(folderStyle.Render("(empty)"))
		b.WriteString("\n")
	} else {
		visibleHeight := m.height - 14
		if visibleHeight < 5 {
			visibleHeight = 5
		}

		start := 0
		if m.cursor >= visibleHeight {
			start = m.cursor - visibleHeight + 1
		}

		end := start + visibleHeight
		if end > len(m.entries) {
			end = len(m.entries)
		}

		for i := start; i < end; i++ {
			entry := m.entries[i]
			var displayName string
			var style, selectedStyle lipgloss.Style

			if entry.IsDir {
				displayName = "📁 " + entry.Name
				style = folderStyle
				selectedStyle = folderSelectedStyle
			} else {
				displayName = "🖼  " + entry.Name
				style = imageStyle
				selectedStyle = imageSelectedStyle
			}

			maxLen := width - 4
			if len(displayName) > maxLen {
				displayName = displayName[:maxLen-3] + "..."
			}

			if i == m.cursor {
				b.WriteString(selectedStyle.Render(displayName))
			} else {
				b.WriteString(style.Render(displayName))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	if m.statusMessage != "" {
		b.WriteString(statusStyle.Render(m.statusMessage))
		b.WriteString("\n")
	}

	help := helpStyle.Render("j/k: up/down • h/l: back/enter • enter: select • q: quit")
	b.WriteString(help)

	return b.String()
}

func (m Model) renderPreview(listWidth, previewWidth, previewHeight int) string {
	entry := m.CurrentEntry()

	// Calculate überzug position
	// x = list pane width + borders
	// y = top border
	x := listWidth + 4
	y := 3
	maxWidth := previewWidth - 2
	maxHeight := previewHeight - 2

	if entry == nil {
		if m.ueberzug != nil {
			m.ueberzug.Hide()
		}
		return "No selection"
	}

	if entry.IsDir {
		if m.ueberzug != nil {
			m.ueberzug.Hide()
		}
		return fmt.Sprintf("📁 %s\n\n(folder)", entry.Name)
	}

	if entry.IsImage {
		if m.ueberzug != nil {
			imagePath := filepath.Join(m.currentPath, entry.Name)
			m.ueberzug.Show(imagePath, x, y, maxWidth, maxHeight)
		}
		return fmt.Sprintf("🖼  %s", entry.Name)
	}

	if m.ueberzug != nil {
		m.ueberzug.Hide()
	}
	return "No preview"
}
