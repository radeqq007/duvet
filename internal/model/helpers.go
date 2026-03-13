package model

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/radeqq007/duvet/internal/alert"
	"github.com/radeqq007/duvet/internal/filesystem"
	"github.com/radeqq007/duvet/internal/git"
	"github.com/radeqq007/duvet/internal/mode"
)

func (m Model) VisibleHeight() int {
	return m.Layout.Height - m.config.Layout.HeaderFooterSize
}

func (m Model) CurrentFile() filesystem.FileNode {
	if len(m.FileTree) == 0 {
		return filesystem.FileNode{}
	}
	return m.FileTree[m.Cursor]
}

func (m Model) getParentDir() string {
	return filepath.Dir(m.CurPath)
}

func (m Model) CurrentFilePath() string {
	return filepath.Join(m.CurPath, m.CurrentFile().Name)
}

func (m *Model) getCurrentFile() filesystem.FileNode {
	return m.FileTree[m.Cursor]
}

func (m *Model) refreshFiles() {
	files, err := filesystem.GetFiles(m.CurPath)
	if err == nil {
		m.FileTree = files
	}
	m.Git = git.GetStatus(m.CurPath)
}

func (m *Model) ShowAlert(alertType alert.AlertType, text ...string) {
	m.IO.Alert.Type = alertType
	m.IO.Alert.Text = strings.Join(text, " ")
	m.IO.Mode = mode.Alert
}

func (m *Model) loadPreview() tea.Cmd {
	if len(m.FileTree) == 0 {
		return nil
	}

	current := m.FileTree[m.Cursor]
	if current.IsDir {
		return nil
	}

	newPath := filepath.Join(m.CurPath, current.Name)
	if newPath == m.Display.Preview.Path {
		return nil
	}

	return func() tea.Msg {
		content, err := filesystem.ReadFileContent(newPath)
		if err != nil {
			return PreviewLoaded{Path: newPath, Content: ""}
		}
		highlighted := filesystem.Highlight(string(content), current.Name, m.config.PreviewTheme)
		return PreviewLoaded{Path: newPath, Content: highlighted}
	}
}

func (m *Model) getTargets() []string {
	if len(m.IO.Selected) > 0 {
		paths := make([]string, 0, len(m.IO.Selected))
		for path := range m.IO.Selected {
			paths = append(paths, path)
		}
		return paths
	}
	return []string{m.CurrentFilePath()}
}

func (m *Model) parseInput() (int, string) {
	if len(m.IO.Input) == 0 {
		return 1, ""
	}

	i := 0
	for i < len(m.IO.Input) && m.IO.Input[i] >= '0' && m.IO.Input[i] <= '9' {
		i++
	}

	numStr := string(m.IO.Input[:i])
	motion := string(m.IO.Input[i:])

	if numStr == "" {
		return 1, motion
	}

	n, err := strconv.Atoi(numStr)
	if err != nil || n == 0 {
		return 1, motion
	}

	return n, motion
}

func (m *Model) clearInput() {
	m.IO.Input = []byte{}
}

func prettifyPath(path string) string {
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(path, home) {
		return "~" + strings.TrimPrefix(path, home)
	}
	return path
}

func isMediaFile(path string) bool {
	mediaExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif",
		".bmp", ".webp", ".svg", ".mp4",
		".mkv", ".avi", ".mov", ".mp3",
		".flac", ".wav", ".ogg", ".pdf",
	}

	return slices.Contains(mediaExtensions, strings.ToLower(filepath.Ext(path)))
}

func wrapLines(content string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Render(content)
}
