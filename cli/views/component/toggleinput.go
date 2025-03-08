package component

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type toggleInput struct {
	isEditing bool
	model     textinput.Model
}

type inputDone struct {
	value string
}

func newToggleInput(value interface{}) *toggleInput {
	input := textinput.New()
	input.Placeholder = fmt.Sprintf("%v", value)
	input.Focus()
	input.CharLimit = 156
	input.Width = 20
	return &toggleInput{
		isEditing: false,
		model:     input,
	}
}

func (i *toggleInput) IsEditing() bool {
	return i.isEditing
}

func (i *toggleInput) SetIsEditing(editing bool) {
	i.isEditing = editing
}

func (i *toggleInput) View() string {
	return i.model.View()
}

func (i *toggleInput) Update(msg tea.Msg) tea.Msg {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
		return inputDone{
			value: i.model.Value(),
		}
	}
	i.model, _ = i.model.Update(msg)
	return nil
}

func (i *toggleInput) Value() string {
	return i.model.Value()
}
