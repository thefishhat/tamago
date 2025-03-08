package entities

import "github.com/thefishhat/tamago/server"

type entitiesItem struct {
	server.EntitySummary
}

func (i entitiesItem) Title() string       { return i.Id }
func (i entitiesItem) Description() string { return i.Name }
func (i entitiesItem) FilterValue() string { return i.Name }
