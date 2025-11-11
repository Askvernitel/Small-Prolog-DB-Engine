package token

type TokenType string

const (
	SELECT_TOKEN = "SELECT"
	FROM_TOKEN   = "FROM"
	INSERT_TOKEN = "INSERT"
	INTO_TOKEN   = "INTO"
	VALUES_TOKEN = "VALUES"
	UPDATE_TOKEN = "UPDATE"
	DELETE_TOKEN = "DELETE"
	SET_TOKEN    = "SET"
	WHERE_TOKEN  = "WHERE"

	IDENT_TOKEN  = "IDENT"
	STRING_TOKEN = "STRING"
	NUMBER_TOKEN = "NUMBER"

	COMMA_TOKEN     = ","
	LPAREN_TOKEN    = "("
	RPAREN_TOKEN    = ")"
	EQUALS_TOKEN    = "="
	ENDLINE_TOKEN   = "ENDLINE"
	SEMICOLON_TOKEN = ";"
)

type Token struct {
	Literal string
	Token   TokenType
}
