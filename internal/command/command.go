package command

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Msg struct {
	Name string
	Args []string
}

func Exec(input string) tea.Cmd {
	return func() tea.Msg {
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return nil
		}

		return Msg{
			Name: parts[0],
			Args: parts[1:],
		}
	}
}

// var commands = map[string]func(...string) tea.Cmd{
// 	"q":    handleQuit,
// 	"quit": handleQuit,
// }

// func handleQuit(args ...string) tea.Cmd {
// 	return tea.Quit
// }
