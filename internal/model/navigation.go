package model

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/radeqq007/duvet/internal/filesystem"
)

func openFile(filepath string) tea.Cmd {
	c := exec.Command("nvim", filepath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return FileClosed{Err: err}
	})
}

func (m Model) NavigateUp() {
	if m.Cursor > 0 {
		m.Cursor--
		if m.Cursor < m.LeftScroll {
			m.LeftScroll = m.Cursor
		}
	}
}

func (m *Model) NavigateDown() {
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

	return nil
}

func (m *Model) NavigateInto() error {
	current := m.CurrentFile()
	newPath := m.CurrentFilePath()

	if current.IsDir {
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
	}

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
