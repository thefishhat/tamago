package server

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
)

type SetComponentRequest struct {
	Value interface{} `json:"value"`
}

// req: /entities/3/components/PlayerData?field=IgnorePlatform
// body: {"value": true}
// resp: 200, etc.
func (s *Server) setComponentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("entity_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid entity ID", http.StatusBadRequest)
		return
	}

	entry := s.store.GetEntry(uint32(id))
	if entry == nil {
		http.Error(w, "Entity not found", http.StatusNotFound)
		return
	}

	fieldPath := r.URL.Query().Get("field")
	componentName := r.PathValue("component_name")

	var component reflect.Value
	for _, componentType := range entry.Archetype().ComponentTypes() {
		if componentType.Name() == componentName {
			ptr := entry.Component(componentType)
			component = reflect.Indirect(reflect.NewAt(componentType.Typ(), ptr))
		}
	}

	// Read the request body and decode into SetComponentRequest
	var req SetComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Pass the value from the request body into SetField
	err = SetField(component, fieldPath, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
