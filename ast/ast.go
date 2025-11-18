package ast

import (
	"strings"
)

type Statement interface {
	Statement()
	String() string // Useful for debugging and printing
}

type QueryStatement interface {
	Statement
	QueryStatement()
}

type DMLStatement interface {
	Statement
	DMLStatement()
}

type SELECTQueryStatement struct {
	Fields []string // Column names to select (or "*" for all)
	Table  string   // Table name to select from
}

func (s *SELECTQueryStatement) Statement() {}

func (s *SELECTQueryStatement) QueryStatement() {}

func (s *SELECTQueryStatement) String() string {
	fields := strings.Join(s.Fields, ", ")
	return "SELECT " + fields + " FROM " + s.Table
}

func NewSELECTQueryStatement(fields []string, table string) *SELECTQueryStatement {
	return &SELECTQueryStatement{
		Fields: fields,
		Table:  table,
	}
}

type INSERTStatement struct {
	Table   string   // Table name
	Columns []string // Column names (optional)
	Values  []string // Values to insert
}

func (i *INSERTStatement) Statement() {}

func (i *INSERTStatement) DMLStatement() {}

func (i *INSERTStatement) String() string {
	result := "INSERT INTO " + i.Table
	if len(i.Columns) > 0 {
		result += " (" + strings.Join(i.Columns, ", ") + ")"
	}
	result += " VALUES (" + strings.Join(i.Values, ", ") + ")"
	return result
}

func NewINSERTStatement(table string, columns []string, values []string) *INSERTStatement {
	return &INSERTStatement{
		Table:   table,
		Columns: columns,
		Values:  values,
	}
}

type UPDATEStatement struct {
	Table       string            // Table name
	Assignments map[string]string // Column = Value pairs
	WhereColumn string            // WHERE clause column (optional)
	WhereValue  string            // WHERE clause value (optional)
}

func (u *UPDATEStatement) Statement() {}

func (u *UPDATEStatement) DMLStatement() {}

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

func NewUPDATEStatement(table string, assignments map[string]string, whereCol string, whereVal string) *UPDATEStatement {
	return &UPDATEStatement{
		Table:       table,
		Assignments: assignments,
		WhereColumn: whereCol,
		WhereValue:  whereVal,
	}
}

type DELETEStatement struct {
	Table       string // Table name
	WhereColumn string // WHERE clause column (optional)
	WhereValue  string // WHERE clause value (optional)
}

func (d *DELETEStatement) Statement() {}

func (d *DELETEStatement) DMLStatement() {}

func (d *DELETEStatement) String() string {
	result := "DELETE FROM " + d.Table

	if d.WhereColumn != "" {
		result += " WHERE " + d.WhereColumn + " = " + d.WhereValue
	}

	return result
}

func NewDELETEStatement(table string, whereCol string, whereVal string) *DELETEStatement {
	return &DELETEStatement{
		Table:       table,
		WhereColumn: whereCol,
		WhereValue:  whereVal,
	}
}

type CREATEStatement struct {
	Table   string
	Columns []string
}

func (d *CREATEStatement) Statement() {}

func (d *CREATEStatement) DMLStatement() {}

// String returns a string representation of the DELETE statement
func (d *CREATEStatement) String() string {
	result := "CREATE" + d.Table
	for _, col := range d.Columns {
		result += " " + col
	}
	return result
}

func NewCREATEStatement(table string, columns []string) *CREATEStatement {
	return &CREATEStatement{
		Table:   table,
		Columns: columns,
	}
}

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
