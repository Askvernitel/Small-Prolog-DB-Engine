package ast

import "strings"

// Statement is the base interface for all AST statements
type Statement interface {
	Statement()
	String() string // Useful for debugging and printing
}

// QueryStatement represents statements that query data
type QueryStatement interface {
	Statement
	QueryStatement()
}

// DMLStatement represents data manipulation statements
type DMLStatement interface {
	Statement
	DMLStatement()
}

// SELECTQueryStatement represents a SELECT query
type SELECTQueryStatement struct {
	Fields []string // Column names to select (or "*" for all)
	Table  string   // Table name to select from
}

// Statement implements the Statement interface
func (s *SELECTQueryStatement) Statement() {}

// QueryStatement implements the QueryStatement interface
func (s *SELECTQueryStatement) QueryStatement() {}

// String returns a string representation of the SELECT statement
func (s *SELECTQueryStatement) String() string {
	fields := strings.Join(s.Fields, ", ")
	return "SELECT " + fields + " FROM " + s.Table
}

// NewSELECTQueryStatement creates a new SELECT query statement
func NewSELECTQueryStatement(fields []string, table string) *SELECTQueryStatement {
	return &SELECTQueryStatement{
		Fields: fields,
		Table:  table,
	}
}

// INSERTStatement represents an INSERT INTO statement
type INSERTStatement struct {
	Table   string   // Table name
	Columns []string // Column names (optional)
	Values  []string // Values to insert
}

// Statement implements the Statement interface
func (i *INSERTStatement) Statement() {}

// DMLStatement implements the DMLStatement interface
func (i *INSERTStatement) DMLStatement() {}

// String returns a string representation of the INSERT statement
func (i *INSERTStatement) String() string {
	result := "INSERT INTO " + i.Table
	if len(i.Columns) > 0 {
		result += " (" + strings.Join(i.Columns, ", ") + ")"
	}
	result += " VALUES (" + strings.Join(i.Values, ", ") + ")"
	return result
}

// NewINSERTStatement creates a new INSERT statement
func NewINSERTStatement(table string, columns []string, values []string) *INSERTStatement {
	return &INSERTStatement{
		Table:   table,
		Columns: columns,
		Values:  values,
	}
}

// UPDATEStatement represents an UPDATE statement
type UPDATEStatement struct {
	Table       string            // Table name
	Assignments map[string]string // Column = Value pairs
	WhereColumn string            // WHERE clause column (optional)
	WhereValue  string            // WHERE clause value (optional)
}

// Statement implements the Statement interface
func (u *UPDATEStatement) Statement() {}

// DMLStatement implements the DMLStatement interface
func (u *UPDATEStatement) DMLStatement() {}

// String returns a string representation of the UPDATE statement
func (u *UPDATEStatement) String() string {
	result := "UPDATE " + u.Table + " SET "

	assignments := make([]string, 0, len(u.Assignments))
	for col, val := range u.Assignments {
		assignments = append(assignments, col+" = "+val)
	}
	result += strings.Join(assignments, ", ")

	if u.WhereColumn != "" {
		result += " WHERE " + u.WhereColumn + " = " + u.WhereValue
	}

	return result
}

// NewUPDATEStatement creates a new UPDATE statement
func NewUPDATEStatement(table string, assignments map[string]string, whereCol string, whereVal string) *UPDATEStatement {
	return &UPDATEStatement{
		Table:       table,
		Assignments: assignments,
		WhereColumn: whereCol,
		WhereValue:  whereVal,
	}
}

// DELETEStatement represents a DELETE FROM statement
type DELETEStatement struct {
	Table       string // Table name
	WhereColumn string // WHERE clause column (optional)
	WhereValue  string // WHERE clause value (optional)
}

// Statement implements the Statement interface
func (d *DELETEStatement) Statement() {}

// DMLStatement implements the DMLStatement interface
func (d *DELETEStatement) DMLStatement() {}

// String returns a string representation of the DELETE statement
func (d *DELETEStatement) String() string {
	result := "DELETE FROM " + d.Table

	if d.WhereColumn != "" {
		result += " WHERE " + d.WhereColumn + " = " + d.WhereValue
	}

	return result
}

// NewDELETEStatement creates a new DELETE statement
func NewDELETEStatement(table string, whereCol string, whereVal string) *DELETEStatement {
	return &DELETEStatement{
		Table:       table,
		WhereColumn: whereCol,
		WhereValue:  whereVal,
	}
}

// Program represents the root AST node containing all statements
type Program struct {
	Statements []Statement
}

// String returns a string representation of the program
func (p *Program) String() string {
	result := ""
	for _, stmt := range p.Statements {
		result += stmt.String() + "\n"
	}
	return result
}
