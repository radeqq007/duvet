package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/radeqq007/duvet/internal/icons"
)

type fileNode struct {
	name  string
	isDir bool
}

type model struct {
	fileTree   []fileNode
	cursor     int
	scroll     int
	width      int
	height     int
	leftPaneW  int
	rightPaneW int
	curPath    string
	parentDir  string
}

type fileClosed struct{ err error }

var (
	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Bold(true)

	dirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
)

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	dir, _ := os.Getwd()

	files, err := getFiles(dir)
	if err != nil {
		files = []fileNode{}
	}

	return model{
		fileTree:   files,
		cursor:     0,
		leftPaneW:  40,
		rightPaneW: 40,
		curPath:    dir,
		parentDir:  filepath.Dir(dir),
	}
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	leftContent := ""
	visibleHeight := m.height - 4

	start := m.scroll
	end := m.scroll + visibleHeight
	end = min(end, len(m.fileTree))

	for i := start; i < end; i++ {
		node := m.fileTree[i]

		var icon string
		if node.isDir {
			icon = "\uf4d3"
		} else {
			icon = icons.GetIcon(filepath.Ext(node.name))
		}

		line := fmt.Sprintf("%s %s", icon, node.name)

		if i == m.cursor {
			line = selectedStyle.Width(m.leftPaneW - 2).Render(line)
		} else if node.isDir {
			line = dirStyle.Render(line)
		} else {
			line = fileStyle.Render(line)
		}

		leftContent += line + "\n"
	}

	for i := len(m.fileTree); i < visibleHeight; i++ {
		leftContent += "\n"
	}

	leftPane := paneStyle.
		Width(m.leftPaneW).
		Height(m.height - 2).
		Render(leftContent)

	var rightContent strings.Builder
	if !m.fileTree[m.cursor].isDir {
		file := filepath.Join(m.curPath, m.fileTree[m.cursor].name)
		content, _ := readFileContent(file)

		lines := strings.Split(content, "\n")
		i := 0
		for i < min(visibleHeight, len(lines)) {
			_, _ = rightContent.WriteString(lines[i] + "\n")
			if len(lines[i]) > m.rightPaneW {
				i++
			}
			
			i++
		}
	} else {
		for range visibleHeight {
			_ = rightContent.WriteByte('\n')
		}
	}

	rightPane := paneStyle.
		Width(m.rightPaneW).
		Height(m.height - 2).
		Render(rightContent.String())

	view := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	return view
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--

				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}

		case "down", "j":
			if m.cursor < len(m.fileTree)-1 {
				m.cursor++

				visibleHeight := m.height - 4
				if m.cursor >= m.scroll+visibleHeight {
					m.scroll = m.cursor - visibleHeight + 1
				}
			}

		case "left", "h":
			files, err := getFiles(m.parentDir)
			if err == nil {
				m.curPath = m.parentDir
				m.parentDir = filepath.Dir(m.curPath)
				m.fileTree = files
				m.cursor = 0
			}

		case "enter", "l", "right":
			path := m.fileTree[m.cursor]
			newPath := filepath.Join(m.curPath, path.name) 
			if path.isDir {
				files, err := getFiles(newPath)
				if err == nil {
					m.parentDir = m.curPath
					m.curPath = newPath
					m.fileTree = files
					m.cursor = 0
					m.scroll = 0
				}
			} else {
				return m, openFile(newPath)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.leftPaneW = msg.Width/2 - 2
		m.rightPaneW = msg.Width/2 - 2
	}

	return m, nil
}

func getFiles(path string) ([]fileNode, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	entries, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	allFiles := make([]fileNode, len(entries))

	for i, entry := range entries {
		allFiles[i] = fileNode{
			name:  entry.Name(),
			isDir: entry.IsDir(),
		}
	}

	return allFiles, nil
}

func readFileContent(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	// Check if it's a binary file by looking for null bytes
	// Prolly not the best way to do that
	checkLen := 512
	if len(content) < checkLen {
		checkLen = len(content)
	}

	for i := 0; i < checkLen; i++ {
		if content[i] == 0 {
			return "", fmt.Errorf("binary file")
		}
	}

	text := string(content)
	return text, err
}

func openFile(filepath string) tea.Cmd {
	c := exec.Command("nvim", filepath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return fileClosed{ err }
	})
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
