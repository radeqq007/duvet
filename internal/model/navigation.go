package model

import (
	"os/exec"
	"runtime"
	"strings"
	
	"charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/radeqq007/duvet/internal/filesystem"
	"github.com/radeqq007/duvet/internal/git"
	"github.com/radeqq007/duvet/internal/pane"
)

func (m *Model) openFile(path string) tea.Cmd {
	if isMediaFile(path) {
		return m.openWithSystem(path)
	}

	c := exec.Command(m.config.DefaultEditor, path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return FileClosed{Err: err}
	})
}

func (m *Model) openWithSystem(path string) tea.Cmd {
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
	if m.Display.Focus != pane.Left {
		return
	}

	if m.Cursor > 0 {
		m.Cursor--

		if m.Cursor < m.Display.LeftScroll {
			m.Display.LeftScroll = m.Cursor
		}
	}
}

func (m *Model) NavigateDown() {
	if m.Display.Focus != pane.Left {
		return
	}

	if m.Cursor < len(m.FileTree)-1 {
		m.Cursor++

		visibleHeight := m.VisibleHeight() - m.config.Layout.StatusBarHeight - m.config.Layout.BorderWidth
		if m.Cursor >= m.Display.LeftScroll+visibleHeight {
			m.Display.LeftScroll = m.Cursor - visibleHeight + 1
		}
	}
}

func (m *Model) NavigateToParent() error {
	parentDir := m.getParentDir()
	files, err := filesystem.GetFiles(parentDir)
	if err != nil {
		return err
	}

	m.CurPath = parentDir
	m.FileTree = files
	m.Cursor = 0
	m.Display.LeftScroll = 0
	m.Display.RightScroll = 0
	m.IO.Selected = make(map[string]struct{})
	m.Git = git.GetStatus(m.CurPath)

	return nil
}

func (m *Model) NavigateInto() error {
	newPath := m.CurrentFilePath()

	files, err := filesystem.GetFiles(newPath)
	if err != nil {
		return err
	}

	m.CurPath = newPath
	m.FileTree = files
	m.Cursor = 0
	m.Display.LeftScroll = 0
	m.Display.RightScroll = 0
	m.IO.Selected = make(map[string]struct{})
	m.Git = git.GetStatus(m.CurPath)

	return nil
}

func (m *Model) ScrollRightUp() {
	if m.Display.RightScroll > 0 {
		m.Display.RightScroll--
	}
}

func (m *Model) ScrollRightDown() {
	// TODO: MAKE THIS CODE NOT SUCK
	// calculate the amount of all rendered lines
	wrapped := lipgloss.NewStyle().
			Width(m.Layout.Width/2 - m.config.Layout.BorderWidth*2).
			Render(m.Display.Preview.Content)

	lines := strings.Split(wrapped, "\n")
	
	if m.Display.RightScroll < len(lines) - 1 {
		m.Display.RightScroll++
	}
}

func (m *Model) ResetRightScroll() {
	m.Display.RightScroll = 0
}
