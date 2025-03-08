package component

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/thefishhat/tamago/server"
)

func TestConstructFieldPath(t *testing.T) {
	t.Parallel()

	type testCase struct {
		component     server.ComponentResponse
		componentType server.ComponentType
		selectedIndex int
		currPath      string
		expectedPath  string
	}

	testCases := []testCase{
		{
			component: server.ComponentResponse{
				Value: false,
				Type:  server.ComponentTypePrimitive,
			},
			componentType: server.ComponentTypePrimitive,
			selectedIndex: 0,
			currPath:      "PersistedPath",
			expectedPath:  "PersistedPath",
		},
		{
			component: server.ComponentResponse{
				Value: map[string]interface{}{
					"PersistedPath": false,
				},
				Type: server.ComponentTypeObject,
			},
			componentType: server.ComponentTypeObject,
			selectedIndex: 0,
			currPath:      "PersistedPath",
			expectedPath:  "PersistedPath.PersistedPath",
		},
		{
			component: server.ComponentResponse{
				Value: map[string]interface{}{
					"PersistedPath": false,
					"AnotherPath":   false,
				},
				Type: server.ComponentTypeObject,
			},
			componentType: server.ComponentTypeObject,
			selectedIndex: 1,
			currPath:      "PersistedPath",
			expectedPath:  "PersistedPath.AnotherPath",
		},
		{
			component: server.ComponentResponse{
				Value: []interface{}{
					false,
					false,
				},
				Type: server.ComponentTypeSlice,
			},
			componentType: server.ComponentTypeSlice,
			selectedIndex: 1,
			currPath:      "PersistedPath",
			expectedPath:  "PersistedPath[1]",
		},
		{
			component: server.ComponentResponse{
				Value: []interface{}{
					[]interface{}{false, false},
					[]interface{}{false, false},
				},
				Type: server.ComponentTypeSlice,
			},
			componentType: server.ComponentTypeSlice,
			selectedIndex: 1,
			currPath:      "PersistedPath[1]",
			expectedPath:  "PersistedPath[1][1]",
		},
		{
			component: server.ComponentResponse{
				Value: nil,
				Type:  server.ComponentTypeNil,
			},
			componentType: server.ComponentTypeNil,
			selectedIndex: 0,
			currPath:      "PersistedPath.AnotherPath[1]",
			expectedPath:  "PersistedPath.AnotherPath[1]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedPath, func(t *testing.T) {
			l := getListFromComponent(&tc.component)
			l.Select(tc.selectedIndex)

			actual := constructFieldPath(l, tc.componentType, tc.currPath)
			if actual != tc.expectedPath {
				t.Errorf("Expected %s, got %s", tc.expectedPath, actual)
			}
		})
	}
}

func getListFromComponent(component *server.ComponentResponse) list.Model {
	items := formatComponentAsItems(component)
	return list.New(items, list.NewDefaultDelegate(), 0, 0)
}
