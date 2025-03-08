package component

import (
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type itemDelegate struct {
	defaultDelegate *list.DefaultDelegate
	help            []key.Binding
}

func (d *itemDelegate) Height() int {
	return d.defaultDelegate.Height()
}

func (d *itemDelegate) Render(w io.Writer, m list.Model, index int, i list.Item) {
	d.defaultDelegate.Render(w, m, index, i)
}

func (d *itemDelegate) Spacing() int {
	return d.defaultDelegate.Spacing()
}

func (d *itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.defaultDelegate.Update(msg, m)
}

func (d *itemDelegate) ShortHelp() []key.Binding {
	return d.help
}

func (d *itemDelegate) LongHelp() [][]key.Binding {
	return [][]key.Binding{d.help}
}

func (d *itemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{d.help}
}

type delegateKeyMap struct {
	back    key.Binding
	choose  key.Binding
	edit    key.Binding
	refresh key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("escape"),
			key.WithHelp("[esc]", "back"),
		),
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("[enter]", "view"),
		),
		edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("[e]", "edit"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("[r]", "refresh"),
		),
	}
}

func newItemDelegate(items []list.Item) list.ItemDelegate {
	keys := newDelegateKeyMap()
	listDelegate := list.NewDefaultDelegate()
	d := &itemDelegate{defaultDelegate: &listDelegate}

	if len(items) == 0 {
		return d
	}

	d.help = []key.Binding{}
	if len(items) == 1 {
		d.help = append(d.help, keys.edit)
	} else {
		d.help = append(d.help, keys.choose)
	}
	d.help = append(d.help, keys.refresh, keys.back)

	return d
}
