package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/thefishhat/tamago/inspector"
	"github.com/thefishhat/tamago/server"
	"github.com/thefishhat/tamago/store"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/ecs"
)

var testCfg = server.Config{
	Addr: ":8080",
}

type MockComponent struct{}

type ServerSuite struct {
	suite.Suite
	ecs    *ecs.ECS
	server *server.Server
	insp   *inspector.Inspector
	st     *store.Store
}

func (s *ServerSuite) SetupTest() {
	s.ecs = ecs.NewECS(donburi.NewWorld())
	s.st = store.NewStore(s.ecs)

	var err error
	s.insp, err = inspector.Start(s.st)
	require.NoError(s.T(), err)

	s.server, err = server.Start(s.st, testCfg)
	require.NoError(s.T(), err)

	err = waitForHealthyServer()
	require.NoError(s.T(), err)
}

func (s *ServerSuite) TearDownTest() {
	s.server.Stop()
	s.insp.Stop()
}

func (s *ServerSuite) AddComponents(components ...component.IComponentType) []donburi.Entity {
	var result []donburi.Entity
	for _, c := range components {
		entity := s.ecs.World.Create(c)
		result = append(result, entity)
	}

	s.insp.IntrospectECS()
	return result
}

func (s *ServerSuite) TestListArchetypes() {
	mockComponent := donburi.NewComponentType[MockComponent]()
	mockComponent.SetName("MyTestComponent")
	s.AddComponents(mockComponent)

	resp, err := http.Get("http://" + testCfg.Addr + "/")
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	var actualResp server.ListArchetypesResponse
	b, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)
	err = json.Unmarshal(b, &actualResp)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), server.ListArchetypesResponse{
		Archetypes: []server.ArchetypeSummary{
			{
				EntityCount: 1,
				Components: []server.ComponentSummary{
					{
						Name: mockComponent.Name(),
						Type: mockComponent.Typ().Name(),
					},
				},
			},
		},
	}, actualResp, "response should match expected")

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
}

func (s *ServerSuite) TestListEntities() {
	mockComponent := donburi.NewComponentType[MockComponent]()
	mockComponent.SetName("MyTestComponent")
	entities := s.AddComponents(mockComponent)
	require.Len(s.T(), entities, 1)
	entity := entities[0]
	entry := s.ecs.World.Entry(entity)

	resp, err := http.Get("http://" + testCfg.Addr + "/entities")
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	var actualResp server.ListEntitiesResponse
	b, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)
	err = json.Unmarshal(b, &actualResp)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), server.ListEntitiesResponse{
		Entities: []server.EntitySummary{
			{
				Id:   "1",
				Name: entry.String(),
				Archetype: server.ArchetypeSummary{
					EntityCount: 1,
					Components: []server.ComponentSummary{
						{
							Name: mockComponent.Name(),
							Type: mockComponent.Typ().Name(),
						},
					},
				},
			},
		},
	}, actualResp, "response should match expected")

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
}

func (s *ServerSuite) TestGetEntity() {
	mockComponent := donburi.NewComponentType[MockComponent]()
	mockComponent.SetName("MyTestComponent")
	entities := s.AddComponents(mockComponent)
	require.Len(s.T(), entities, 1)
	entity := entities[0]
	entry := s.ecs.World.Entry(entity)

	resp, err := http.Get("http://" + testCfg.Addr + fmt.Sprintf("/entities/%d", entity.Id()))
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	var actualResp server.GetEntityResponse
	b, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)
	err = json.Unmarshal(b, &actualResp)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), server.GetEntityResponse{
		Entity: server.Entity{
			EntitySummary: server.EntitySummary{
				Id:   fmt.Sprintf("%d", entry.Id()),
				Name: entry.String(),
				Archetype: server.ArchetypeSummary{
					EntityCount: 1,
					Components: []server.ComponentSummary{
						{
							Name: mockComponent.Name(),
							Type: mockComponent.Typ().Name(),
						},
					},
				},
			},
			Components: []server.Component{
				{
					Name:  mockComponent.Name(),
					Type:  mockComponent.Typ().Name(),
					Value: map[string]interface{}{},
				},
			},
		},
	}, actualResp, "response should match expected")

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
}

func (s *ServerSuite) TestGetEntityComponents() {
	mockComponent := donburi.NewComponentType[MockComponent]()
	mockComponent.SetName("MyTestComponent")
	entities := s.AddComponents(mockComponent)
	require.Len(s.T(), entities, 1)
	entity := entities[0]

	resp, err := http.Get("http://" + testCfg.Addr + fmt.Sprintf("/entities/%d", entity.Id()))
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	var entitiesResp server.GetEntityResponse
	b, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)
	err = json.Unmarshal(b, &entitiesResp)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)

	resp, err = http.Get("http://" + testCfg.Addr + fmt.Sprintf("/entities/%d/components", entity.Id()))
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	var entitiesCompResp server.GetEntityResponse
	b, err = io.ReadAll(resp.Body)
	require.NoError(s.T(), err)
	err = json.Unmarshal(b, &entitiesCompResp)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), entitiesResp, entitiesCompResp, "/entities should match /entities/components")
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
}

func (s *ServerSuite) TestGetComponent() {
	type Person struct {
		Name string
	}
	mockComponent := donburi.NewComponentType[Person]()
	mockComponent.SetName("MyPersonComponent")
	entities := s.AddComponents(mockComponent)
	require.Len(s.T(), entities, 1)
	entity := entities[0]

	resp, err := http.Get("http://" + testCfg.Addr + fmt.Sprintf("/entities/%d/components/%s", entity.Id(), mockComponent.Name()))
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	var actualResp server.ComponentResponse
	b, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)
	err = json.Unmarshal(b, &actualResp)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), server.ComponentResponse{
		Value: map[string]interface{}{
			"Name": "string",
		},
		Type: server.ComponentTypeObject,
	}, actualResp, "response should match expected")

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
}

func (s *ServerSuite) TestGetComponentField() {
	type Person struct {
		Name string
	}
	mockComponent := donburi.NewComponentType[Person]()
	mockComponent.SetName("MyPersonComponent")
	entities := s.AddComponents(mockComponent)
	require.Len(s.T(), entities, 1)
	entity := entities[0]
	mockComponent.SetValue(s.ecs.World.Entry(entity), Person{Name: "donburi"})
	s.insp.IntrospectECS()

	resp, err := http.Get("http://" + testCfg.Addr + fmt.Sprintf("/entities/%d/components/%s?field=Name", entity.Id(), mockComponent.Name()))
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	var actualResp server.ComponentResponse
	b, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)
	err = json.Unmarshal(b, &actualResp)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), server.ComponentResponse{
		Value: "\"donburi\"",
		Type:  server.ComponentTypePrimitive,
	}, actualResp, "response should match expected")

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
}

func (s *ServerSuite) TestSetComponentField() {
	type Person struct {
		Name string
	}
	mockComponent := donburi.NewComponentType[Person]()
	mockComponent.SetName("MyPersonComponent")
	entities := s.AddComponents(mockComponent)
	require.Len(s.T(), entities, 1)
	entity := entities[0]
	mockComponent.SetValue(s.ecs.World.Entry(entity), Person{Name: "donburi"})
	s.insp.IntrospectECS()

	var setReq = server.SetComponentRequest{
		Value: "tamago",
	}
	b, err := json.Marshal(setReq)
	require.NoError(s.T(), err)

	req, err := http.NewRequest(http.MethodPut, "http://"+testCfg.Addr+fmt.Sprintf("/entities/%d/components/%s?field=Name", entity.Id(), mockComponent.Name()), bytes.NewReader(b))
	require.NoError(s.T(), err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
	s.insp.IntrospectECS()
	v := (*Person)(s.ecs.World.Entry(entity).Component(mockComponent))
	assert.Equal(s.T(), "tamago", v.Name)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

func waitForHealthyServer() error {
	// exponential backoff
	// 50ms, 100ms, 200ms, 400ms, 800ms, 1600ms, 3200ms, 6400ms, 12800ms, 25600ms
	delay := 50 * time.Millisecond
	for i := 0; i < 10; i++ {
		resp, err := http.Get("http://" + testCfg.Addr + "/healthcheck")
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		time.Sleep(delay)
		delay *= 2
	}

	return fmt.Errorf("server unhealthy after 10 retries")
}
