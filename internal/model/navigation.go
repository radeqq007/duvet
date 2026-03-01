package model

import (
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/radeqq007/duvet/internal/config"
	"github.com/radeqq007/duvet/internal/filesystem"
	"github.com/radeqq007/duvet/internal/pane"
)

func openFile(path string) tea.Cmd {
	if isMediaFile(path) {
		return openWithSystem(path)
	}

	c := exec.Command(config.DefaultEditor, path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return FileClosed{Err: err}
	})
}

func openWithSystem(path string) tea.Cmd {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", "", path)
	default: // linux, bsd, etc.
		cmd = exec.Command("xdg-open", path)
	}

	return func() tea.Msg {
		err := cmd.Start()
		return FileClosed{Err: err}
	}
}

func (m *Model) NavigateUp() {
	if m.Focus != pane.Left {
		return
	}

	if m.Cursor > 0 {
		m.Cursor--

		if m.Cursor < m.LeftScroll {
			m.LeftScroll = m.Cursor
		}
	}
}

func (m *Model) NavigateDown() {
	if m.Focus != pane.Left {
		return
	}

	if m.Cursor < len(m.FileTree)-1 {
		m.Cursor++
		visibleHeight := m.VisibleHeight()
		if m.Cursor >= m.LeftScroll+visibleHeight {
			m.LeftScroll = m.Cursor - visibleHeight + 1
		}
	}
}

func (m *Model) NavigateToParent() error {
	files, err := filesystem.GetFiles(m.ParentDir)
	if err != nil {
		return err
	}

	m.CurPath = m.ParentDir
	m.ParentDir = filesystem.ParentDir(m.CurPath)
	m.FileTree = files
	m.Cursor = 0
	m.LeftScroll = 0
	m.RightScroll = 0
	m.Selected = make(map[string]struct{})

	return nil
}

func (m *Model) NavigateInto() error {
	newPath := m.CurrentFilePath()

	files, err := filesystem.GetFiles(newPath)
	if err != nil {
		return err
	}

	m.ParentDir = m.CurPath
	m.CurPath = newPath
	m.FileTree = files
	m.Cursor = 0
	m.LeftScroll = 0
	m.RightScroll = 0
	m.Selected = make(map[string]struct{})

	return nil
}

func (m *Model) ScrollRightUp() {
	if m.RightScroll > 0 {
		m.RightScroll--
	}
}

func (m *Model) ScrollRightDown() {
	m.RightScroll++
}

func (m *Model) ResetRightScroll() {
	m.RightScroll = 0
}
