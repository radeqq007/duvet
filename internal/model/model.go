package model

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/radeqq007/duvet/internal/alert"
	"github.com/radeqq007/duvet/internal/config"
	"github.com/radeqq007/duvet/internal/filesystem"
	"github.com/radeqq007/duvet/internal/mode"
	"github.com/radeqq007/duvet/internal/pane"
)

type Model struct {
	FileTree    []filesystem.FileNode
	Cursor      int
	Focus       pane.Pane
	LeftScroll  int
	RightScroll int
	Width       int
	Height      int
	CurPath     string
	ParentDir   string
	Mode        mode.Mode
	CmdInput    string
	Alert       alert.Alert
	Preview     Preview
	Selected    map[string]struct{}
}

type Preview struct {
	Path    string
	Content string
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

	config.Load()
	config.LoadBookmarks()

	return Model{
		FileTree:    files,
		Cursor:      0,
		LeftScroll:  0,
		RightScroll: 0,
		Focus:       0,
		CurPath:     dir,
		ParentDir:   filepath.Dir(dir),
		Alert: alert.Alert{
			Type: alert.Normal,
			Text: "",
		},
		Preview: Preview{
			Path:    "",
			Content: "",
		},
		Selected: make(map[string]struct{}),
	}
}
