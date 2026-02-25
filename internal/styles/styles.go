package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/radeqq007/duvet/internal/config"
)

var (
	PaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(config.Colors.PaneBorder))

	FocusedPaneStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(config.Colors.FocusedPaneBorder))

	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(config.Colors.SelectedFileBG)).
			Foreground(lipgloss.Color(config.Colors.SelectedFileFG)).
			Bold(true)

	DirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(config.Colors.DirFG)).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(config.Colors.FileFG))

	CmdBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(40).
			BorderForeground(lipgloss.Color(config.Colors.CmdBoxBorder)).
			Foreground(lipgloss.Color(config.Colors.CmdBoxFG))

	AlertBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 4).
			Width(40).
			BorderForeground(lipgloss.Color(config.Colors.AlertBorder)).
			Foreground(lipgloss.Color(config.Colors.AlertFG))
)
