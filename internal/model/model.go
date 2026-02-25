package model

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/radeqq007/duvet/internal/config"
	"github.com/radeqq007/duvet/internal/filesystem"
	"github.com/radeqq007/duvet/internal/icons"
	"github.com/radeqq007/duvet/internal/mode"
	"github.com/radeqq007/duvet/internal/pane"
	"github.com/radeqq007/duvet/internal/styles"
	"github.com/radeqq007/duvet/internal/ui"
)

type Model struct {
	FileTree    []filesystem.FileNode
	Cursor      int
	Focus       pane.Pane
	LeftScroll  int
	RightScroll int
	Width       int
	Height      int
	LeftPaneW   int
	RightPaneW  int
	CurPath     string
	ParentDir   string
	Mode        mode.Mode
	CmdInput    string
}

type FileClosed struct{ Err error }

func (m Model) Init() tea.Cmd {
	return nil
}

func New() Model {
	dir, _ := os.Getwd()

	files, err := filesystem.GetFiles(dir)
	if err != nil {
		files = []filesystem.FileNode{}
	}

	return Model{
		FileTree:    files,
		Cursor:      0,
		LeftScroll:  0,
		RightScroll: 0,
		Focus:       0,
		LeftPaneW:   40,
		RightPaneW:  40,
		CurPath:     dir,
		ParentDir:   filepath.Dir(dir),
	}
}

func (m Model) VisibleHeight() int {
	return m.Height - config.HeaderFooterSize
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

func (m Model) View() string {
	if m.Width == 0 {
		return "loading..."
	}

	var leftContent strings.Builder
	visibleHeight := m.Height - 4

	start := m.LeftScroll
	end := m.LeftScroll + visibleHeight
	end = min(end, len(m.FileTree))

	for i := start; i < end; i++ {
		node := m.FileTree[i]

		var icon string
		if node.IsDir {
			icon = "\uf4d3"
		} else {
			icon = icons.GetIcon(filepath.Ext(node.Name))
		}

		line := fmt.Sprintf("%s %s", icon, node.Name)

		if i == m.Cursor {
			line = styles.SelectedStyle.Width(m.LeftPaneW - 2).Render(line)
		} else if node.IsDir {
			line = styles.DirStyle.Render(line)
		} else {
			line = styles.FileStyle.Render(line)
		}

		_, _ = leftContent.WriteString(line + "\n")
	}

	for i := len(m.FileTree); i < visibleHeight; i++ {
		leftContent.WriteByte('\n')
	}

	var leftPane string
	if m.Focus == pane.Left {
		leftPane = styles.FocusedPaneStyle.
			Width(m.LeftPaneW).
			Height(m.Height - 2).
			Render(leftContent.String())
	} else {
		leftPane = styles.PaneStyle.
			Width(m.LeftPaneW).
			Height(m.Height - 2).
			Render(leftContent.String())
	}

	var rightContent strings.Builder
	if !m.FileTree[m.Cursor].IsDir {
		file := filepath.Join(m.CurPath, m.FileTree[m.Cursor].Name)
		content, _ := filesystem.ReadFileContent(file)

		lines := strings.Split(content, "\n")
		i := m.RightScroll
		for i < min(visibleHeight+m.RightScroll, len(lines)) {
			_, _ = rightContent.WriteString(lines[i] + "\n")
			if len(lines[i]) > m.RightPaneW {
				i++
			}

			i++
		}
	} else {
		for range visibleHeight {
			_ = rightContent.WriteByte('\n')
		}
	}

	var rightPane string
	if m.Focus == pane.Right {
		rightPane = styles.FocusedPaneStyle.
			Width(m.RightPaneW).
			Height(m.Height - 2).
			Render(rightContent.String())
	} else {
		rightPane = styles.PaneStyle.
			Width(m.RightPaneW).
			Height(m.Height - 2).
			Render(rightContent.String())
	}

	view := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	if m.Mode == mode.Command {
		content := "$ " + m.CmdInput
		cmdBox := styles.CmdBoxStyle.Render(content)

		x := m.Width/2 - lipgloss.Width(cmdBox)/2
		y := m.Height/2 - lipgloss.Height(cmdBox)/2

		view = ui.PlaceOverlay(x, y, cmdBox, view)
	}

	return view
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			if m.Mode == mode.Normal && msg.String() == ":" {
				m.Mode = mode.Command
				m.CmdInput = ""
				return m, nil
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Focus == pane.Left && m.Cursor > 0 {
				m.Cursor--

				if m.Cursor < m.LeftScroll {
					m.LeftScroll = m.Cursor
				}
			} else {
				if m.RightScroll == 0 {
					break
				}

				m.RightScroll--
			}

		case "down", "j":
			if m.Focus == pane.Left && m.Cursor < len(m.FileTree)-1 {
				m.Cursor++

				visibleHeight := m.Height - 4
				if m.Cursor >= m.LeftScroll+visibleHeight {
					m.LeftScroll = m.Cursor - visibleHeight + 1
				}
			} else {
				m.RightScroll++
			}

		case "left", "h":
			files, err := filesystem.GetFiles(m.ParentDir)
			if err == nil {
				m.CurPath = m.ParentDir
				m.ParentDir = filepath.Dir(m.CurPath)
				m.FileTree = files
				m.Cursor = 0
			}

		case "ctrl+right":
			if m.Focus == pane.Left {
				m.Focus = 1
			}

		case "ctrl+left":
			if m.Focus == pane.Right {
				m.Focus = 0
			}

		case "enter", "l", "right":
			path := m.FileTree[m.Cursor]
			newPath := filepath.Join(m.CurPath, path.Name)
			if path.IsDir {
				files, err := filesystem.GetFiles(newPath)
				if err == nil {
					m.ParentDir = m.CurPath
					m.CurPath = newPath
					m.FileTree = files
					m.Cursor = 0
					m.LeftScroll = 0
				}
			} else {
				return m, openFile(newPath)
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.LeftPaneW = msg.Width/2 - 2
		m.RightPaneW = msg.Width/2 - 2
	}

	return m, nil
}

func openFile(filepath string) tea.Cmd {
	c := exec.Command("nvim", filepath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return FileClosed{Err: err}
	})
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

func (m *Model) UpdateDimensions(width, height int) {
	m.Width = width
	m.Height = height
	m.LeftPaneW = width/2 - config.BorderWidth
	m.RightPaneW = width/2 - config.BorderWidth
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
