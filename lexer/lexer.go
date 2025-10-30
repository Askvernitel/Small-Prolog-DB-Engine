package lexer

import (
	"bytes"
	"strings"
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

	switch strings.ToUpper(literal) {
	case "SELECT":
		tok.Token = token.SELECT_TOKEN
	case "FROM":
		tok.Token = token.FROM_TOKEN
	default:
		tok.Token = token.IDENT_TOKEN
	}

	l.tokens = append(l.tokens, tok)
}

func (l *Lexer) Tokenize(input string) []token.Token {
	l.tokens = make([]token.Token, 0)
	l.ReadBuffer.Reset()

	for _, char := range input {
		switch char {
		case ',':
			// Flush any buffered content first
			l.flushBuffer()
			// Add comma token
			l.tokens = append(l.tokens, token.Token{
				Literal: ",",
				Token:   token.COMMA_TOKEN,
			})
		case '\n':
			// Flush buffer before newline
			l.flushBuffer()
			// Add endline token
			l.tokens = append(l.tokens, token.Token{
				Literal: "\n",
				Token:   token.ENDLINE_TOKEN,
			})
		case ' ', '\t', '\r':
			// Whitespace ends current token
			l.flushBuffer()
		default:
			// Accumulate character in buffer
			l.ReadBuffer.WriteRune(char)
		}
	}

	// Flush any remaining content
	l.flushBuffer()

	return l.tokens
}

func (l *Lexer) GetTokens() []token.Token {
	return l.tokens
}
