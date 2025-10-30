package lexer

import (
	"bytes"
	"strings"
	"unicode"
	"weird/db/engine/token"
)

type Lexer struct {
	ReadBuffer bytes.Buffer
	tokens     []token.Token
}

func New() *Lexer {
	return &Lexer{
		tokens: make([]token.Token, 0),
	}
}

func (l *Lexer) isEndByte(char rune) bool {
	return char == ' ' || char == '\n' || char == '\t' || char == '\r'
}

func (l *Lexer) flushBuffer() {
	if l.ReadBuffer.Len() == 0 {
		return
	}

	literal := l.ReadBuffer.String()
	l.ReadBuffer.Reset()

	// Determine token type based on literal
	var tok token.Token
	tok.Literal = literal

	// Check if it's a number
	if isNumber(literal) {
		tok.Token = token.NUMBER_TOKEN
		l.tokens = append(l.tokens, tok)
		return
	}

	// Check for keywords (case-insensitive)
	switch strings.ToUpper(literal) {
	case "SELECT":
		tok.Token = token.SELECT_TOKEN
	case "FROM":
		tok.Token = token.FROM_TOKEN
	case "INSERT":
		tok.Token = token.INSERT_TOKEN
	case "INTO":
		tok.Token = token.INTO_TOKEN
	case "VALUES":
		tok.Token = token.VALUES_TOKEN
	case "UPDATE":
		tok.Token = token.UPDATE_TOKEN
	case "DELETE":
		tok.Token = token.DELETE_TOKEN
	case "SET":
		tok.Token = token.SET_TOKEN
	case "WHERE":
		tok.Token = token.WHERE_TOKEN
	default:
		tok.Token = token.IDENT_TOKEN
	}

	l.tokens = append(l.tokens, tok)
}

func (l *Lexer) Tokenize(input string) []token.Token {
	l.tokens = make([]token.Token, 0)
	l.ReadBuffer.Reset()

	inString := false
	var stringDelimiter rune

	for _, char := range input {
		if char == '\'' || char == '"' {
			if !inString {
				l.flushBuffer()
				inString = true
				stringDelimiter = char
				l.ReadBuffer.WriteRune(char)
			} else if char == stringDelimiter {
				l.ReadBuffer.WriteRune(char)
				tok := token.Token{
					Literal: l.ReadBuffer.String(),
					Token:   token.STRING_TOKEN,
				}
				l.tokens = append(l.tokens, tok)
				l.ReadBuffer.Reset()
				inString = false
			} else {
				l.ReadBuffer.WriteRune(char)
			}
			continue
		}

		if inString {
			l.ReadBuffer.WriteRune(char)
			continue
		}

		switch char {
		case ',':
			l.flushBuffer()
			l.tokens = append(l.tokens, token.Token{
				Literal: ",",
				Token:   token.COMMA_TOKEN,
			})
		case '(':
			l.flushBuffer()
			l.tokens = append(l.tokens, token.Token{
				Literal: "(",
				Token:   token.LPAREN_TOKEN,
			})
		case ')':
			l.flushBuffer()
			l.tokens = append(l.tokens, token.Token{
				Literal: ")",
				Token:   token.RPAREN_TOKEN,
			})
		case '=':
			l.flushBuffer()
			l.tokens = append(l.tokens, token.Token{
				Literal: "=",
				Token:   token.EQUALS_TOKEN,
			})
		case ';':
			l.flushBuffer()
			l.tokens = append(l.tokens, token.Token{
				Literal: ";",
				Token:   token.SEMICOLON_TOKEN,
			})
		case '\n':
			l.flushBuffer()
			l.tokens = append(l.tokens, token.Token{
				Literal: "\n",
				Token:   token.ENDLINE_TOKEN,
			})
		case ' ', '\t', '\r':
			l.flushBuffer()
		default:
			l.ReadBuffer.WriteRune(char)
		}
	}

	l.flushBuffer()

	return l.tokens
}

func (l *Lexer) GetTokens() []token.Token {
	return l.tokens
}

func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !unicode.IsDigit(c) && c != '.' && c != '-' {
			return false
		}
	}
	return true
}
