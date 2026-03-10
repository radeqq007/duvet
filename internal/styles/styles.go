package styles

import (
	"charm.land/lipgloss/v2"
	"github.com/radeqq007/duvet/internal/config"
)

func PaneStyle(c config.ColorsConfig) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(c.PaneBorder))
}

func FocusedPaneStyle(c config.ColorsConfig) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(c.FocusedPaneBorder))
}

func SelectedStyle(c config.ColorsConfig) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(c.SelectedFileBG)).
		Foreground(lipgloss.Color(c.SelectedFileFG)).
		Bold(true)
}

func DirStyle(c config.ColorsConfig) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.DirFG)).
		Bold(true)
}

func FileStyle(c config.ColorsConfig) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.FileFG))
}

func CmdBoxStyle(c config.ColorsConfig) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Width(50).
		BorderForeground(lipgloss.Color(c.CmdBoxBorder)).
		Foreground(lipgloss.Color(c.CmdBoxFG))
}

func alertBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 1).
		Width(60)
}

func AlertNormalStyle(c config.ColorsConfig) lipgloss.Style {
	return alertBoxStyle().
		BorderForeground(lipgloss.Color(c.AlertNormalBorder)).
		Foreground(lipgloss.Color(c.AlertNormalFG))
}

func AlertInfoStyle(c config.ColorsConfig) lipgloss.Style {
	return alertBoxStyle().
		BorderForeground(lipgloss.Color(c.AlertInfoBorder)).
		Foreground(lipgloss.Color(c.AlertInfoFG))
}

func AlertErrorStyle(c config.ColorsConfig) lipgloss.Style {
	return alertBoxStyle().
		BorderForeground(lipgloss.Color(c.AlertErrorBorder)).
		Foreground(lipgloss.Color(c.AlertErrorFG))
}

func AlertWarningStyle(c config.ColorsConfig) lipgloss.Style {
	return alertBoxStyle().
		BorderForeground(lipgloss.Color(c.AlertWarningBorder)).
		Foreground(lipgloss.Color(c.AlertWarningFG))
}

func StatusBarStyle(c config.ColorsConfig) lipgloss.Style {
	return PaneStyle(c).
		Height(1).
		Padding(0, 1)
}
