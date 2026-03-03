package model

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/radeqq007/duvet/internal/alert"
	"github.com/radeqq007/duvet/internal/config"
	"github.com/radeqq007/duvet/internal/filesystem"
	"github.com/radeqq007/duvet/internal/mode"
)

func (m Model) VisibleHeight() int {
	return m.Height - config.Layout.HeaderFooterSize
}

func (m Model) CurrentFile() filesystem.FileNode {
	if len(m.FileTree) == 0 {
		return filesystem.FileNode{}
	}
	return m.FileTree[m.Cursor]
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
}

func (m *Model) ShowAlert(alertType alert.AlertType, text ...string) {
	m.Alert.Type = alertType
	m.Alert.Text = strings.Join(text, " ")
	m.Mode = mode.Alert
}

func (m *Model) updatePreview() {
	if len(m.FileTree) == 0 {
		m.Preview.Path = ""
		m.Preview.Content = ""
		return
	}

	current := m.FileTree[m.Cursor]
	if current.IsDir {
		m.Preview.Path = ""
		m.Preview.Content = ""
		return
	}

	newPath := filepath.Join(m.CurPath, current.Name)
	if newPath == m.Preview.Path {
		return
	}

	res := <-filesystem.ReadFileContent(newPath)

	if res.Err != nil {
		res.Content = nil
	}

	content := string(res.Content)
	content = filesystem.Highlight(content, current.Name)

	m.Preview.Path = newPath
	m.Preview.Content = content
}

func (m *Model) getTargets() []string {
	if len(m.Selected) > 0 {
		paths := make([]string, 0, len(m.Selected))
		for path := range m.Selected {
			paths = append(paths, path)
		}
		return paths
	}
	return []string{m.CurrentFilePath()}
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
