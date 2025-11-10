package gitPage

import (
	"charm.land/lipgloss/v2"
)

func (g *GitPage) FooterView() string {
	h := g.help.View(keys)

	return lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Render(h)
}
