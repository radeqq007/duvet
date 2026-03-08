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
		m.Layout.Width = msg.Width
		m.Layout.Height = msg.Height
	}

	switch m.IO.Mode {
	case mode.Normal:
		return m.handleNormalModeUpdate(msg)
	case mode.Command:
		return m.handleCommandModeUpdate(msg)
	case mode.Alert:
		if _, ok := msg.(tea.KeyMsg); ok {
			m.IO.Mode = mode.Normal
		}
		return m, nil
	}

	m.IO.Mode = mode.Normal

	return m, nil
}

func (m Model) handleNormalModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			if m.IO.Mode == mode.Normal && msg.String() == ":" {
				m.IO.Mode = mode.Command
				m.IO.CmdInput = ""
				return m, nil
			}

		case "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.Display.Focus == pane.Left {
				m.NavigateUp()
				return m, m.loadPreview()
			} else {
				m.ScrollRightUp()
			}

		case "down", "j":
			if m.Display.Focus == pane.Left {
				m.NavigateDown()
				return m, m.loadPreview()
			} else {
				m.ScrollRightDown()
			}

		case "left", "h":
			if m.Display.Focus == pane.Left {
				if err := m.NavigateToParent(); err != nil {
					m.ShowAlert(alert.Error, "Cannot navigate to parent:", err.Error())
				}

				m.Display.Preview = Preview{}
			} else {
				m.Display.Focus = pane.Left
			}

		case "right", "l":
			m.Display.Focus = pane.Right

		case "enter":
			path := m.FileTree[m.Cursor]
			if path.IsDir {
				if err := m.NavigateInto(); err != nil {
					m.ShowAlert(alert.Error, "Cannot navigate into:", err.Error())
				}
				m.Display.Preview = Preview{}
				return m, m.loadPreview()
			} else {
				newPath := filepath.Join(m.CurPath, path.Name)
				return m, openFile(newPath)
			}

		case " ":
			path := m.CurrentFilePath()
			if _, ok := m.IO.Selected[path]; ok {
				delete(m.IO.Selected, path)
			} else {
				m.IO.Selected[path] = struct{}{}
			}
			// TODO: navigating down kinda gives nice UX but also can be annoying
			// m.NavigateDown()
		}

	case command.Msg:
		return m.handleCommand(msg)

	case PreviewLoaded:
		m.Display.Preview = Preview{
			Path:    msg.Path,
			Content: msg.Content,
		}
	}

	return m, nil
}

func (m Model) handleCommandModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(tea.KeyMsg).String(); msg {
	case "esc":
		m.IO.Mode = mode.Normal
		m.IO.CmdInput = ""

	case "enter":
		m.IO.Mode = mode.Normal
		return m, command.Exec(m.IO.CmdInput)

	case "backspace":
		if len(m.IO.CmdInput) >= 1 {
			m.IO.CmdInput = m.IO.CmdInput[:len(m.IO.CmdInput)-1]
		}

	default:
		m.IO.CmdInput += msg
	}

	return m, nil
}
