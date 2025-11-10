package gitPage

import (
	"charm.land/bubbles/v2/key"
)

type keyMap struct {
	Enter    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	PrevPage key.Binding
	NextPage key.Binding
	Search   key.Binding
	Open     key.Binding
	Help     key.Binding
	Quit     key.Binding
	Apply    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter},
		{k.Tab},
		{k.ShiftTab},
		{k.PrevPage},
		{k.NextPage},
		{k.Apply},
		{k.Open},
		{k.Help},
		{k.Quit},
	}
}

var keys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "panel scroll focus"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "move down"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "move up"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next page"),
	),
	PrevPage: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("shift+n", "prev page"),
	),
	Search: key.NewBinding(
		key.WithKeys(":"),
		key.WithHelp(":", "search"),
	),
	Open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open github"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Apply: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "apply stash"),
		key.WithDisabled(),
	),
}
