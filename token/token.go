package token

type TokenType string

const (
	SELECT_TOKEN  = "SELECT"
	FROM_TOKEN    = "FROM"
	IDENT_TOKEN   = "IDENT"
	ENDLINE_TOKEN = "ENDLINE"
	COMMA_TOKEN   = ","
)

type Token struct {
	Literal string
	Token   TokenType
}
