package component

import (
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	errMsgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type errorMsg struct {
	timer timer.Model
	msg   string
}

func newErrorMsg() *errorMsg {
	return &errorMsg{}
}

func (e *errorMsg) View() string {
	return errMsgStyle.Render(e.msg)
}

func (e *errorMsg) SetMsg(msg string) {
	e.msg = msg
	e.timer = timer.New(3)
	e.timer.Init()
}

func (e *errorMsg) Update(msg tea.Msg) {
	e.timer, _ = e.timer.Update(msg)

	if e.timer.Timedout() {
		e.msg = ""
	}
}
