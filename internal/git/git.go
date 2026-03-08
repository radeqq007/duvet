package git

import (
	"os/exec"
	"strings"
)

// TODO: add file status to the tree, diffs in the file preview etc.
type Status struct {
	Branch string
}

func GetStatus(dir string) *Status {
	branch, err := run(dir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return &Status{Branch: ""}
	}

	return &Status{Branch: branch}
}

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
