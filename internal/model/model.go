package model

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/radeqq007/duvet/internal/alert"
	"github.com/radeqq007/duvet/internal/config"
	"github.com/radeqq007/duvet/internal/filesystem"
	"github.com/radeqq007/duvet/internal/git"
	"github.com/radeqq007/duvet/internal/mode"
	"github.com/radeqq007/duvet/internal/pane"
)

type Model struct {
	config   *config.Config
	FileTree []filesystem.FileNode
	Cursor   int
	CurPath  string
	Git      *git.Status
	Layout   LayoutState
	Display  ViewState
	IO       IOState
}

type LayoutState struct {
	Width  int
	Height int
}

type ViewState struct {
	LeftScroll  int
	RightScroll int
	Preview     Preview
	Focus       pane.Pane
}

type IOState struct {
	Mode     mode.Mode
	CmdInput string
	Input    []byte
	Alert    alert.Alert
	Selected map[string]struct{}
	Yanked   []string
}

type Preview struct {
	Path    string
	Content string
}

type FileClosed struct{ Err error }

type PreviewLoaded struct {
	Path    string
	Content string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New(cfg *config.Config) Model {
	dir, _ := os.Getwd()

	files, err := filesystem.GetFiles(dir)
	if err != nil {
		files = []filesystem.FileNode{}
	}

	_ = config.LoadBookmarks()

	return Model{
		config:   cfg,
		FileTree: files,
		CurPath:  dir,
		Git:      git.GetStatus(dir),
		IO: IOState{
			Selected: make(map[string]struct{}),
		},
	}
}
