package hotswapmodel

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ModelSwapper interface {
	GetModel() tea.Model
}

type SwitchToLastModel struct{}
