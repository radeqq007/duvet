package model

import (
	"path/filepath"
	"strings"

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
			if m.IO.Mode == mode.Normal {
				m.IO.Mode = mode.Command
				m.IO.CmdInput = ""
				return m, nil
			}

		case "ctrl+c":
			return m, tea.Quit

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
				return m, m.openFile(newPath)
			}

		case "tab":
			if m.Display.Focus == pane.Left {
				m.Display.Focus = pane.Right
			} else {
				m.Display.Focus = pane.Left
			}

		case "esc":
			m.clearInput()

		case " ":
			path := m.CurrentFilePath()
			if _, ok := m.IO.Selected[path]; ok {
				delete(m.IO.Selected, path)
			} else {
				m.IO.Selected[path] = struct{}{}
			}

			m.NavigateDown()

		case "up":
			m.IO.Input = append(m.IO.Input, 'k')
			return m.handleInput()

		case "down":
			m.IO.Input = append(m.IO.Input, 'j')
			return m.handleInput()

		case "left":
			m.IO.Input = append(m.IO.Input, 'h')
			return m.handleInput()

		case "right":
			m.IO.Input = append(m.IO.Input, 'l')
			return m.handleInput()

		default:
			if len(msg.Runes) > 0 {
				m.IO.Input = append(m.IO.Input, byte(msg.Runes[0]))
			}
			return m.handleInput()
		}

	case command.Msg:
		return m.handleCommand(msg)

	case PreviewLoaded:
		m.Display.Preview = Preview(msg)
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

var pendingPrefixes = map[string]struct{}{
	"d": {}, "y": {}, "g": {},
}

func (m *Model) handleInput() (tea.Model, tea.Cmd) {
	count, motion := m.parseInput()
	// defer m.clearInput()

	switch motion {
	case "k":
		if m.Display.Focus == pane.Left {
			for range count {
				m.NavigateUp()
			}
			m.clearInput()
			return m, m.loadPreview()
		} else {
			for range count {
				m.ScrollRightUp()
			}
			m.clearInput()
		}

	case "j":
		if m.Display.Focus == pane.Left {
			for range count {
				m.NavigateDown()
			}
			m.clearInput()
			return m, m.loadPreview()
		} else {
			for range count {
				m.ScrollRightDown()
			}
			m.clearInput()
		}

	case "h":
		if err := m.NavigateToParent(); err != nil {
			m.ShowAlert(alert.Error, "Cannot navigate to parent:", err.Error())
		}

		m.Display.Preview = Preview{}

		m.clearInput()

	case "l":
		path := m.FileTree[m.Cursor]
		if path.IsDir {
			if err := m.NavigateInto(); err != nil {
				m.ShowAlert(alert.Error, "Cannot navigate into:", err.Error())
			}
			m.Display.Preview = Preview{}
			m.clearInput()
			return m, m.loadPreview()
		} else {
			newPath := filepath.Join(m.CurPath, path.Name)
			m.clearInput()
			return m, m.openFile(newPath)
		}

	case "yy":
		m.clearInput()
		m.yank()

	case "p":
		m.clearInput()
		m.paste()

	case "dd":
		m.clearInput()
		m.delete()

	case "gg":
		if m.Display.Focus == pane.Left {
			m.Display.LeftScroll = 0
			m.Cursor = 0
			m.clearInput()
			return m, m.loadPreview()
		} else {
			m.Display.RightScroll = 0
			m.clearInput()
		}

	case "G":
		if m.Display.Focus == pane.Left {
			visibleHeight := m.VisibleHeight() - m.config.Layout.StatusBarHeight - m.config.Layout.BorderWidth

			if len(m.IO.Input) > 1 {
				// has a line number
				m.Cursor = count - 1
				m.Display.LeftScroll = max(0, m.Cursor-visibleHeight+1)
				return m, m.loadPreview()
			}

			m.Cursor = len(m.FileTree) - 1
			m.Display.LeftScroll = max(0, m.Cursor-visibleHeight+1)
			m.clearInput()

			return m, m.loadPreview()

		} else {
			// TODO: add <line number>G motion support for the right pane

			lines := strings.Split(wrapLines(m.Display.Preview.Content, m.Layout.Width/2-m.config.Layout.BorderWidth*2), "\n")
			visibleHeight := m.VisibleHeight() - m.config.Layout.StatusBarHeight - m.config.Layout.BorderWidth
			m.Display.RightScroll = max(0, len(lines)-visibleHeight)
			m.clearInput()
		}

	case "":
		// sequence still being built, do nothin

	default:
		if _, ok := pendingPrefixes[motion]; !ok {
			// unknown sequence, discard
			m.clearInput()
		}
	}

	return m, nil
}
