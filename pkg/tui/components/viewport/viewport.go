package viewport

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type Viewport struct {
	Viewport map[string]*viewport.Model
	FocusId  string
}

type ViewContentMsg = struct {
	key     string
	content string
	Height,
	Width int
}

type EnableMouseWheelMsg = struct {
	key string
}

func ViewContentCmd(key, content string, width, height int) tea.Cmd {
	return func() tea.Msg {
		return ViewContentMsg{key: key, content: content, Width: width, Height: height}
	}
}

func EnableMouseWheelCmd(key string) tea.Cmd {
	return func() tea.Msg {
		return EnableMouseWheelMsg{key}
	}
}

func (v *Viewport) Update(msg tea.Msg) (*Viewport, tea.Cmd) {
	switch msg := msg.(type) {
	case ViewContentMsg:
		view, ok := v.Viewport[msg.key]
		if !ok {
			newView := viewport.Model{}
			newView.KeyMap = viewport.KeyMap{}

			v.Viewport[msg.key] = &newView
			view = &newView
		}
		if v.FocusId == msg.key {
			view.MouseWheelEnabled = true
		} else {
			view.MouseWheelEnabled = false
		}

		view.SetContent(msg.content)
		if msg.Width > 0 {
			view.SetWidth(msg.Width)
		}
		if msg.Height > 0 {
			view.SetHeight(msg.Height)
		}
	case EnableMouseWheelMsg:
		v.FocusId = msg.key
	}

	cmds := make([]tea.Cmd, len(v.Viewport))
	index := 0
	for key, view := range v.Viewport {
		var m viewport.Model
		if !view.MouseWheelEnabled {
			continue
		}
		m, cmds[index] = view.Update(msg)
		v.Viewport[key] = &m
		index++
	}
	return v, tea.Batch(cmds...)
}
