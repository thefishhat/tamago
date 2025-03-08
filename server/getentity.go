package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/yohamta/donburi"
)

type EntitySummary struct {
	Id        string           `json:"id"`
	Name      string           `json:"name"`
	Archetype ArchetypeSummary `json:"archetype"`
}

type Component struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type Entity struct {
	EntitySummary
	Components []Component `json:"components"`
}

type GetEntityResponse struct {
	Entity Entity `json:"entity"`
}

func (s *Server) getEntityHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
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

	var summary EntitySummary = entitySummaryFromEntry(entry)
	var entity Entity
	entity.EntitySummary = summary
	entity.Components = getComponentsFromEntry(entry)
	response := GetEntityResponse{Entity: entity}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}
}

func entitySummaryFromEntry(entry *donburi.Entry) EntitySummary {
	var entity EntitySummary
	entity.Id = fmt.Sprintf("%d", entry.Id())
	entity.Name = entry.String()
	entity.Archetype.EntityCount = len(entry.Archetype().Entities())
	for _, components := range entry.Archetype().ComponentTypes() {
		entity.Archetype.Components = append(entity.Archetype.Components, struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}{
			Name: components.Name(),
			Type: components.Typ().Name(),
		})
	}
	return entity
}

func getComponentsFromEntry(entry *donburi.Entry) []Component {
	var components []Component
	componentTypes := entry.Archetype().ComponentTypes()
	for _, componentType := range componentTypes {
		ptr := entry.Component(componentType)
		component := reflect.Indirect(reflect.NewAt(componentType.Typ(), ptr))
		fields := recursivelyConstructValue(component, 1)

		resp := Component{
			Name:  componentType.Name(),
			Type:  componentType.Typ().Name(),
			Value: fields,
		}
		components = append(components, resp)
	}
	return components
}
