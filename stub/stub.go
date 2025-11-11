package stub

import (
	"weird/db/engine/client"
)

type StubDbExecutor struct {
}

func (s *StubDbExecutor) ExecuteQuery(_ string) ([]*client.Response, error) {
	return []*client.Response{
		{
			Status:  "success",
			Message: "Query executed successfully",
			Table:   "users",
			Columns: []string{"id", "name", "email", "created_at"},
			Rows: []client.Row{
				{ID: 1, Data: []string{"1", "Alice Smith", "alice@example.com", "2024-01-15"}},
				{ID: 2, Data: []string{"2", "Bob Jones", "bob@example.com", "2024-02-20"}},
				{ID: 3, Data: []string{"3", "Carol White", "carol@example.com", "2024-03-10"}},
			},
			ID:    10,
			Count: 3,
		},
	}, nil
}
