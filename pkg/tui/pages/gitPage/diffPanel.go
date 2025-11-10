package gitPage

import (
	"charm.land/lipgloss/v2"
)

func (g *GitPage) DiffPanelView() string {
	diff := g.viewport[DiffPanel]
	color := blurColor
	faint := true
	if g.focusPanel == DiffPanel {
		color = primaryColor
		faint = false
	}
	return lipgloss.
		NewStyle().
		Faint(faint).
		BorderForeground(color).
		Border(lipgloss.RoundedBorder()).
		Render(diff.View())
}
