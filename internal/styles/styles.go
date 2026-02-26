package styles

import (
	"charm.land/lipgloss/v2"
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
			Width(40)

	AlertNormalStyle = AlertBoxStyle.
				BorderForeground(lipgloss.Color(config.Colors.AlertNormalBorder)).
				Foreground(lipgloss.Color(config.Colors.AlertNormalFG))

	AlertInfoStyle = AlertBoxStyle.
			BorderForeground(lipgloss.Color(config.Colors.AlertInfoBorder)).
			Foreground(lipgloss.Color(config.Colors.AlertInfoFG))

	AlertErrorStyle = AlertBoxStyle.
			BorderForeground(lipgloss.Color(config.Colors.AlertErrorBorder)).
			Foreground(lipgloss.Color(config.Colors.AlertErrorFG))

	AlertWarningStyle = AlertBoxStyle.
				BorderForeground(lipgloss.Color(config.Colors.AlertWarningBorder)).
				Foreground(lipgloss.Color(config.Colors.AlertWarningFG))
)
