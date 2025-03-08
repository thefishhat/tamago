package editor

import (
	"fmt"

	"github.com/thefishhat/tamago/config"
	"github.com/thefishhat/tamago/inspector"
	"github.com/thefishhat/tamago/server"
	"github.com/thefishhat/tamago/store"
	"github.com/yohamta/donburi/ecs"
)

// Editor is currently empty, but will be extended in the future.
type Editor struct{}

// Attach creates an in-memory store to format and cache the ECS data.
// It also creates an inspector that periodically updates the store with the latest ECS data.
// Finally, it starts a server that can be accessed using a CLI client.
//
// The editor can be configured using env variables. See [config.Config].
func Attach(ecs *ecs.ECS) (*Editor, error) {
	cfg := config.LoadConfig()

	editor := &Editor{}

	store := store.NewStore(ecs)

	_, err := inspector.Start(store)
	if err != nil {
		return nil, fmt.Errorf("starting inspector: %w", err)
	}

	_, err = server.Start(store, server.Config{
		Addr: cfg.Addr,
	})
	if err != nil {
		return nil, fmt.Errorf("starting server: %w", err)
	}

	return editor, nil
}
