package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/radeqq007/duvet/internal/alert"
	"github.com/radeqq007/duvet/internal/filesystem"
	"github.com/radeqq007/duvet/internal/icons"
	"github.com/radeqq007/duvet/internal/mode"
	"github.com/radeqq007/duvet/internal/pane"
	"github.com/radeqq007/duvet/internal/styles"
	"github.com/radeqq007/duvet/internal/ui"
)

func (m Model) View() string {
	if m.Width == 0 {
		return "loading..."
	}

	var leftContent strings.Builder
	visibleHeight := m.VisibleHeight()

	start := m.LeftScroll
	end := m.LeftScroll + visibleHeight
	end = min(end, len(m.FileTree))

	for i := start; i < end; i++ {
		node := m.FileTree[i]

		_, isSelected := m.Selected[filepath.Join(m.CurPath, node.Name)]

		var icon string
		if node.IsDir {
			icon = "\uf4d3"
		} else {
			icon = icons.GetIcon(filepath.Ext(node.Name))
		}

		if isSelected {
			icon = "▌ " + icon
		}

		line := fmt.Sprintf("%s %s", icon, node.Name)

		if i == m.Cursor {
			line = styles.SelectedStyle.Width(m.Width / 2).Render(line)
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
			Width(m.Width / 2).
			Height(m.Height - 2).
			Render(leftContent.String())
	} else {
		leftPane = styles.PaneStyle.
			Width(m.Width / 2).
			Height(m.Height - 2).
			Render(leftContent.String())
	}

	var rightContent strings.Builder
	if !m.FileTree[m.Cursor].IsDir {
		file := filepath.Join(m.CurPath, m.FileTree[m.Cursor].Name)
		content, _ := filesystem.ReadFileContent(file)

		wrapped := lipgloss.NewStyle().
			Width(m.Width/2 - 2).
			Render(content)

		visualLines := strings.Split(wrapped, "\n")

		start := m.RightScroll
		end := min(start+visibleHeight, len(visualLines))

		for i := start; i < end; i++ {
			rightContent.WriteString(visualLines[i] + "\n")
		}

		linesRendered := end - start
		for i := linesRendered; i < visibleHeight; i++ {
			rightContent.WriteByte('\n')
		}
	} else {
		for range visibleHeight {
			_ = rightContent.WriteByte('\n')
		}
	}

	var rightPane string
	if m.Focus == pane.Right {
		rightPane = styles.FocusedPaneStyle.
			Width(m.Width / 2).
			Height(m.Height - 2).
			Render(rightContent.String())
	} else {
		rightPane = styles.PaneStyle.
			Width(m.Width / 2).
			Height(m.Height - 2).
			Render(rightContent.String())
	}

	view := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	switch m.Mode {
	case mode.Command:
		content := ":" + m.CmdInput
		if strings.HasPrefix(m.CmdInput, "!") {
			content = "$" + m.CmdInput[1:]
		}
		cmdBox := styles.CmdBoxStyle.Render(content)

		x := m.Width/2 - lipgloss.Width(cmdBox)/2
		y := m.Height/2 - lipgloss.Height(cmdBox)/2

		view = ui.PlaceOverlay(x, y, cmdBox, view)

	case mode.Alert:
		var alertBox string
		switch m.Alert.Type {
		case alert.Normal:
			alertBox = styles.AlertNormalStyle.Render(m.Alert.Text)
		case alert.Info:
			alertBox = styles.AlertInfoStyle.Render(m.Alert.Text)
		case alert.Error:
			alertBox = styles.AlertErrorStyle.Render(m.Alert.Text)
		case alert.Warning:
			alertBox = styles.AlertWarningStyle.Render(m.Alert.Text)
		}

		x := m.Width/2 - lipgloss.Width(alertBox)/2
		y := m.Height/2 - lipgloss.Height(alertBox)/2

		view = ui.PlaceOverlay(x, y, alertBox, view)
	}

	return view
}

func (m *Model) UpdateDimensions(width, height int) {
	m.Width = width
	m.Height = height
}
