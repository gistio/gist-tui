package gitPage

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

func (g *GitPage) HeaderView() string {
	title := fmt.Sprintf("Gist * %s", g.command)
	header := lipgloss.
		NewStyle().
		Render(title)
	w, _ := lipgloss.Size(header)
	space := lipgloss.NewStyle().Width((g.width - w) / 2).Render("")
	return lipgloss.NewStyle().Render(space, header)
}
