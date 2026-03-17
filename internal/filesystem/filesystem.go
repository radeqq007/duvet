package filesystem

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
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
	defer func() { _ = dir.Close() }()

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

	return sortFiles(allFiles), nil
}

func sortFiles(files []FileNode) []FileNode {
	// TODO: add different sortings (by name, date, size etc.)
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}

		return files[i].Name < files[j].Name
	})

	return files
}

func ReadFileContent(file string) ([]byte, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if IsBinary(content) {
		return nil, fmt.Errorf("binary file")
	}

	return content, nil
}

func GetFileSize(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "error: " + err.Error()
	}

	size := float64(fileInfo.Size())
	units := []string{"B", "KB", "MB", "GB", "TB"}

	i := 0
	for size >= 1024 && i < len(units)-1 {
		size /= 1024
		i++
	}

	return fmt.Sprintf("%.2f %s", size, units[i])
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

func Highlight(content, filename, theme string) string {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, content, filename, "terminal256", theme)
	if err != nil {
		return content
	}
	return strings.TrimRight(buf.String(), "\n")
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func CreateDirectory(path string) error {
	return os.Mkdir(path, 0755)
}

func CreateFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}