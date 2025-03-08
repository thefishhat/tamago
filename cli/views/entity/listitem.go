package entity

import (
	"fmt"

	"github.com/thefishhat/tamago/server"
)

type entityItem struct {
	server.Component
}

func (i entityItem) Title() string       { return i.Name }
func (i entityItem) Description() string { return i.Type + " " + fmt.Sprintf("%v", i.Value) }
func (i entityItem) FilterValue() string { return i.Name + i.Type }
