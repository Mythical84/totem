package lexer

import (
	"main/errors"
	"strconv"
	"unicode"
)

var tokens = []Token{}
var start_line = 1
var line = 1
var start = 0
var current = 0
var text string
var file string

func Lexer(input string, filename string) ([]Token, []error) {
	error_list := []error{}
	text = input
	file = filename
	for current < len(text) {
		start = current
		start_line = line
		char := advance()
		switch char {
		//skip whitespace
		case ' ', '\t':

		case '\n':
			add_token_value(EOL, "")
			line++

		case '(', '[', '{', '}', ')', ']', ',', ';', ':':
			add_token(token_type[char])

		case '*':
			if match('*') {
				add_token(POWER)
			} else if match('=') {
				add_token(ADD_EQUAL)
			} else {
				add_token(MUL)
			}

		case '-', '+', '/', '%', '!', '=', '>', '<':
			if match('=') {
				add_token(token_type[char] + 1)
			} else {
				add_token(token_type[char])
			}

		case '#':
			if match('=') {
				err := multiline_comment()
				if err != nil {
					error_list = append(error_list, err)
				}
			} else {
				for peek() != '\n' && !is_at_end() {
					advance()
				}
			}

		case '"', '\'':
			err := add_string_token(char)
			if err != nil {
				error_list = append(error_list, err)
			}

		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			add_num_token(false)

		case '.':
			if unicode.IsDigit(char) {
				add_num_token(true)
			} else {
				add_token(DOT)
			}

		default:
			if unicode.IsLetter(char) || char == '_' {
				add_identifier_token()
			} else {
				error_list = append(error_list, errors.SyntaxError("Unkown Token",
					line, current, text, file))
			}

		}
	}

	add_token_value(EOF, "")

	if len(error_list) != 0 {
		return nil, error_list
	} else {
		return tokens, nil
	}
}

func add_string_token(quote rune) error {
	// skip opening quotation mark
	for peek() != quote {
		if is_at_end() {
			return errors.SyntaxError("Expected token "+string(quote),
				start_line, start, text, file)
		}
		if peek() == '\n' {
			line++
		}
		advance()
	}
	// skip closing quotation mark
	advance()

	add_token_value(STR, text[start+1:current-1])
	return nil
}

func multiline_comment() error {
	for text[current-1:current+1] != "=#" {
		advance()
		if peek() == '\n' {
			line++
		}
		if current+1 >= len(text) {
			return errors.SyntaxError("Unterminated block comment",
				start_line, start, text, file)
		}
	}
	return nil
}

func add_num_token(leading_decimal bool) error {
	decimal := leading_decimal
	for (unicode.IsDigit(peek()) || (peek() == '.' && !decimal) ||
		peek() == '_') && !is_at_end() {
		if peek() == '.' {
			decimal = true
		}
		advance()
	}

	val, _ := strconv.ParseFloat(text[start:current], 10)
	add_token_value(NUM, val)

	return nil
}

func add_identifier_token() {
	for unicode.IsLetter(peek()) || peek() == '_' || unicode.IsDigit(peek()) {
		advance()
	}
	k := get_keyword_token_type(text[start:current])
	if k != EOF {
		add_token(k)
		return
	}
	add_token(IDENTIFIER)
}

func get_keyword_token_type(keyword string) TokenType {
	switch keyword {
	case "and":
		return AND
	case "or":
		return OR
	case "if":
		return IF
	case "else":
		return ELSE
	case "class":
		return CLASS
	case "self":
		return SELF
	case "return":
		return RETURN
	case "continue":
		return CONTINUE
	case "break":
		return BREAK
	case "loop":
		return LOOP
	case "while":
		return WHILE
	case "for":
		return FOR
	case "fn":
		return FN
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "null":
		return NULL
	case "try":
		return TRY
	case "catch":
		return CATCH
	case "local":
		return LOCAL
	default:
		return EOF
	}

}

func add_token(token_type TokenType) {
	tokens = append(tokens, Token{
		Type:     token_type,
		Value:    text[start:current],
		Line:     line,
		Location: start,
	})
}

func add_token_value(token_type TokenType, val any) {
	tokens = append(tokens, Token{
		Type:     token_type,
		Value:    val,
		Line:     line,
		Location: start,
	})
}

func advance() rune {
	char := peek()
	current++
	return char
}

func peek() rune {
	if current >= len(text) {
		return '\000'
	}
	return []rune(text)[current]
}

func match(expected rune) bool {
	if is_at_end() {
		return false
	} else if peek() != expected {
		return false
	}

	advance()
	return true
}

func is_at_end() bool {
	return current == len(text)+1
}
