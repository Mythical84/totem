package lexer

import "fmt"

type Token struct {
	Value    any
	Line     int
	Type     TokenType
	Location int
}

func (token Token) String() string {
	return fmt.Sprintf("{ type: %s, value: %v, line: %d }",
		TokenName[token.Type],
		token.Value,
		token.Line)
}

type TokenType int

const (
	EOF TokenType = iota

	// Single Characters
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_SQUARE
	RIGHT_SQUARE
	COMMA
	DOT
	EOL
	COLON

	// Single or Double Characters
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	SUB
	SUB_EQUAL
	ADD
	ADD_EQUAL
	DIV
	DIV_EQUAL
	MOD
	MOD_EQUAL
	MUL
	MUL_EQUAL
	POWER

	// Literals
	STR
	NUM
	TRUE
	FALSE
	NULL
	IDENTIFIER

	// Keywords
	AND
	OR
	IF
	ELSE
	CLASS
	SELF
	RETURN
	CONTINUE
	BREAK
	LOOP
	WHILE
	FOR
	FN
	TRY
	CATCH
	LOCAL
	IMPORT
)

var TokenName = map[TokenType]string{
	EOF:           "EOF",
	LEFT_PAREN:    "LEFT_PAREN",
	RIGHT_PAREN:   "RIGHT_PAREN",
	LEFT_BRACE:    "LEFT_BRACE",
	RIGHT_BRACE:   "RIGHT_BRACE",
	LEFT_SQUARE:   "LEFT_SQUARE",
	RIGHT_SQUARE:  "RIGHT_SQUARE",
	COMMA:         "COMMA",
	DOT:           "DOT",
	SUB:           "SUB",
	SUB_EQUAL:     "SUB_EQUAL",
	ADD:           "ADD",
	ADD_EQUAL:     "ADD_EQUAL",
	DIV:           "DIV",
	DIV_EQUAL:     "DIV_EQUAL",
	MOD:           "MOD",
	MOD_EQUAL:     "MOD_EQUAL",
	MUL:           "MUL",
	POWER:         "POWER",
	MUL_EQUAL:     "MUL_EQUAL",
	EOL:           "EOL",
	BANG:          "BANG",
	BANG_EQUAL:    "BANG_EQUAL",
	EQUAL:         "EQUAL",
	EQUAL_EQUAL:   "EQUAL_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",
	STR:           "STR",
	NUM:           "NUM",
	NULL:          "NULL",
	IDENTIFIER:    "IDENTIFIER",
	AND:           "AND",
	OR:            "OR",
	IF:            "IF",
	ELSE:          "ELSE",
	CLASS:         "CLASS",
	SELF:          "SELF",
	RETURN:        "RETURN",
	CONTINUE:      "CONTINUE",
	BREAK:         "BREAK",
	LOOP:          "LOOP",
	WHILE:         "WHILE",
	FOR:           "FOR",
	FN:            "FN",
	TRUE:          "TRUE",
	FALSE:         "FALSE",
	TRY:           "TRY",
	CATCH:         "CATCH",
	LOCAL:         "LOCAL",
	IMPORT:        "IMPORT",
}

var token_type = map[rune]TokenType{
	'(': LEFT_PAREN,
	')': RIGHT_PAREN,
	'{': LEFT_BRACE,
	'}': RIGHT_BRACE,
	'[': LEFT_SQUARE,
	']': RIGHT_SQUARE,
	',': COMMA,
	'.': DOT,
	'-': SUB,
	'+': ADD,
	'*': MUL,
	'/': DIV,
	'%': MOD,
	';': EOL,
	'!': BANG,
	'=': EQUAL,
	'>': GREATER,
	'<': LESS,
	':': COLON,
}
