package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/radeqq007/duvet/internal/config"
	"github.com/radeqq007/duvet/internal/model"
)

func main() {
	cfg, _ := config.Get()

	p := tea.NewProgram(model.New(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
