package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/gistio/gist-tui/pkg/tui"
)

func getPage() string {
	args := os.Args
	if args[1] != "" {
		return args[1]
	}
	return "log"
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic("No dir")
	}
	model := tui.NewModel(cwd, getPage())
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
