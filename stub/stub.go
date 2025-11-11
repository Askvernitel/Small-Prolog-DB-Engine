package stub

import "weird/db/engine/client"

const (
	DB_EXECUTOR = iota
)

type Stub interface {
}

func GetStub(stubId int) Stub {
	switch stubId {
	case DB_EXECUTOR:
		return &StubDbExecutor{}
	}
	return nil
}

type StubDbExecutor struct {
}

func (s *StubDbExecutor) ExecuteQuery(_ string) ([]*client.Response, error) {
	return []*client.Response{
		{
			Status:  "success",
			Message: "Query executed successfully",
			Table:   "users",
			Columns: []string{"id", "name", "email", "created_at"},
			Rows: []Row{
				{"id": 1, "name": "Alice Smith", "email": "alice@example.com", "created_at": "2024-01-15"},
				{"id": 2, "name": "Bob Jones", "email": "bob@example.com", "created_at": "2024-02-20"},
				{"id": 3, "name": "Carol White", "email": "carol@example.com", "created_at": "2024-03-10"},
			},
			ID:    10,
			Count: 3,
		},
	}, nil
}
