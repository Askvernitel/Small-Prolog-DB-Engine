package ast

type Statement interface {
	Statement()
	String() string
}

type QueryStatement interface {
	Statement
	QueryStatement()
}

type SELECTQueryStatement struct {
	Fields []string
	Table  string
}

func (s *SELECTQueryStatement) Statement() {}

func (s *SELECTQueryStatement) QueryStatement() {}

func (s *SELECTQueryStatement) String() string {
	fields := ""
	for i, f := range s.Fields {
		if i > 0 {
			fields += ", "
		}
		fields += f
	}
	return "SELECT " + fields + " FROM " + s.Table
}

func NewSELECTQueryStatement(fields []string, table string) *SELECTQueryStatement {
	return &SELECTQueryStatement{
		Fields: fields,
		Table:  table,
	}
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	result := ""
	for _, stmt := range p.Statements {
		result += stmt.String() + "\n"
	}
	return result
}
