package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gistio/gist-tui/pkg/tui/components/focus"
	page "github.com/gistio/gist-tui/pkg/tui/pages"
	"github.com/gistio/gist-tui/pkg/tui/pages/gitPage"
)

const (
	StashPage   page.PageID = "splash"
	GitListPage page.PageID = "git-list-page"
)

type window struct {
	Height,
	Width int
}

type model struct {
	window
	*focus.Focus
	currentPage page.PageID
	pages       map[page.PageID]tea.Model
	view        *tea.View
}

func (m *model) Init() tea.Cmd {
	page, _ := m.getCurrentPage()
	return page.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.window.Width = msg.Width
		m.window.Height = msg.Height
	case tea.KeyMsg:
		if msg.Key().Text == "q" {
			return m, tea.Quit
		}
	case tea.FocusMsg:
		return m, tea.Quit
	}

	cmds := make([]tea.Cmd, 2)
	page, _ := m.getCurrentPage()
	m.Focus, cmds[0] = m.Focus.Update(msg)
	m.pages[m.currentPage], cmds[1] = page.Update(msg)
	return m, tea.Sequence(cmds...)
}

func (m *model) View() tea.View {
	page, ok := m.getCurrentPage()
	if !ok {
		m.view.SetContent("404")
	} else {
		m.view.SetContent(page.View().Content)
	}
	return *m.view
}

func (m *model) getCurrentPage() (tea.Model, bool) {
	page, ok := m.pages[m.currentPage]
	return page, ok
}

func NewModel(cwd string, command string) *model {
	pages := map[page.PageID]tea.Model{}
	i := gitPage.New(cwd, command)
	pages[GitListPage] = &i

	return &model{
		view: &tea.View{
			AltScreen: true,
			MouseMode: tea.MouseModeCellMotion,
		},
		Focus:       &focus.Focus{},
		currentPage: GitListPage,
		pages:       pages,
	}
}
