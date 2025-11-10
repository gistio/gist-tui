package gitPage

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
)

func (g *GitPage) PagView(color color.Color) string {
	pagBorder := lipgloss.NormalBorder()
	pagBorder.TopLeft = pagBorder.Top
	pagBorder.TopLeft = pagBorder.Top

	viewedCount := min(g.cursor*g.limit+g.limit, g.total)

	page := g.cursor + 1

	container := lipgloss.NewStyle().Height(1).
		Width(ListPanelWidth + 1).
		Border(pagBorder).
		BorderForeground(color).
		BorderBottom(false).
		BorderRight(false).
		BorderLeft(false)
	view := lipgloss.NewStyle().
		Render(fmt.Sprintf("Page %d | %d/%d", page, viewedCount, g.total))
	w, _ := lipgloss.Size(view)
	space := lipgloss.NewStyle().Width(((ListPanelWidth + 1) - w) / 2).Render()
	return container.Render(lipgloss.JoinHorizontal(lipgloss.Left, space, view))
}

func (g *GitPage) ListPanelView() string {
	panel := g.viewport[ListPanel].View()
	input := g.textinput.View()
	color := blurColor
	if g.focusPanel == ListPanel {
		color = primaryColor
	}

	listContent := lipgloss.JoinVertical(lipgloss.Left, input, panel)
	pag := g.PagView(color)

	return lipgloss.
		NewStyle().
		BorderForeground(color).
		Border(lipgloss.RoundedBorder()).
		Render(lipgloss.JoinVertical(lipgloss.Left, listContent, pag))
}
