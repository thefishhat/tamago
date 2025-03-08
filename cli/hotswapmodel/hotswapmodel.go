package hotswapmodel

import (
	tea "github.com/charmbracelet/bubbletea"
)

type HotSwapModel struct {
	activeModel   tea.Model
	modelStack    []tea.Model
	width, height int
}

func New(activeModel tea.Model) *HotSwapModel {
	return &HotSwapModel{
		activeModel: activeModel,
	}
}

func (m HotSwapModel) Init() tea.Cmd {
	if m.activeModel == nil {
		return nil
	}
	return m.activeModel.Init()
}

func (m HotSwapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case ModelSwapper:
		m.swapWithHistory(msg.GetModel())
		return m, nil
	case SwitchToLastModel:
		if len(m.modelStack) > 0 {
			m.activeModel = m.modelStack[len(m.modelStack)-1]
			m.modelStack = m.modelStack[:len(m.modelStack)-1]
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.activeModel == nil {
		return m, nil
	}

	var cmd tea.Cmd
	m.activeModel, cmd = m.activeModel.Update(msg)
	return m, cmd
}

func (m HotSwapModel) View() string {
	if m.activeModel == nil {
		return "Waiting for active model to be set..."
	}
	return m.activeModel.View()
}

func (m *HotSwapModel) swap(model tea.Model) {
	resizedModel, _ := model.Update(tea.WindowSizeMsg{
		Width:  m.width,
		Height: m.height,
	})
	m.activeModel = resizedModel
}

func (m *HotSwapModel) pushCurrentModel() {
	if m.activeModel != nil {
		m.modelStack = append(m.modelStack, m.activeModel)
	}
}

func (m *HotSwapModel) swapWithHistory(model tea.Model) {
	m.pushCurrentModel()
	m.swap(model)
}
