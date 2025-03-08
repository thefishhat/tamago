package component

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/thefishhat/tamago/cli/hotswapmodel"
	"github.com/thefishhat/tamago/server"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type Client interface {
	GetComponent(entityID string, componentName string, fieldPath string) (*server.ComponentResponse, error)
	SetComponent(entityID string, componentName string, fieldPath string, value interface{}) error
}

type ComponentModel struct {
	list          list.Model
	entityID      string
	componentName string
	componentType server.ComponentType
	fieldPath     string
	client        Client
}

func NewComponentModel(client Client, entityID string, componentName string, fieldPath string) *ComponentModel {
	response, err := client.GetComponent(entityID, componentName, fieldPath)
	if err != nil {
		log.Fatal("fetching component:", err)
	}

	items := formatComponentAsItems(response)
	delegate := newItemDelegate(items)
	list := list.New(items, delegate, 0, 0)
	list.Title = "Entities > Entity " + entityID + " > " + componentName
	if fieldPath != "" {
		list.Title += " : " + fieldPath
	}

	return &ComponentModel{
		list:          list,
		entityID:      entityID,
		componentName: componentName,
		componentType: response.Type,
		fieldPath:     fieldPath,
		client:        client,
	}
}

func (m *ComponentModel) Init() tea.Cmd {
	return nil
}

func (m *ComponentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	selectedItem, ok := m.list.SelectedItem().(componentItem)
	if !ok {
		return m, nil
	}

	if selectedItem.input.IsEditing() {
		inputMsg := selectedItem.input.Update(msg)
		if inputDone, ok := inputMsg.(inputDone); ok {
			selectedItem.input.SetIsEditing(false)
			err := m.setValue(inputDone.value)
			if err != nil {
				selectedItem.errMsg.SetMsg(err.Error())
			} else {
				m.reloadItems()
			}
			m.list, _ = m.list.Update(inputMsg)
			msg = nil
		} else {
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return hotswapmodel.SwitchToLastModel{} }
		case "r":
			return m, func() tea.Msg {
				m.reloadItems()
				return nil
			}
		case "enter":
			return m, func() tea.Msg {
				newFieldPath := constructFieldPath(m.list, m.componentType, m.fieldPath)
				if newFieldPath == m.fieldPath {
					return m
				}
				return Open(m.client, m.entityID, m.componentName, newFieldPath)
			}
		case "e":
			return m, func() tea.Msg {
				if len(m.list.Items()) == 1 {
					selectedItem.input.SetIsEditing(true)
				}
				return nil
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

func (m *ComponentModel) View() string {
	return docStyle.Render(
		m.list.View(),
	)
}

func (m *ComponentModel) setValue(value string) error {
	err := m.client.SetComponent(m.entityID, m.componentName, m.fieldPath, value)
	if err != nil {
		return fmt.Errorf("setting component: %w", err)
	}
	return nil
}

func (m *ComponentModel) reloadItems() {
	response, err := m.client.GetComponent(m.entityID, m.componentName, m.fieldPath)
	if err != nil {
		log.Fatal("fetching component:", err)
	}
	items := formatComponentAsItems(response)
	m.list.SetItems(items)
}

func constructFieldPath(l list.Model, componentType server.ComponentType, currPath string) string {
	if l.SelectedItem() == nil {
		return currPath
	}

	selectedItem, ok := l.SelectedItem().(componentItem)
	if !ok {
		return currPath
	}

	switch componentType {
	case server.ComponentTypePrimitive, server.ComponentTypeNil:
		return currPath

	case server.ComponentTypeObject:
		var res string
		if currPath != "" {
			res = currPath + "."
		}
		return res + selectedItem.name

	case server.ComponentTypeSlice:
		return currPath + "[" + strconv.Itoa(l.Index()) + "]"

	default:
		return currPath
	}
}

func formatComponentAsItems(component *server.ComponentResponse) []list.Item {
	var items []list.Item

	switch component.Type {
	case server.ComponentTypeObject:
		var obj map[string]interface{}
		obj = component.Value.(map[string]interface{})
		for key, value := range obj {
			items = append(items, newComponentItem(key, value))
		}
	case server.ComponentTypeSlice:
		var arr []interface{}
		arr = component.Value.([]interface{})
		for i, value := range arr {
			items = append(items, newComponentItem(fmt.Sprintf("[%d]", i), value))
		}
	case server.ComponentTypePrimitive:
		items = append(items, newComponentItem("value", component.Value))
	case server.ComponentTypeNil:
		items = append(items, newComponentItem("value", ""))
	}

	return items
}
