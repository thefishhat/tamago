package entity

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thefishhat/tamago/cli/hotswapmodel"
)

type open struct {
	EntityId string
	Client   Client
}

func Open(client Client, entityId string) hotswapmodel.ModelSwapper {
	return open{
		EntityId: entityId,
		Client:   client,
	}
}

func (msg open) GetModel() tea.Model {
	return NewEntityModel(msg.Client, msg.EntityId)
}
