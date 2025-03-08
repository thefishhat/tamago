package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func TestNewStore(t *testing.T) {
	ecs := &ecs.ECS{}
	store := NewStore(ecs)

	assert.NotNil(t, store)
	assert.Equal(t, ecs, store.ecs)
	assert.NotNil(t, store.entries)
}

func TestStore_GetWorld(t *testing.T) {
	ecs := &ecs.ECS{}
	store := NewStore(ecs)

	assert.Equal(t, ecs.World, store.GetWorld())
}

func TestStore_GetEntry(t *testing.T) {
	ecs := &ecs.ECS{}
	store := NewStore(ecs)
	entry := &donburi.Entry{}
	store.entries[1] = entry

	assert.Equal(t, entry, store.GetEntry(1))
	assert.Nil(t, store.GetEntry(2))
}

func TestStore_GetEntries(t *testing.T) {
	ecs := &ecs.ECS{}
	store := NewStore(ecs)
	entry1 := &donburi.Entry{}
	entry2 := &donburi.Entry{}
	store.entries[1] = entry1
	store.entries[2] = entry2

	entries := store.GetEntries()
	assert.Equal(t, 2, len(entries))
	assert.Equal(t, entry1, entries[1])
	assert.Equal(t, entry2, entries[2])
}

func TestStore_SetEntries(t *testing.T) {
	ecs := &ecs.ECS{}
	store := NewStore(ecs)
	entry1 := &donburi.Entry{}
	entry2 := &donburi.Entry{}
	entries := map[uint32]*donburi.Entry{
		1: entry1,
		2: entry2,
	}
	store.SetEntries(entries)

	assert.Equal(t, entries, store.GetEntries())
}
