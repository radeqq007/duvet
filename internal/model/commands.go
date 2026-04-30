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
	"github.com/radeqq007/duvet/internal/bookmarks"
	"github.com/radeqq007/duvet/internal/command"
	"github.com/radeqq007/duvet/internal/filesystem"
)

func (m *Model) handleCommand(msg command.Msg) (tea.Model, tea.Cmd) {
	switch msg.Name {

	case "q", "quit":
		return m, tea.Quit

	case "rename":
		return m.rename(msg.Args)

	case "delete":
		return m.delete()

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

	case "yank":
		return m.yank()

	case "paste":
		return m.paste()

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

func (m *Model) delete() (tea.Model, tea.Cmd) {
	files := m.getTargets()

	for _, file := range files {
		err := os.RemoveAll(file)
		if err != nil {
			m.ShowAlert(alert.Error, "Error removing a file: ", err.Error())
		}
	}

	m.refreshFiles()

	// ensure that the cursor isn't out of bounds
	m.Cursor = min(m.Cursor, len(m.FileTree)-1)

	return m, nil
}

func (m *Model) touch(args []string) (tea.Model, tea.Cmd) {
	if len(args) < 1 {
		return m, nil
	}

	path := filepath.Join(m.CurPath, args[0])

	err := filesystem.CreateFile(path)
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

	err := filesystem.CreateDirectory(path)
	if err != nil {
		m.ShowAlert(alert.Error, "Error creating a directory: ", err.Error())
	}

	m.refreshFiles()

	return m, nil
}

func (m *Model) cd(args []string) (tea.Model, tea.Cmd) {
	m.Display.LeftScroll = 0
	m.Display.RightScroll = 0
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
		err := bookmarks.SetBookmark(name, m.CurPath)
		if err != nil {
			m.ShowAlert(alert.Error, "Error saving the bookmark:", err.Error())
		}

	case "load":
		if len(args) < 2 {
			return m, nil
		}

		name := args[1]

		if path, ok := bookmarks.GetBookmark(name); !ok {
			m.ShowAlert(alert.Error, "No bookmark '"+name+"' found.")
		} else {
			m.NavigateInto(path)
			m.refreshFiles()
		}

	case "list":
		var text strings.Builder
		text.WriteString("Bookmark list:\n")
		for name, path := range bookmarks.GetBookmarks() {
			fmt.Fprintf(&text, "%s: %s\n", name, path)
		}
		m.ShowAlert(alert.Info, text.String())

	case "delete":
		if len(args) < 2 {
			return m, nil
		}
		name := args[1]
		err := bookmarks.DeleteBookmark(name)
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

			if m.Cursor < m.Display.LeftScroll {
				m.Display.LeftScroll = m.Cursor
			} else if m.Cursor >= m.Display.LeftScroll+visibleHeight {
				m.Display.LeftScroll = m.Cursor - visibleHeight + 1
			}

			break
		}
	}

	return m, m.loadPreview()
}

func (m *Model) selectFiles(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		return m, nil
	}

	pattern := strings.Join(args, " ")

	for _, f := range m.FileTree {
		matched, _ := filepath.Match(pattern, f.Name)
		if matched {
			m.IO.Selected[filepath.Join(m.CurPath, f.Name)] = struct{}{}
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
			delete(m.IO.Selected, filepath.Join(m.CurPath, f.Name))
		}
	}

	return m, nil
}

func (m *Model) yank() (tea.Model, tea.Cmd) {
	files := m.getTargets()
	m.IO.Yanked = make([]string, len(files))
	copy(m.IO.Yanked, files)

	return m, nil
}

func (m *Model) paste() (tea.Model, tea.Cmd) {
	for _, src := range m.IO.Yanked {
		name := filepath.Base(src)

		i := 1
		dstPath := filepath.Join(m.CurPath, name)
		for {
			if _, err := os.Stat(dstPath); err != nil {
				break
			}

			ext := filepath.Ext(name)
			base := strings.TrimSuffix(name, ext)

			if i == 1 {
				if base == "" {
					// dotfiles like .gitignore
					dstPath = filepath.Join(m.CurPath, ext+" copy")

				} else {
					dstPath = filepath.Join(m.CurPath, base+" copy"+ext)
				}
			} else {
				if base == "" {
					// dotfiles like .gitignore
					dstPath = filepath.Join(m.CurPath, fmt.Sprintf("%s copy %v", ext, i))
				} else {
					dstPath = filepath.Join(m.CurPath, fmt.Sprintf("%s copy %v%s", base, i, ext))
				}
			}

			i++
		}

		if err := filesystem.CopyFile(src, dstPath); err != nil {
			m.ShowAlert(alert.Error, "Paste error:", err.Error())
			return m, nil
		}
	}

	m.refreshFiles()

	return m, nil
}
