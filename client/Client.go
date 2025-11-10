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

// Prolog DB Response format
type Response struct {
	Status  string   `json:"status"`
	Message string   `json:"message,omitempty"`
	Table   string   `json:"table,omitempty"`
	Columns []string `json:"columns,omitempty"`
	Rows    []Row    `json:"rows,omitempty"`
	ID      int      `json:"id,omitempty"`
	Count   int      `json:"count,omitempty"`
}

type Row struct {
	ID   int           `json:"id"`
	Data []interface{} `json:"data"`
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
	fmt.Println(payload)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	fmt.Println(string(jsonData))
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
	fmt.Println(string(body))
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Status != "success" {
		return &response, fmt.Errorf("query failed: %s", response.Message)
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

func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// Helper method to get row data as a map
func (r *Row) AsMap(columns []string) map[string]interface{} {
	result := make(map[string]interface{})
	for i, col := range columns {
		if i < len(r.Data) {
			result[col] = r.Data[i]
		}
	}
	return result
}

// Helper method to print response in a readable format
func (r *Response) Print() {
	if r.Status == "success" {
		fmt.Printf("✓ Success: %s\n", r.Message)

		if r.ID > 0 {
			fmt.Printf("  ID: %d\n", r.ID)
		}

		if r.Count > 0 {
			fmt.Printf("  Count: %d\n", r.Count)
		}

		if len(r.Rows) > 0 {
			fmt.Printf("  Table: %s\n", r.Table)
			fmt.Printf("  Columns: %v\n", r.Columns)
			fmt.Printf("  Rows (%d):\n", len(r.Rows))
			for _, row := range r.Rows {
				fmt.Printf("    ID %d: %v\n", row.ID, row.Data)
			}
		}
	} else {
		fmt.Printf("✗ Error: %s\n", r.Message)
	}
}
