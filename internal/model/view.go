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
	if m.Layout.Width == 0 {
		return "loading..."
	}

	leftPane := m.RenderLeftPane()
	rightPane := m.RenderRightPane()
	bar := m.RenderStatusBar()

	view := lipgloss.JoinHorizontal(lipgloss.Bottom, leftPane, rightPane)
	view = lipgloss.JoinVertical(lipgloss.Left, view, bar)

	switch m.IO.Mode {
	case mode.Command:
		content := ":" + m.IO.CmdInput
		if strings.HasPrefix(m.IO.CmdInput, "!") {
			content = "$" + m.IO.CmdInput[1:]
		}
		cmdBox := styles.CmdBoxStyle.Render(content)

		x := m.Layout.Width/2 - lipgloss.Width(cmdBox)/2
		y := m.Layout.Height/2 - lipgloss.Height(cmdBox)/2

		view = ui.PlaceOverlay(x, y, cmdBox, view)

	case mode.Alert:
		var alertBox string
		switch m.IO.Alert.Type {
		case alert.Normal:
			alertBox = styles.AlertNormalStyle.Render(m.IO.Alert.Text)
		case alert.Info:
			alertBox = styles.AlertInfoStyle.Render(m.IO.Alert.Text)
		case alert.Error:
			alertBox = styles.AlertErrorStyle.Render(m.IO.Alert.Text)
		case alert.Warning:
			alertBox = styles.AlertWarningStyle.Render(m.IO.Alert.Text)
		}

		x := m.Layout.Width/2 - lipgloss.Width(alertBox)/2
		y := m.Layout.Height/2 - lipgloss.Height(alertBox)/2

		view = ui.PlaceOverlay(x, y, alertBox, view)
	}

	return view
}

func (m *Model) RenderLeftPane() string {
	var leftContent strings.Builder

	visibleHeight := m.VisibleHeight() - config.Layout.StatusBarHeight - config.Layout.BorderWidth

	start := m.Display.LeftScroll
	end := m.Display.LeftScroll + visibleHeight
	end = min(end, len(m.FileTree))

	for i := start; i < end; i++ {
		node := m.FileTree[i]

		_, isSelected := m.IO.Selected[filepath.Join(m.CurPath, node.Name)]

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
			line = styles.SelectedStyle.Width(m.Layout.Width / 2).Render(line)
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
	if m.Display.Focus == pane.Left {
		leftPane = styles.FocusedPaneStyle.
			Width(m.Layout.Width / 2).
			Height(visibleHeight).
			Render(leftContent.String())
	} else {
		leftPane = styles.PaneStyle.
			Width(m.Layout.Width / 2).
			Height(visibleHeight).
			Render(leftContent.String())
	}

	return leftPane
}

func (m *Model) RenderRightPane() string {
	var rightContent strings.Builder
	visibleHeight := m.VisibleHeight() - config.Layout.StatusBarHeight - config.Layout.BorderWidth

	if m.Display.Preview.Content != "" {
		wrapped := lipgloss.NewStyle().
			Width(m.Layout.Width/2 - config.Layout.BorderWidth*2).
			Render(m.Display.Preview.Content)

		visualLines := strings.Split(wrapped, "\n")
		start := m.Display.RightScroll
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
	if m.Display.Focus == pane.Right {
		rightPane = styles.FocusedPaneStyle.
			Width(m.Layout.Width / 2).
			Height(visibleHeight).
			Render(rightContent.String())
	} else {
		rightPane = styles.PaneStyle.
			Width(m.Layout.Width / 2).
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
	if m.Git.Branch != "" {
		left += " on \ue0a0 " + m.Git.Branch
	}

	left += seperator + strconv.Itoa(len(m.FileTree)) + " items"

	selected := len(m.IO.Selected)
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
	spaceCount := m.Layout.Width - ansi.StringWidth(
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
	m.Layout.Width = width
	m.Layout.Height = height
}
