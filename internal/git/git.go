package git

import (
	"image/color"
	"os/exec"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
)

// TODO: add file status to the tree, diffs in the file preview etc.
type Status struct {
	Root   string
	Branch string
	Files  map[string]string // path: status
}

func GetStatus(dir string) *Status {
	root, err := run(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		// Not inside a git repository
		return &Status{Branch: "", Files: nil, Root: ""}
	}

	branch, _ := run(dir, "rev-parse", "--abbrev-ref", "HEAD")

	filesStatus, _ := run(dir, "status", "--porcelain")
	lines := strings.Split(filesStatus, "\n")

	files := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			continue
		}

		status := line[:2]

		path := strings.TrimSpace(line[2:])

		files[filepath.Join(root, path)] = status
	}

	return &Status{Branch: branch, Files: files, Root: root}
}

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func ColorStatus(status string) string {
	var c color.Color
	switch status {
	case "??":
		c = lipgloss.Green

	case "!!":
		c = lipgloss.Color("#6c6c6c")

	case "A ", "AM", "A?":
		c = lipgloss.BrightGreen

	case " M", "M ", "MM":
		c = lipgloss.BrightYellow

	case " D", "D ", "MD":
		c = lipgloss.BrightRed

	case "R ", "RM":
		c = lipgloss.BrightCyan

	case "C ":
		c = lipgloss.Cyan

	case "UU", "U ", "UD", "DU":
		c = lipgloss.Magenta

	default:
		return status
	}

	return lipgloss.NewStyle().Foreground(c).Render(status)
}
