package inspector

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/yohamta/donburi"
)

type Store interface {
	GetWorld() donburi.World
	GetEntries() map[uint32]*donburi.Entry
	SetEntries(map[uint32]*donburi.Entry)
}

type Inspector struct {
	store     Store
	scheduler gocron.Scheduler
}

// The inspector will periodically introspect the ECS and update the store with the latest entries.
//
// The interval of introspection is 3 seconds.
func Start(store Store) (*Inspector, error) {
	log := log.New(log.Writer(), "[inspector] ", log.LstdFlags)

	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("creating scheduler: %w", err)
	}

	inspector := &Inspector{
		store:     store,
		scheduler: s,
	}

	_, err = s.NewJob(
		gocron.DurationJob(
			3*time.Second,
		),
		gocron.NewTask(inspector.IntrospectECS),
	)
	if err != nil {
		return nil, fmt.Errorf("creating ECS job: %w", err)
	}
	log.Println("Created ECS introspection job")

	s.Start()
	log.Println("Inspector started")

	return inspector, nil
}

// IntrospectECS inspects the ECS and updates the store with the latest entries.
func (i *Inspector) IntrospectECS() {
	world := i.store.GetWorld()
	entries := i.store.GetEntries()
	newEntries := make(map[uint32]*donburi.Entry, len(entries))
	for _, arch := range world.Archetypes() {
		entities := arch.Entities()
		for _, entity := range entities {
			entry := world.Entry(entity)
			id := uint32(entry.Id())
			newEntries[id] = entry
		}
	}
	i.store.SetEntries(newEntries)
}

// Stop stops the inspector.
func (i *Inspector) Stop() {
	i.scheduler.Shutdown()
}
