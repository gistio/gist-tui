package focus

import (
	tea "charm.land/bubbletea/v2"
)

type Focus struct {
	nodes      []string
	activeNode int
}

func (f Focus) FocusNext() int {
	if len(f.nodes) == 0 {
		return -1
	}

	index := f.activeNode + 1
	if index > len(f.nodes)-1 {
		index = 0
	}

	return index
}

func (f Focus) FocusPrev() int {
	if len(f.nodes) == 0 {
		return -1
	}

	index := f.activeNode - 1
	if index < 0 {
		index = len(f.nodes) - 1
	}
	return index
}

func (f Focus) FocusNodes() []string {
	return f.nodes
}
func (f Focus) FocusNode() string {
	if len(f.nodes) == 0 {
		return ""
	}

	return f.nodes[f.activeNode]
}

func (f Focus) FindNodeIndex(nodeId string) int {
	if len(f.nodes) == 0 || nodeId == "" {
		return -1
	}

	for index, node := range f.nodes {
		if node == nodeId {
			return index
		}
	}

	return -1
}

type FetFocusNodesMsg struct {
	nodes []string
}

func SetFocusNodesCmd(nodes []string) tea.Cmd {
	return func() tea.Msg {
		return FetFocusNodesMsg{nodes: nodes}
	}
}

type FocusMsg struct {
	NodeId string
}

func FocusCmd(nodeId string) tea.Cmd {
	return func() tea.Msg {
		return FocusMsg{NodeId: nodeId}
	}
}

type BlurMsg struct {
	nodeId string
}

func BlurCmd(nodeId string) tea.Cmd {
	return func() tea.Msg {
		return BlurMsg{nodeId: nodeId}
	}
}

func (f *Focus) Update(msg tea.Msg) (*Focus, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case FocusMsg:
		index := f.FindNodeIndex(msg.NodeId)
		if index > -1 {
			f.activeNode = index
		}
	case FetFocusNodesMsg:
		f.nodes = msg.nodes
		break
	case tea.KeyMsg:
		key := msg.Key()
		if key.Mod == tea.ModShift && key.Code == tea.KeyTab {
			index := f.FocusPrev()
			f.activeNode = index
		} else if msg.Key().Code == tea.KeyTab {
			index := f.FocusNext()
			f.activeNode = index
		}
		cmds = append(cmds, FocusCmd(f.FocusNode()))
	}
	return f, tea.Batch(cmds...)
}
