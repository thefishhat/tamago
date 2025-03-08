package entities

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/thefishhat/tamago/cli/views/entity"
	"github.com/thefishhat/tamago/server"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Client interface {
	GetEntities() (*server.ListEntitiesResponse, error)
	GetEntity(entityID string) (*server.GetEntityResponse, error)
	GetComponent(entityID string, componentName string, fieldPath string) (*server.ComponentResponse, error)
	SetComponent(entityID string, componentName string, fieldPath string, value interface{}) error
}

type EntitiesModel struct {
	list   list.Model
	client Client
}

func NewEntitiesModel(client Client) *EntitiesModel {
	response, err := client.GetEntities()
	if err != nil {
		log.Fatal("fetching entities: ", err)
	}

	items := formatEntitiesAsItems(response.Entities)

	delegate := newItemDelegate()
	list := list.New(items, delegate, 0, 0)
	list.Title = "Entities"

	return &EntitiesModel{
		list:   list,
		client: client,
	}
}

func (m *EntitiesModel) Init() tea.Cmd {
	return nil
}

func (m *EntitiesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "r":
			response, _ := m.client.GetEntities()
			items := formatEntitiesAsItems(response.Entities)
			m.list.SetItems(items)
		case "enter":
			selected, ok := m.list.SelectedItem().(entitiesItem)
			if !ok {
				break
			}
			return m, func() tea.Msg {
				return entity.Open(m.client, selected.Id)
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

func (m *EntitiesModel) View() string {
	return docStyle.Render(m.list.View())
}

func formatEntitiesAsItems(entities []server.EntitySummary) []list.Item {
	items := make([]list.Item, 0, len(entities))
	for _, entity := range entities {
		items = append(items, entitiesItem{entity})
	}
	return items
}
