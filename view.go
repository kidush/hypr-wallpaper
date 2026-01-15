package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Layout constants
const (
	// ListPaneWidthPercent is the percentage of terminal width for the file list
	ListPaneWidthPercent = 35

	// PaneBorderWidth accounts for borders between panes
	PaneBorderWidth = 4

	// PaneContentPadding is the padding inside each pane for content
	PaneContentPadding = 4

	// PaneHeightOffset accounts for top/bottom borders and margins
	PaneHeightOffset = 4

	// PreviewHeightOffset accounts for preview pane header and borders
	PreviewHeightOffset = 6

	// ListHeaderOffset accounts for title, path, and spacing in file list
	ListHeaderOffset = 14

	// MinVisibleItems is the minimum number of items to show in the list
	MinVisibleItems = 5

	// TruncationSuffix length for "..."
	TruncationSuffixLen = 3

	// PreviewXOffset is the X position offset for überzug preview
	PreviewXOffset = 4

	// PreviewYOffset is the Y position offset for überzug preview
	PreviewYOffset = 3

	// PreviewPadding is internal padding for the preview image
	PreviewPadding = 2
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

	listWidth := m.width * ListPaneWidthPercent / 100
	previewWidth := m.width - listWidth - PaneBorderWidth

	listContent := m.renderList(listWidth - PaneContentPadding)
	listPane := listPaneStyle.
		Width(listWidth).
		Height(m.height - PaneHeightOffset).
		Render(listContent)

	previewContent := m.renderPreview(listWidth, previewWidth-PaneContentPadding, m.height-PreviewHeightOffset)
	previewPane := previewPaneStyle.
		Width(previewWidth).
		Height(m.height - PaneHeightOffset).
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
		visibleHeight := m.height - ListHeaderOffset
		if visibleHeight < MinVisibleItems {
			visibleHeight = MinVisibleItems
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

			maxLen := width - PaneContentPadding
			if len(displayName) > maxLen {
				displayName = displayName[:maxLen-TruncationSuffixLen] + "..."
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
	x := listWidth + PreviewXOffset
	y := PreviewYOffset
	maxWidth := previewWidth - PreviewPadding
	maxHeight := previewHeight - PreviewPadding

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
