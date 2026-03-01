package model

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/radeqq007/duvet/internal/alert"
	"github.com/radeqq007/duvet/internal/command"
	"github.com/radeqq007/duvet/internal/config"
)

func (m *Model) handleCommand(msg command.Msg) (tea.Model, tea.Cmd) {
	switch msg.Name {

	case "q", "quit":
		return m, tea.Quit

	case "rename":
		return m.rename(msg.Args)

	case "delete":
		return m.delete(msg.Args)

	case "touch":
		return m.touch(msg.Args)

	case "mkdir":
		return m.mkdir(msg.Args)

	case "cd":
		return m.cd(msg.Args)

	case "alert":
		return m.alertCommand(msg.Args)

	case "bm":
		return m.bookmark(msg.Args)

	case "find":
		return m.find(msg.Args)

	case "select":
		return m.selectFiles(msg.Args)

	case "deselect":
		return m.deselectFiles(msg.Args)

	default:
		if strings.HasPrefix(msg.Name, "!") {
			return m.execCommand(msg.Name[1:], msg.Args)
		}
	}

	return m, nil
}

func (m *Model) rename(args []string) (tea.Model, tea.Cmd) {
	file := m.getCurrentFile()

	if len(args) < 1 {
		return m, nil
	}

	oldPath := filepath.Join(m.CurPath, file.Name)
	newPath := filepath.Join(m.CurPath, args[0])

	err := os.Rename(oldPath, newPath)
	if err != nil {
		m.ShowAlert(alert.Error, "Error renaming: ", err.Error())
	}

	m.refreshFiles()

	return m, nil
}

func (m *Model) delete(args []string) (tea.Model, tea.Cmd) {
	files := m.getTargets()

	for _, file := range files {
		err := os.RemoveAll(file)
		if err != nil {
			m.ShowAlert(alert.Error, "Error removing a file: ", err.Error())
		}
	}

	m.refreshFiles()

	return m, nil
}

func (m *Model) touch(args []string) (tea.Model, tea.Cmd) {
	if len(args) < 1 {
		return m, nil
	}

	path := filepath.Join(m.CurPath, args[0])

	_, err := os.Create(path)
	if err != nil {
		m.ShowAlert(alert.Error, "Error creating a file: ", err.Error())
	}

	m.refreshFiles()

	return m, nil
}

func (m *Model) mkdir(args []string) (tea.Model, tea.Cmd) {
	if len(args) < 1 {
		return m, nil
	}

	path := filepath.Join(m.CurPath, args[0])

	err := os.Mkdir(path, os.FileMode(os.O_CREATE))
	if err != nil {
		m.ShowAlert(alert.Error, "Error creating a directory: ", err.Error())
	}

	m.refreshFiles()

	return m, nil
}

func (m *Model) cd(args []string) (tea.Model, tea.Cmd) {
	m.LeftScroll = 0
	m.RightScroll = 0
	m.Cursor = 0

	var target string

	if len(args) == 0 {
		home, _ := os.UserHomeDir()
		target = home
	} else {
		target = args[0]

		if strings.HasPrefix(target, "~") {
			home, _ := os.UserHomeDir()

			target = filepath.Join(home, strings.TrimPrefix(target, "~"))
		}

		if !filepath.IsAbs(target) {
			target = filepath.Join(m.CurPath, target)
		}
	}

	target = filepath.Clean(target)
	info, err := os.Stat(target)
	if err != nil || !info.IsDir() {
		return m, nil
	}

	m.CurPath = target
	m.ParentDir = filepath.Dir(target)

	m.refreshFiles()

	return m, nil
}

func (m *Model) bookmark(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		return m, nil
	}

	switch args[0] {
	case "save":
		if len(args) < 2 {
			return m, nil
		}

		name := args[1]
		err := config.SetBookmark(name, m.CurPath)
		if err != nil {
			m.ShowAlert(alert.Error, "Error saving the bookmark:", err.Error())
		}

	case "load":
		if len(args) < 2 {
			return m, nil
		}

		name := args[1]

		if path, ok := config.GetBookmark(name); !ok {
			m.ShowAlert(alert.Error, "No bookmark '"+name+"' found.")
		} else {
			m.CurPath = path
			m.ParentDir = filepath.Dir(path)
			m.refreshFiles()
		}

	case "list":
		var text strings.Builder
		text.WriteString("Bookmark list:\n")
		for name, path := range config.GetBookmarks() {
			fmt.Fprintf(&text, "%s: %s\n", name, path)
		}
		m.ShowAlert(alert.Info, text.String())

	case "delete":
		if len(args) < 2 {
			return m, nil
		}
		name := args[1]
		err := config.DeleteBookmark(name)
		if err != nil {
			m.ShowAlert(alert.Error, "Error deleting the bookmark:", err.Error())
		}
	}

	return m, nil
}

func (m *Model) alertCommand(args []string) (tea.Model, tea.Cmd) {
	if len(args) < 1 {
		return m, nil
	}
	switch args[0] {
	case "normal":
		m.ShowAlert(alert.Normal, strings.Join(args[1:], " "))
	case "info":
		m.ShowAlert(alert.Info, strings.Join(args[1:], " "))
	case "warning":
		m.ShowAlert(alert.Warning, strings.Join(args[1:], " "))
	case "error":
		m.ShowAlert(alert.Error, strings.Join(args[1:], " "))
	default:
		m.ShowAlert(alert.Normal, strings.Join(args, " "))

	}
	return m, nil
}

func (m *Model) execCommand(name string, args []string) (tea.Model, tea.Cmd) {
	if name == "" {
		return m, nil
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = m.CurPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.ShowAlert(alert.Error, "Shell error: ", err.Error())
		return m, nil
	}

	if len(output) > 0 {
		m.ShowAlert(alert.Info, string(output))
	}

	m.refreshFiles()

	return m, nil
}

func (m *Model) find(args []string) (tea.Model, tea.Cmd) {
	if len(args) < 1 {
		return m, nil
	}

	query := strings.Join(args, " ")

	names := make([]string, len(m.FileTree))
	for i, f := range m.FileTree {
		names[i] = f.Name
	}

	matches := fuzzy.RankFindFold(query, names)
	if len(matches) == 0 {
		m.ShowAlert(alert.Normal, "No match found for: "+query)
		return m, nil
	}

	sort.Sort(matches)
	bestMatch := matches[0].Target

	for i, f := range m.FileTree {
		if f.Name == bestMatch {
			m.Cursor = i
			visibleHeight := m.VisibleHeight()

			if m.Cursor < m.LeftScroll {
				m.LeftScroll = m.Cursor
			} else if m.Cursor >= m.LeftScroll+visibleHeight {
				m.LeftScroll = m.Cursor - visibleHeight + 1
			}

			m.updatePreview()
			break
		}
	}

	return m, nil
}

func (m *Model) selectFiles(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		return m, nil
	}

	pattern := strings.Join(args, " ")

	for _, f := range m.FileTree {
		matched, _ := filepath.Match(pattern, f.Name)
		if matched {
			m.Selected[filepath.Join(m.CurPath, f.Name)] = struct{}{}
		}
	}

	return m, nil
}

func (m *Model) deselectFiles(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		return m, nil
	}

	pattern := strings.Join(args, " ")

	for _, f := range m.FileTree {
		matched, _ := filepath.Match(pattern, f.Name)
		if matched {
			delete(m.Selected, filepath.Join(m.CurPath, f.Name))
		}
	}

	return m, nil
}
