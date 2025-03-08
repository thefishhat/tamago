package server

import (
	"encoding/json"
	"net/http"
)

type ComponentSummary struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ArchetypeSummary struct {
	EntityCount int                `json:"entity_count"`
	Components  []ComponentSummary `json:"components"`
}

type ListArchetypesResponse struct {
	Archetypes []ArchetypeSummary `json:"archetypes"`
}

func (s *Server) listArchetypesHandler(w http.ResponseWriter, _ *http.Request) {
	var response ListArchetypesResponse
	for _, arch := range s.store.GetWorld().Archetypes() {
		entities := arch.Entities()
		if len(entities) == 0 {
			continue
		}
		var archetype ArchetypeSummary
		archetype.EntityCount = len(entities)
		for _, components := range arch.ComponentTypes() {
			archetype.Components = append(archetype.Components, struct {
				Name string `json:"name"`
				Type string `json:"type"`
			}{
				Name: components.Name(),
				Type: components.Typ().Name(),
			})
		}
		response.Archetypes = append(response.Archetypes, archetype)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
