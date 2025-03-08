package component

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thefishhat/tamago/cli/hotswapmodel"
)

type open struct {
	EntityID      string
	ComponentName string
	FieldPath     string
	Client        Client
}

func Open(client Client, entityID, componentName, fieldPath string) hotswapmodel.ModelSwapper {
	return open{
		EntityID:      entityID,
		ComponentName: componentName,
		FieldPath:     fieldPath,
		Client:        client,
	}
}

func (msg open) GetModel() tea.Model {
	return NewComponentModel(msg.Client, msg.EntityID, msg.ComponentName, msg.FieldPath)
}
