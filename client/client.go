package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/thefishhat/tamago/server"
)

// Client is an HTTP client for the [server.Server].
type Client struct {
	Addr string
}

// NewClient creates a new client for the given address.
func NewClient(addr string) *Client {
	return &Client{
		Addr: addr,
	}
}

// GetEntity fetches the entity with the given ID from the server.
func (c *Client) GetEntity(entityID string) (*server.GetEntityResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/entities/%s", c.Addr, entityID))
	if err != nil {
		return nil, fmt.Errorf("fetching entity: %w", err)
	}
	defer resp.Body.Close()

	var response server.GetEntityResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return &response, nil
}

// GetEntities fetches all entities from the server.
func (c *Client) GetEntities() (*server.ListEntitiesResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/entities", c.Addr))
	if err != nil {
		return nil, fmt.Errorf("fetching entities: %w", err)
	}
	defer resp.Body.Close()

	var response server.ListEntitiesResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return &response, nil
}

// GetComponent fetches the component with the given name from the entity with the given ID.
// If the fieldPath is not empty, it will fetch the field at the given path.
// The fieldPath is a dot-separated path to the field in the component.
// Example:
//
//	"position.x" // will fetch the x field from the position component.
//	"inventory.items[0].name" // will fetch the name field from the first item in the inventory component.
func (c *Client) GetComponent(entityID string, componentName string, fieldPath string) (*server.ComponentResponse, error) {
	componentUrl := fmt.Sprintf("http://%s/entities/%s/components/%s", c.Addr, entityID, componentName)
	if fieldPath != "" {
		componentUrl += "?field=" + fieldPath
	}
	resp, err := http.Get(componentUrl)
	if err != nil {
		return nil, fmt.Errorf("fetching component: %w", err)
	}
	defer resp.Body.Close()

	var response server.ComponentResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return &response, nil
}

// SetComponent sets the value of the field at the given path in the component with the given name.
// The value can be any JSON-serializable value.
// Example:
//
//	client.SetComponent("1", "position", "x", 10) // sets the x field in the position component to 10.
func (c *Client) SetComponent(entityID string, componentName string, fieldPath string, value interface{}) error {
	url := fmt.Sprintf("http://%s/entities/%s/components/%s?field=%s", c.Addr, entityID, componentName, fieldPath)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	buf := []byte(fmt.Sprintf(`{"value": %s}`, value))
	req.Body = io.NopCloser(bytes.NewReader(buf))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed request: %d", res.StatusCode)
	}

	return nil
}
