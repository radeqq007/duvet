package command

import tea "github.com/charmbracelet/bubbletea"

func Exec(cmd string, args ...string) tea.Cmd {
	if fn, ok := commands[cmd]; ok {
		return fn(args...)
	}

	return nil
}

var commands = map[string]func(...string) tea.Cmd{
	"q":    handleQuit,
	"quit": handleQuit,
}

func handleQuit(args ...string) tea.Cmd {
	return tea.Quit
}
