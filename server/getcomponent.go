package server

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
)

type ComponentType string

const (
	ComponentTypePrimitive ComponentType = "primitive"
	ComponentTypeObject    ComponentType = "object"
	ComponentTypeSlice     ComponentType = "slice"
	ComponentTypeNil       ComponentType = "nil"
)

type ComponentResponse struct {
	Value interface{}   `json:"value"`
	Type  ComponentType `json:"type"`
}

// req: /entities/3/components/PlayerData?field=IgnorePlatform
// resp: {"value": false, "type": "primitive"}
func (s *Server) getComponentHandler(w http.ResponseWriter, r *http.Request) {
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

	field, err := GetField(component, fieldPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := ComponentResponse{
		Value: field,
		Type:  reflectToComponentType(field),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}
}

func reflectToComponentType(v interface{}) ComponentType {
	if v == nil {
		return ComponentTypeNil
	}
	switch v.(type) {
	case []interface{}:
		return ComponentTypeSlice
	case map[string]interface{}:
		return ComponentTypeObject
	default:
		return ComponentTypePrimitive
	}
}
