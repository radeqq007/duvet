package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/radeqq007/duvet/internal/config"
)

var (
	PaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(config.PaneBorderColor))

	FocusedPaneStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(config.FocusedPaneBorderColor))

	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(config.SelectedFileBackground)).
			Foreground(lipgloss.Color(config.SelectedFileForeground)).
			Bold(true)

	DirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(config.DirForeground)).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(config.FileForeground))

	CmdBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(40).
			BorderForeground(lipgloss.Color(config.CmdBoxBorderColor)).
			Foreground(lipgloss.Color(config.CmdBoxForeground))

	AlertBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 4).
			Width(40).
			BorderForeground(lipgloss.Color(config.AlertBoxBorderColor)).
			Foreground(lipgloss.Color(config.AlertBoxForeground))
)
