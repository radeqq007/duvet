package model

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/radeqq007/duvet/internal/alert"
	"github.com/radeqq007/duvet/internal/command"
	"github.com/radeqq007/duvet/internal/mode"
	"github.com/radeqq007/duvet/internal/pane"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.Width = msg.Width
		m.Height = msg.Height
	}

	switch m.Mode {
	case mode.Normal:
		return m.handleNormalModeUpdate(msg)
	case mode.Command:
		return m.handleCommandModeUpdate(msg)
	case mode.Alert:
		if _, ok := msg.(tea.KeyMsg); ok {
			m.Mode = mode.Normal
		}
		return m, nil
	}

	m.Mode = mode.Normal

	return m, nil
}

func (m Model) handleNormalModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			if m.Mode == mode.Normal && msg.String() == ":" {
				m.Mode = mode.Command
				m.CmdInput = ""
				return m, nil
			}

		case "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.Focus == pane.Left {
				m.NavigateUp()
			} else {
				m.ScrollRightUp()
			}

		case "down", "j":
			if m.Focus == pane.Left {
				m.NavigateDown()
			} else {
				m.ScrollRightDown()
			}

		case "left", "h":
			if m.Focus == pane.Left {
				if err := m.NavigateToParent(); err != nil {
					m.ShowAlert(alert.Error, "Cannot navigate to parent:", err.Error())
				}

			} else {
				m.Focus = pane.Left
			}

		case "right", "l":
			m.Focus = pane.Right

		case "enter":
			path := m.FileTree[m.Cursor]
			if path.IsDir {
				if err := m.NavigateInto(); err != nil {
					m.ShowAlert(alert.Error, "Cannot navigate into:", err.Error())
				}
			} else {
				newPath := filepath.Join(m.CurPath, path.Name)
				return m, openFile(newPath)
			}

		case " ":
			path := m.CurrentFilePath()
			if _, ok := m.Selected[path]; ok {
				delete(m.Selected, path)
			} else {
				m.Selected[path] = struct{}{}
			}
			// TODO: navigating down kinda gives nice UX but also can be annoying
			// m.NavigateDown()
		}

	case command.Msg:
		return m.handleCommand(msg)
	}

	return m, nil
}

func (m Model) handleCommandModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(tea.KeyMsg).String(); msg {
	case "esc":
		m.Mode = mode.Normal
		m.CmdInput = ""

	case "enter":
		m.Mode = mode.Normal
		return m, command.Exec(m.CmdInput)

	case "backspace":
		if len(m.CmdInput) >= 1 {
			m.CmdInput = m.CmdInput[:len(m.CmdInput)-1]
		}

	default:
		m.CmdInput += msg
	}

	return m, nil
}
