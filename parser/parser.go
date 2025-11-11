package parser

import (
	"fmt"
	"weird/db/engine/ast"
	"weird/db/engine/lexer"
	"weird/db/engine/token"
)

type QueryParser interface {
	Parse() (*ast.Program, error)
}
type Parser struct {
	tokens  []token.Token
	pos     int
	current token.Token
}

func New(tokens []token.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    0,
	}
	if len(tokens) > 0 {
		p.current = tokens[0]
	}
	return p
}

func (p *Parser) advance() {
	p.pos++
	if p.pos < len(p.tokens) {
		p.current = p.tokens[p.pos]
	}
}

func (p *Parser) peek() *token.Token {
	if p.pos+1 < len(p.tokens) {
		return &p.tokens[p.pos+1]
	}
	return nil
}

func (p *Parser) expect(tokenType token.TokenType) error {
	if p.current.Token != tokenType {
		return fmt.Errorf("expected %s, got %s", tokenType, p.current.Token)
	}
	p.advance()
	return nil
}

func (p *Parser) skipWhitespace() {
	for p.pos < len(p.tokens) && (p.current.Token == token.ENDLINE_TOKEN || p.current.Token == token.SEMICOLON_TOKEN) {
		p.advance()
	}
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{
		Statements: make([]ast.Statement, 0),
	}

	for p.pos < len(p.tokens) {
		p.skipWhitespace()

		if p.pos >= len(p.tokens) {
			break
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.skipWhitespace()
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.current.Token {
	case token.SELECT_TOKEN:
		return p.parseSELECTStatement()
	case token.INSERT_TOKEN:
		return p.parseINSERTStatement()
	case token.UPDATE_TOKEN:
		return p.parseUPDATEStatement()
	case token.DELETE_TOKEN:
		return p.parseDELETEStatement()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.current.Literal)
	}
}

func (p *Parser) parseSELECTStatement() (*ast.SELECTQueryStatement, error) {
	if err := p.expect(token.SELECT_TOKEN); err != nil {
		return nil, err
	}

	fields := make([]string, 0)

	for {
		p.skipWhitespace()

		if p.current.Token != token.IDENT_TOKEN {
			return nil, fmt.Errorf("expected field name, got %s", p.current.Token)
		}

		fields = append(fields, p.current.Literal)
		p.advance()

		p.skipWhitespace()

		if p.current.Token == token.COMMA_TOKEN {
			p.advance()
			continue
		}

		break
	}

	p.skipWhitespace()

	if err := p.expect(token.FROM_TOKEN); err != nil {
		return nil, err
	}

	p.skipWhitespace()

	if p.current.Token != token.IDENT_TOKEN {
		return nil, fmt.Errorf("expected table name, got %s", p.current.Token)
	}

	tableName := p.current.Literal
	p.advance()

	return ast.NewSELECTQueryStatement(fields, tableName), nil
}

func (p *Parser) parseINSERTStatement() (*ast.INSERTStatement, error) {
	if err := p.expect(token.INSERT_TOKEN); err != nil {
		return nil, err
	}

	p.skipWhitespace()

	if err := p.expect(token.INTO_TOKEN); err != nil {
		return nil, err
	}

	p.skipWhitespace()

	if p.current.Token != token.IDENT_TOKEN {
		return nil, fmt.Errorf("expected table name, got %s", p.current.Token)
	}

	tableName := p.current.Literal
	p.advance()
	p.skipWhitespace()

	columns := make([]string, 0)
	if p.current.Token == token.LPAREN_TOKEN {
		p.advance()
		p.skipWhitespace()

		for {
			if p.current.Token != token.IDENT_TOKEN {
				return nil, fmt.Errorf("expected column name, got %s", p.current.Token)
			}

			columns = append(columns, p.current.Literal)
			p.advance()
			p.skipWhitespace()

			if p.current.Token == token.COMMA_TOKEN {
				p.advance()
				p.skipWhitespace()
				continue
			}

			break
		}

		if err := p.expect(token.RPAREN_TOKEN); err != nil {
			return nil, err
		}
		p.skipWhitespace()
	}

	if err := p.expect(token.VALUES_TOKEN); err != nil {
		return nil, err
	}

	p.skipWhitespace()

	if err := p.expect(token.LPAREN_TOKEN); err != nil {
		return nil, err
	}

	p.skipWhitespace()

	// Parse values
	values := make([]string, 0)
	for {
		if p.current.Token != token.STRING_TOKEN && p.current.Token != token.NUMBER_TOKEN && p.current.Token != token.IDENT_TOKEN {
			return nil, fmt.Errorf("expected value, got %s", p.current.Token)
		}

		values = append(values, p.current.Literal)
		p.advance()
		p.skipWhitespace()

		if p.current.Token == token.COMMA_TOKEN {
			p.advance()
			p.skipWhitespace()
			continue
		}

		break
	}

	if err := p.expect(token.RPAREN_TOKEN); err != nil {
		return nil, err
	}

	return ast.NewINSERTStatement(tableName, columns, values), nil
}

func (p *Parser) parseUPDATEStatement() (*ast.UPDATEStatement, error) {
	if err := p.expect(token.UPDATE_TOKEN); err != nil {
		return nil, err
	}

	p.skipWhitespace()

	if p.current.Token != token.IDENT_TOKEN {
		return nil, fmt.Errorf("expected table name, got %s", p.current.Token)
	}

	tableName := p.current.Literal
	p.advance()
	p.skipWhitespace()

	if err := p.expect(token.SET_TOKEN); err != nil {
		return nil, err
	}

	p.skipWhitespace()

	assignments := make(map[string]string)
	for {
		if p.current.Token != token.IDENT_TOKEN {
			return nil, fmt.Errorf("expected column name, got %s", p.current.Token)
		}

		colName := p.current.Literal
		p.advance()
		p.skipWhitespace()

		if err := p.expect(token.EQUALS_TOKEN); err != nil {
			return nil, err
		}

		p.skipWhitespace()

		if p.current.Token != token.STRING_TOKEN && p.current.Token != token.NUMBER_TOKEN && p.current.Token != token.IDENT_TOKEN {
			return nil, fmt.Errorf("expected value, got %s", p.current.Token)
		}

		value := p.current.Literal
		assignments[colName] = value
		p.advance()
		p.skipWhitespace()

		if p.current.Token == token.COMMA_TOKEN {
			p.advance()
			p.skipWhitespace()
			continue
		}

		break
	}

	whereCol := ""
	whereVal := ""
	if p.current.Token == token.WHERE_TOKEN {
		p.advance()
		p.skipWhitespace()

		if p.current.Token != token.IDENT_TOKEN {
			return nil, fmt.Errorf("expected column name in WHERE clause, got %s", p.current.Token)
		}

		whereCol = p.current.Literal
		p.advance()
		p.skipWhitespace()

		if err := p.expect(token.EQUALS_TOKEN); err != nil {
			return nil, err
		}

		p.skipWhitespace()

		if p.current.Token != token.STRING_TOKEN && p.current.Token != token.NUMBER_TOKEN && p.current.Token != token.IDENT_TOKEN {
			return nil, fmt.Errorf("expected value in WHERE clause, got %s", p.current.Token)
		}

		whereVal = p.current.Literal
		p.advance()
	}

	return ast.NewUPDATEStatement(tableName, assignments, whereCol, whereVal), nil
}
func (p *Parser) parseDELETEStatement() (*ast.DELETEStatement, error) {
	if err := p.expect(token.DELETE_TOKEN); err != nil {
		return nil, err
	}

	p.skipWhitespace()

	if err := p.expect(token.FROM_TOKEN); err != nil {
		return nil, err
	}

	p.skipWhitespace()

	if p.current.Token != token.IDENT_TOKEN {
		return nil, fmt.Errorf("expected table name, got %s", p.current.Token)
	}

	tableName := p.current.Literal
	p.advance()
	p.skipWhitespace()

	// Optional WHERE clause
	whereCol := ""
	whereVal := ""
	if p.current.Token == token.WHERE_TOKEN {
		p.advance()
		p.skipWhitespace()

		if p.current.Token != token.IDENT_TOKEN {
			return nil, fmt.Errorf("expected column name in WHERE clause, got %s", p.current.Token)
		}

		whereCol = p.current.Literal
		p.advance()
		p.skipWhitespace()

		if err := p.expect(token.EQUALS_TOKEN); err != nil {
			return nil, err
		}

		p.skipWhitespace()

		if p.current.Token != token.STRING_TOKEN && p.current.Token != token.NUMBER_TOKEN && p.current.Token != token.IDENT_TOKEN {
			return nil, fmt.Errorf("expected value in WHERE clause, got %s", p.current.Token)
		}

		whereVal = p.current.Literal
		p.advance()
	}

	return ast.NewDELETEStatement(tableName, whereCol, whereVal), nil
}

func ParseSingle(tokens []token.Token) (ast.Statement, error) {
	p := New(tokens)
	program, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if len(program.Statements) == 0 {
		return nil, fmt.Errorf("no statements found")
	}

	return program.Statements[0], nil
}

func ParseString(q string) (*ast.Program, error) {
	l := lexer.Lexer{}
	t := l.Tokenize(q)
	p := New(t)
	return p.Parse()
}
