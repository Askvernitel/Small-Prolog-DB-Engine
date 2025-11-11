package executor

import (
	"fmt"
	"weird/db/engine/ast"
	"weird/db/engine/client"
	"weird/db/engine/parser"
)

type DbExecutor interface {
	ExecuteQuery(query string) ([]*client.Response, error)
}
type StubDbExecutor struct {
}

func (s *StubDbExecutor) ExecuteQuery(query string) ([]*client.Response, error) {

}

type Executor struct {
	client client.DbClient
}

// NewExecutor creates a new executor with a database client
func NewExecutor(dbClient client.DbClient) *Executor {
	/*
		resp, err := dbClient.CreateTable("Users", []string{"name", "email", "password"})
		if err != nil {
			fmt.Println("ERROR OCCURED DURING CREATION" + err.Error())
		}
		fmt.Println(resp)*/
	return &Executor{
		client: dbClient,
	}
}

// Execute executes an AST statement and returns the response
func (e *Executor) Execute(stmt ast.Statement) (*client.Response, error) {
	switch s := stmt.(type) {
	case *ast.SELECTQueryStatement:
		return e.executeSelect(s)
	case *ast.INSERTStatement:
		return e.executeInsert(s)
	case *ast.UPDATEStatement:
		return e.executeUpdate(s)
	case *ast.DELETEStatement:
		return e.executeDelete(s)
	default:
		return nil, fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

func (e *Executor) executeSelect(stmt *ast.SELECTQueryStatement) (*client.Response, error) {
	return e.client.SelectAll(stmt.Table)
}

func (e *Executor) executeInsert(stmt *ast.INSERTStatement) (*client.Response, error) {
	values := make([]interface{}, len(stmt.Values))
	for i, v := range stmt.Values {
		values[i] = cleanValue(v)
	}
	return e.client.Insert(stmt.Table, values)
}

// executeUpdate executes an UPDATE statement
func (e *Executor) executeUpdate(stmt *ast.UPDATEStatement) (*client.Response, error) {
	// Convert assignments to map[string]interface{}
	set := make(map[string]interface{})
	for col, val := range stmt.Assignments {
		set[col] = cleanValue(val)
	}

	// Build WHERE clause if present
	var where map[string]interface{}
	if stmt.WhereColumn != "" {
		where = map[string]interface{}{
			stmt.WhereColumn: cleanValue(stmt.WhereValue),
		}
	}

	return e.client.Update(stmt.Table, set, where)
}

// executeDelete executes a DELETE statement
func (e *Executor) executeDelete(stmt *ast.DELETEStatement) (*client.Response, error) {
	// Build WHERE clause if present
	var where map[string]interface{}
	if stmt.WhereColumn != "" {
		where = map[string]interface{}{
			stmt.WhereColumn: cleanValue(stmt.WhereValue),
		}
	}

	if where == nil {
		// No WHERE clause means delete all
		return e.client.DeleteAll(stmt.Table)
	}

	return e.client.Delete(stmt.Table, where)
}

// cleanValue removes quotes from string literals and converts to appropriate type
func cleanValue(value string) interface{} {
	// Remove surrounding quotes if present
	if len(value) >= 2 {
		if (value[0] == '\'' && value[len(value)-1] == '\'') ||
			(value[0] == '"' && value[len(value)-1] == '"') {
			return value[1 : len(value)-1]
		}
	}

	// Try to detect if it's a number
	// Simple check: if it contains only digits and possibly a decimal point
	if isNumeric(value) {
		return value // Return as string, let JSON marshaling handle it
	}

	return value
}

// isNumeric checks if a string represents a numeric value
func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}

	hasDot := false
	for i, c := range s {
		if c == '-' && i == 0 {
			continue
		}
		if c == '.' {
			if hasDot {
				return false
			}
			hasDot = true
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// ExecuteProgram executes all statements in a program
func (e *Executor) ExecuteProgram(program *ast.Program) ([]*client.Response, error) {
	results := make([]*client.Response, 0, len(program.Statements))

	for _, stmt := range program.Statements {
		resp, err := e.Execute(stmt)
		if err != nil {
			return results, fmt.Errorf("failed to execute statement '%s': %w", stmt.String(), err)
		}
		results = append(results, resp)
	}

	return results, nil
}

// ExecuteMultiple executes multiple statements and returns all results
func (e *Executor) ExecuteMultiple(statements []ast.Statement) ([]*client.Response, error) {
	results := make([]*client.Response, 0, len(statements))

	for _, stmt := range statements {
		resp, err := e.Execute(stmt)
		if err != nil {
			return results, fmt.Errorf("failed to execute statement '%s': %w", stmt.String(), err)
		}
		results = append(results, resp)
	}

	return results, nil
}
func (e *Executor) ExecuteQuery(q string) ([]*client.Response, error) {
	program, err := parser.ParseString(q)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}
	return e.ExecuteProgram(program)
}

// Close closes the executor and its underlying client
func (e *Executor) Close() error {
	return e.client.Close()
}
