package store

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// Store is an in-memory cache for the ECS data.
type Store struct {
	ecs     *ecs.ECS
	entries map[uint32]*donburi.Entry
}

// NewStore creates a new store for the given ECS.
func NewStore(ecs *ecs.ECS) *Store {
	return &Store{
		ecs:     ecs,
		entries: make(map[uint32]*donburi.Entry),
	}
}

// GetWorld returns the ECS world the store was created against.
func (s *Store) GetWorld() donburi.World {
	return s.ecs.World
}

// GetEntry returns the entry with the given ID from the store.
func (s *Store) GetEntry(id uint32) *donburi.Entry {
	if entry, ok := s.entries[id]; ok {
		return entry
	}
	return nil
}

// GetEntries returns all entries in the store.
func (s *Store) GetEntries() map[uint32]*donburi.Entry {
	return s.entries
}

// SetEntries sets the entries in the store.
func (s *Store) SetEntries(entries map[uint32]*donburi.Entry) {
	s.entries = entries
}
