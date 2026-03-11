package git

import (
	"os/exec"
	"path/filepath"
	"strings"
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
		status := strings.TrimSpace(line[:2])
		path := line[3:]
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
