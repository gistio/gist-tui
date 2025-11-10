package gitPage

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/gistio/gist-tui/pkg/git"
	"github.com/gistio/gist-tui/pkg/link"
	"github.com/gistio/gist-tui/pkg/tui/components/focus"
)

type PanelType = string

type GitCommand = string

const (
	ListPanel       PanelType  = "list"
	DiffPanel       PanelType  = "diff"
	StashCommand    GitCommand = "stash"
	LogCommand      GitCommand = "log"
	TextinputID     string     = "search-input"
	ListPanelWidth  int        = 35
	FooterHeight               = 2
	HeaderHeight               = 1
	TextinputHeight            = 1
)

var blurColor = lipgloss.Color("#949494ff")
var focusColor = lipgloss.Color("#ffffffff")
var primaryColor = lipgloss.Color("#F25D94")

type gitCommand struct {
	command GitCommand
	log     git.Git
}

type GitPage struct {
	focusId string
	gitCommand
	height, width int
	git           git.Git
	textinput     *textinput.Model
	viewport      map[string]*viewport.Model
	focusPanel    PanelType
	help          help.Model

	cursor int
	limit  int
	total  int
}

func (g *GitPage) Init() tea.Cmd {
	g.git.LogList()
	g.git.StashList()
	g.git.GetRemoteUrl()
	g.textinput.Placeholder = "Search..."
	g.textinput.SetWidth(ListPanelWidth - 2)
	var cmd tea.Cmd
	g, cmd = g.updateListPanel()
	g.viewport[ListPanel].MouseWheelEnabled = true
	g.viewport[ListPanel].KeyMap = viewport.KeyMap{}
	return tea.Sequence(
		cmd,
		focus.FocusCmd(TextinputID),
	)
}

func (g *GitPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	var cmd tea.Cmd
	*g.textinput, cmd = g.textinput.Update(msg)
	cmds = append(cmds, cmd)

	if view, ok := g.viewport[g.focusPanel]; ok {
		var cmd tea.Cmd
		*g.viewport[g.focusPanel], cmd = view.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		g, cmd = g.SetSize(msg.Height, msg.Width)
		cmds = append(cmds, cmd)
	case focus.FocusMsg:
		var cmd tea.Cmd
		g.focusId = msg.NodeId
		cmds = append(cmds, g.handleFocus(msg.NodeId))
		g, cmd = g.updateListPanel()
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		var cmd tea.Cmd
		g, cmd = g.updateListPanel()
		cmds = append(cmds, cmd)
		switch true {
		case key.Matches(msg, keys.Help):
			g.help.ShowAll = !g.help.ShowAll
		case key.Matches(msg, keys.Search):
			cmds = append(cmds, focus.FocusCmd(TextinputID))
			return g, tea.Sequence(cmds...)
		case key.Matches(msg, keys.Enter):
			if g.focusPanel == DiffPanel {
				g.focusPanel = ListPanel
			} else {
				g.focusPanel = DiffPanel
			}
		case key.Matches(msg, keys.NextPage) && !g.textinput.Focused():
			g = g.handleNextPage()
		case key.Matches(msg, keys.PrevPage) && !g.textinput.Focused():
			g = g.handlePrevPage()
		case key.Matches(msg, keys.Open) && !g.textinput.Focused():
			if g.focusId != "" {
				link.Open(fmt.Sprintf("%s/commit/%s", g.git.Remote.Url, g.focusId))
			}
		case key.Matches(msg, keys.Apply) && !g.textinput.Focused() && g.command == StashCommand:
			g.git.StashApplyIndex(0)
		}
	}
	return g, tea.Sequence(cmds...)
}

func (g *GitPage) handleNextPage() *GitPage {
	nextPage := g.cursor + 1
	if g.total > nextPage*g.limit {
		g.cursor = nextPage
	}
	return g
}
func (g *GitPage) handlePrevPage() *GitPage {
	prevPage := g.cursor - 1
	if prevPage > -1 {
		g.cursor = prevPage
	}
	return g
}

func (g *GitPage) updateListPanel() (*GitPage, tea.Cmd) {
	var cmd tea.Cmd
	if g.command == StashCommand {
		g, cmd = g.updateStashList()
	} else if g.command == LogCommand {
		g, cmd = g.updateLogList()
	}

	return g, cmd
}

func (g *GitPage) updateStashList() (*GitPage, tea.Cmd) {
	searchTerm := g.textinput.Value()
	var list []git.Stash
	if searchTerm != "" {
		list = g.git.StashSearch(searchTerm)
		if g.textinput.Focused() {
			g.cursor = 0
		}
	} else {
		list = g.git.StashList()
	}

	g.total = len(list)
	start := g.cursor * g.limit
	end := min(g.cursor*g.limit+g.limit, len(list))
	base := lipgloss.NewStyle()
	content := []string{}

	focusIds := []string{}
	focusIds = append(focusIds, TextinputID)
	for _, item := range list[start:end] {
		focusIds = append(focusIds, item.ID)
		message := strings.Trim(item.Message, " ")
		if item.ID == g.focusId {
			content = append(content, base.Foreground(primaryColor).Render(message))
			diff := g.git.StashDiff(item.Index)
			g.viewport[DiffPanel].SetContent(diff)
		} else {
			content = append(content, base.Render(message))
		}
	}

	g.viewport[ListPanel].SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
	return g, focus.SetFocusNodesCmd(focusIds)
}

func (g *GitPage) updateLogList() (*GitPage, tea.Cmd) {
	searchTerm := g.textinput.Value()
	var list []git.Commit
	if searchTerm != "" {
		list = g.git.LogSearch(searchTerm)
		if g.textinput.Focused() {
			g.cursor = 0
		}
	} else {
		list = g.git.LogList()
	}

	g.total = len(list)
	start := g.cursor * g.limit
	end := min(g.cursor*g.limit+g.limit, len(list))
	base := lipgloss.NewStyle().Width(ListPanelWidth + 1)
	content := []string{}

	focusIds := []string{}
	focusIds = append(focusIds, TextinputID)
	for _, item := range list[start:end] {
		focusIds = append(focusIds, item.CommitHash)
		messageLines := strings.Split(item.Message, "\n")
		message := strings.Trim(strings.Join(messageLines, " "), " ")
		if len(message) > ListPanelWidth {
			message = message[0:ListPanelWidth]
		}
		if item.CommitHash == g.focusId {
			content = append(content, base.Background(primaryColor).Render(message))
			diff := g.git.ShowCommitDiff(item.CommitHash)
			g.viewport[DiffPanel].SetContent(diff)
		} else {
			content = append(content, base.Render(message))
		}
	}

	g.viewport[ListPanel].SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
	return g, focus.SetFocusNodesCmd(focusIds)
}

func (g *GitPage) handleFocus(nodeId string) tea.Cmd {
	if nodeId == TextinputID {
		g.textinput.Focus()
		keys.Apply.SetEnabled(false)
		return textinput.Blink
	} else if g.textinput.Focused() {
		if g.command == StashCommand && keys.Apply.Enabled() == false {
			keys.Apply.SetEnabled(true)
		}
		g.textinput.Blur()
	}
	return nil
}

func (g *GitPage) View() tea.View {
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		g.ListPanelView(),
		g.DiffPanelView(),
	)
	layout := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, g.HeaderView(), mainContent, g.FooterView()))
	return layout
}

func (g *GitPage) SetSize(h, w int) (*GitPage, tea.Cmd) {
	cmds := []tea.Cmd{}
	g.height = h
	g.width = w
	if panel, ok := g.viewport[ListPanel]; ok {
		panel.SetWidth(ListPanelWidth + 1)
		panel.SetHeight(h - 3 - FooterHeight - TextinputHeight - HeaderHeight)
	}
	if panel, ok := g.viewport[DiffPanel]; ok {
		panel.SetWidth(w - (ListPanelWidth + 5))
		panel.SetHeight(h - 1 - FooterHeight - HeaderHeight)
	}
	return g, tea.Batch(cmds...)
}

func New(cwd string, command GitCommand) GitPage {
	git := git.Git{
		WorkDir: cwd,
	}
	t := textinput.New()
	page := GitPage{
		git:        git,
		focusPanel: ListPanel,
		textinput:  &t,
		help:       help.New(),
		gitCommand: gitCommand{
			command: command,
		},
		viewport: map[string]*viewport.Model{
			ListPanel: {},
			DiffPanel: {},
		},
		limit:  50,
		cursor: 0,
	}
	return page
}
