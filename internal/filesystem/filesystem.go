package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileNode struct {
	Name  string
	IsDir bool
}

func GetFiles(path string) ([]FileNode, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	entries, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	allFiles := make([]FileNode, len(entries))

	for i, entry := range entries {
		allFiles[i] = FileNode{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		}
	}

	return allFiles, nil
}

func ReadFileContent(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	if IsBinary(content) {
		return "", fmt.Errorf("binary file")
	}

	text := string(content)
	return text, err
}

func ParentDir(path string) string {
	return filepath.Dir(path)
}

func IsBinary(content []byte) bool {
	checkBytes := 512
	checkLen := min(checkBytes, len(content))

	for i := range checkLen {
		if content[i] == 0 {
			return true
		}
	}
	return false
}
