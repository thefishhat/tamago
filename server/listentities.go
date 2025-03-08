package server

import (
	"encoding/json"
	"net/http"
)

type ListEntitiesResponse struct {
	Entities []EntitySummary `json:"entities"`
}

func (s *Server) listEntitiesHandler(w http.ResponseWriter, _ *http.Request) {
	var response ListEntitiesResponse
	for _, entry := range s.store.GetEntries() {
		var entity EntitySummary = entitySummaryFromEntry(entry)
		response.Entities = append(response.Entities, entity)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
