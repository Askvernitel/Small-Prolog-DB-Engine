package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type CreateTableRequest struct {
	Type    string   `json:"type"`
	Table   string   `json:"table"`
	Columns []string `json:"columns"`
}

type InsertRequest struct {
	Type   string        `json:"type"`
	Table  string        `json:"table"`
	Values []interface{} `json:"values"`
}

type SelectRequest struct {
	Type  string                 `json:"type"`
	Table string                 `json:"table"`
	Where map[string]interface{} `json:"where,omitempty"`
}

type UpdateRequest struct {
	Type  string                 `json:"type"`
	Table string                 `json:"table"`
	Set   map[string]interface{} `json:"set"`
	Where map[string]interface{} `json:"where,omitempty"`
}

type DeleteRequest struct {
	Type  string                 `json:"type"`
	Table string                 `json:"table"`
	Where map[string]interface{} `json:"where,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) sendRequest(payload interface{}) (*Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/query", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	fmt.Println(response)
	if !response.Success {
		return &response, fmt.Errorf("query failed: %s", response.Error)
	}

	return &response, nil
}

func (c *Client) CreateTable(table string, columns []string) (*Response, error) {
	req := CreateTableRequest{
		Type:    "create_table",
		Table:   table,
		Columns: columns,
	}
	return c.sendRequest(req)
}

func (c *Client) Insert(table string, values []interface{}) (*Response, error) {
	req := InsertRequest{
		Type:   "insert",
		Table:  table,
		Values: values,
	}
	return c.sendRequest(req)
}

func (c *Client) Select(table string, where map[string]interface{}) (*Response, error) {
	req := SelectRequest{
		Type:  "select",
		Table: table,
		Where: where,
	}
	return c.sendRequest(req)
}

func (c *Client) SelectAll(table string) (*Response, error) {
	return c.Select(table, nil)
}

func (c *Client) Update(table string, set map[string]interface{}, where map[string]interface{}) (*Response, error) {
	req := UpdateRequest{
		Type:  "update",
		Table: table,
		Set:   set,
		Where: where,
	}
	return c.sendRequest(req)
}

func (c *Client) Delete(table string, where map[string]interface{}) (*Response, error) {
	req := DeleteRequest{
		Type:  "delete",
		Table: table,
		Where: where,
	}
	return c.sendRequest(req)
}

func (c *Client) DeleteAll(table string) (*Response, error) {
	return c.Delete(table, nil)
}

func (c *Client) Ping() error {
	req, err := http.NewRequest("GET", c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create ping request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("server unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}
