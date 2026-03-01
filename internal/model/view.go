package model

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/radeqq007/duvet/internal/alert"
	"github.com/radeqq007/duvet/internal/config"
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

	leftPane := m.RenderLeftPane()
	rightPane := m.RenderRightPane()
	bar := m.RenderStatusBar()

	view := lipgloss.JoinHorizontal(lipgloss.Bottom, leftPane, rightPane)
	view = lipgloss.JoinVertical(lipgloss.Left, view, bar)

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

func (m *Model) RenderLeftPane() string {
	var leftContent strings.Builder

	visibleHeight := m.VisibleHeight() - config.Layout.StatusBarHeight - config.Layout.BorderWidth

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
			Height(visibleHeight).
			Render(leftContent.String())
	} else {
		leftPane = styles.PaneStyle.
			Width(m.Width / 2).
			Height(visibleHeight).
			Render(leftContent.String())
	}

	return leftPane
}

func (m *Model) RenderRightPane() string {
	var rightContent strings.Builder
	visibleHeight := m.VisibleHeight() - config.Layout.StatusBarHeight - config.Layout.BorderWidth

	if m.Preview.Content != "" {
		wrapped := lipgloss.NewStyle().
			Width(m.Width/2 - config.Layout.BorderWidth*2).
			Render(m.Preview.Content)

		visualLines := strings.Split(wrapped, "\n")
		start := m.RightScroll
		end := min(start+visibleHeight, len(visualLines))

		for i := start; i < end; i++ {
			rightContent.WriteString(visualLines[i] + "\n")
		}
		for i := end - start; i < visibleHeight; i++ {
			rightContent.WriteByte('\n')
		}
	} else {
		for range visibleHeight {
			rightContent.WriteByte('\n')
		}
	}

	var rightPane string
	if m.Focus == pane.Right {
		rightPane = styles.FocusedPaneStyle.
			Width(m.Width / 2).
			Height(visibleHeight).
			Render(rightContent.String())
	} else {
		rightPane = styles.PaneStyle.
			Width(m.Width / 2).
			Height(visibleHeight).
			Render(rightContent.String())
	}

	return rightPane
}

func (m *Model) RenderStatusBar() string {
	file := m.CurrentFile()
	var icon string
	if file.IsDir {
		icon = "\uf4d3"
	} else {
		icon = icons.GetIcon(file.Name)
	}

	seperator := " | "

	var left string
	left += prettifyPath(m.CurPath)

	left += seperator + strconv.Itoa(len(m.FileTree)) + " items"

	selected := len(m.Selected)
	if selected > 0 {
		left += seperator + "Selected: " + strconv.Itoa(selected)
	}

	var right string

	fileSize := ""
	if !file.IsDir {
		fileSize = filesystem.GetFileSize(m.CurrentFilePath())
	}

	right += icon + " " + file.Name

	if fileSize != "" {
		right += seperator + fileSize
	}

	// TODO: that -2 is padding, save that in the config or something cause now it's just a magic number
	spaceCount := m.Width - ansi.StringWidth(
		left,
	) - ansi.StringWidth(
		right,
	) - config.Layout.BorderWidth*2 - 2
	spaceCount = max(spaceCount, 0)

	spacer := strings.Repeat(" ", spaceCount)

	content := left + spacer + right

	bar := styles.StatusBarStyle.Render(content)

	return bar
}

func (m *Model) UpdateDimensions(width, height int) {
	m.Width = width
	m.Height = height
}
