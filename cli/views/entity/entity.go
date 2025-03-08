package entity

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/thefishhat/tamago/cli/hotswapmodel"
	component "github.com/thefishhat/tamago/cli/views/component"
	"github.com/thefishhat/tamago/server"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Client interface {
	GetEntity(entityID string) (*server.GetEntityResponse, error)
	GetComponent(entityID string, componentName string, fieldPath string) (*server.ComponentResponse, error)
	SetComponent(entityID string, componentName string, fieldPath string, value interface{}) error
}

type EntityModel struct {
	list   list.Model
	entity server.Entity
	client Client
}

func NewEntityModel(client Client, entityID string) *EntityModel {
	response, err := client.GetEntity(entityID)
	if err != nil {
		log.Fatal("fetching entity: ", err)
	}

	delegate := newItemDelegate()
	items := formatEntityAsItems(response.Entity)
	list := list.New(items, delegate, 0, 0)
	list.Title = "Entities > Entity " + entityID

	return &EntityModel{
		list:   list,
		entity: response.Entity,
		client: client,
	}
}

func (m *EntityModel) Init() tea.Cmd {
	return nil
}

func (m *EntityModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return hotswapmodel.SwitchToLastModel{} }
		case "r":
			response, _ := m.client.GetEntity(m.entity.Id)
			items := formatEntityAsItems(response.Entity)
			m.list.SetItems(items)
		case "enter":
			selected, ok := m.list.SelectedItem().(entityItem)
			if !ok {
				break
			}
			return m, func() tea.Msg {
				return component.Open(m.client, m.entity.Id, selected.Component.Name, "")
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *EntityModel) View() string {
	return docStyle.Render(
		m.list.View(),
	)
}

func formatEntityAsItems(entity server.Entity) []list.Item {
	var items []list.Item

	for _, attr := range entity.Components {
		items = append(items, entityItem{attr})
	}

	return items
}
