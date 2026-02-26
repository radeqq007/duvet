package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/radeqq007/duvet/internal/alert"
	"github.com/radeqq007/duvet/internal/command"
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
	var file string
	var path string
	if len(args) < 1 {
		file = m.getCurrentFile().Name
		path = filepath.Join(m.CurPath, file)
	} else {
		file = args[0]
		path = filepath.Join(m.CurPath, file)
	}

	err := os.RemoveAll(path)
	if err != nil {
		m.ShowAlert(alert.Error, "Error removing a file: ", err.Error())
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

	m.refreshFiles()

	return m, nil
}

func (m *Model) bookmark(args []string) (tea.Model, tea.Cmd) {
	switch args[0] {
	case "save":
		if len(args) < 2 {
			return m, nil
		}

		name := args[1]
		m.Bookmarks[name] = m.CurPath

	case "load":
		if len(args) < 2 {
			return m, nil
		}

		name := args[1]
		if path, ok := m.Bookmarks[name]; ok {
			m.CurPath = path
			m.refreshFiles()
		} else {
			m.ShowAlert(alert.Error, "No bookmark", name, "found.")
		}

	case "list":
		var text strings.Builder
		text.WriteString("Bookmark list:\n")
		for name, path := range m.Bookmarks {
			fmt.Fprintf(&text, "%s: %s\n", name, path)
		}
		m.ShowAlert(alert.Info, text.String())

	case "remove":
		if len(args) < 2 {
			return m, nil
		}

		name := args[1]
		delete(m.Bookmarks, name)
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
