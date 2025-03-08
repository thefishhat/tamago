package entities

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type delegateKeyMap struct {
	choose  key.Binding
	refresh key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.refresh,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.refresh,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("[enter]", "view"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("[r]", "refresh"),
		),
	}
}

func newItemDelegate() list.DefaultDelegate {
	keys := newDelegateKeyMap()
	d := list.NewDefaultDelegate()
	help := []key.Binding{keys.choose, keys.refresh}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	return d
}
