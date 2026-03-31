package parser

import (
	"fmt"
	. "main/errors"
	. "main/lexer"
	. "main/visitor"
	"slices"
)

var tokens []Token
var index int
var grouping int
var filename string
var text string

func Parse(token_list []Token, content string, file string) ([]Stmt, []error) {
	fmt.Print("")
	tokens = token_list
	filename = file
	text = content
	is_error := false
	stmts := []Stmt{}
	errors := []error{}
	index = 0
	grouping = 0

	for !at_end() {
		stmt, errs := declaration()
		if errs != nil {
			is_error = true
			errors = append(errors, errs...)
			synchronize()
		} else {
			stmts = append(stmts, stmt)
		}
	}

	if is_error {
		return nil, errors
	}

	return stmts, nil
}

func declaration() (Stmt, []error) {
	if match(FN) {
		return func_declaration()
	} else if match(CLASS) {
		return class_declaration()
	} else {
		return statement()
	}
}

func func_declaration() (Stmt, []error) {
	name, err := consume("Expected Identifier", peek().Location, IDENTIFIER)
	if err != nil {
		return nil, err
	}
	left_paren, err := consume("Expected '('", peek().Location, LEFT_PAREN)
	grouping++
	if err != nil {
		return nil, err
	}

	args := []string{}
	if !check(RIGHT_PAREN) {
		arg, err := consume("Expected argument name", peek().Location, IDENTIFIER)
		if err != nil {
			return nil, err
		}
		args = append(args, arg.Value.(string))
		for match(COMMA) {
			arg, err := consume("Expected argument name", peek().Location, IDENTIFIER)
			if err != nil {
				return nil, err
			}
			args = append(args, arg.Value.(string))
		}
	}

	_, err = consume("Unclosed '('", left_paren.Location, RIGHT_PAREN)
	body, err := statement()
	if err != nil {
		return nil, err
	}
	return Func{Name: name, Args: args, Body: body}, nil
}

func class_declaration() (Stmt, []error) {
	stmts := []Stmt{}
	errors := []error{}
	left := previous()

	name, err := consume("Expected identifier", peek().Location, IDENTIFIER)
	if err != nil {
		errors = append(errors, err...)
	}

	var superclass Expr = nil
	if match(COLON) {
		consume("Expected superclass name", peek().Location, IDENTIFIER)
		superclass = VarGet{Name: previous()}
	}
	_, err = consume("Expected '{'", peek().Location, LEFT_BRACE)
	if err != nil {
		errors = append(errors, err...)
	}

	match(EOL)

	for !check(RIGHT_BRACE) && !at_end() {
		stmt, err := declaration()
		if err != nil {
			errors = append(errors, err...)
			synchronize()
		}
		stmts = append(stmts, stmt)
	}

	_, err = consume("'{' was not closed", left.Location, RIGHT_BRACE)
	if err != nil {
		errors = append(errors, err...)
	}

	match(EOL)

	if len(errors) != 0 {
		return nil, errors
	} else {
		return ClassDecl{Name: name, Body: stmts, Superclass: superclass}, nil
	}
}

func statement() (Stmt, []error) {
	if match(LEFT_BRACE) {
		return block_statement()
	} else if match(TRY) {
		return try_catch_statement()
	} else if match(IF) {
		return if_statement()
	} else if match(WHILE) {
		return while_statement()
	} else if match(LOOP) {
		return loop_statement()
	} else if match(RETURN) {
		return return_statement()
	} else if match(BREAK) {
		return break_statement()
	} else if match(CONTINUE) {
		return continue_statement()
	} else if match(IMPORT) {
		return import_statement()
	} else {
		return expression_statement()
	}
}

func block_statement() (Stmt, []error) {
	stmts := []Stmt{}
	errors := []error{}
	left := previous()

	match(EOL)

	for !check(RIGHT_BRACE) && !at_end() {
		stmt, err := statement()
		if err != nil {
			errors = append(errors, err...)
			synchronize()
		}
		stmts = append(stmts, stmt)
	}

	_, err := consume("'{' was not closed", left.Location, RIGHT_BRACE)
	if err != nil {
		errors = append(errors, err...)
	}

	match(EOL)

	if len(errors) != 0 {
		return nil, errors
	} else {
		return Block{Body: stmts}, nil
	}
}

func try_catch_statement() (Stmt, []error) {
	try := previous()
	try_body, err := statement()
	if err != nil {
		return nil, err
	}
	_, err = consume("Incomplete Try Catch statement", try.Location, CATCH)
	if err != nil {
		return nil, err
	}

	identifier, _ := consume("Expected identifier", peek().Location, IDENTIFIER)

	catch_body, err := statement()
	if err != nil {
		return nil, err
	}

	return TryCatch{TryBody: try_body, CatchBody: catch_body,
		ErrName: identifier}, nil
}

func if_statement() (Stmt, []error) {
	expr, err := expression()
	if err != nil {
		return nil, err
	}
	body, err := statement()
	if err != nil {
		return nil, err
	}

	var else_body Stmt = nil
	if match(ELSE) {
		else_body, err = statement()
		if err != nil {
			return nil, err
		}
	}

	return If{IfBody: body, ElseBody: else_body, Condition: expr}, nil
}

func while_statement() (Stmt, []error) {
	expr, err := expression()
	if err != nil {
		return nil, err
	}

	body, err := statement()
	if err != nil {
		return nil, err
	}

	return While{Body: body, Condition: expr}, nil
}

func loop_statement() (Stmt, []error) {
	body, err := statement()
	if err != nil {
		return nil, err
	}

	return Loop{Body: body}, nil
}

func return_statement() (Stmt, []error) {
	if tokens[index-1].Type == EOL {
		return Return{Expr: Literal{Value: nil}, Line: previous().Line}, nil
	}
	expr, err := expression()
	if err != nil {
		return nil, err
	}

	_, err = consume("Invalid Syntax", peek().Location, EOL)

	return Return{Expr: expr, Line: previous().Line}, nil
}

func break_statement() (Stmt, []error) {
	_, err := consume("Invalid Syntax", peek().Location, EOL)
	if err != nil {
		return nil, err
	}
	return Break{Line: previous().Line}, nil
}

func continue_statement() (Stmt, []error) {
	_, err := consume("Invalid Syntax", peek().Location, EOL)
	if err != nil {
		return nil, err
	}
	return Continue{Line: previous().Line}, nil
}

func import_statement() (Stmt, []error) {
	path, err := consume("Expected import path", peek().Location, STR)
	if err != nil {
		return nil, err
	}
	_, err = consume("Invalid syntax", previous().Location, EOL)
	if err != nil {
		return nil, err
	}
	return Import{Path: path}, nil
}

func expression_statement() (Stmt, []error) {
	expr, err := expression()
	if err != nil {
		return nil, err
	}
	_, err = consume("Invalid Syntax", peek().Location, EOL)
	if err != nil {
		return nil, err
	}

	return Expression{Expr: expr}, nil
}

func expression() (Expr, []error) {
	return assignment()
}

func assignment() (Expr, []error) {
	local := check(LOCAL)
	if local {
		advance()
	}
	expr, err := compound_assignment()
	if err != nil {
		return nil, err
	}
	if match(EQUAL) {
		equals := previous()
		val, err := assignment()
		if err != nil {
			return nil, err
		}
		switch expr := expr.(type) {
		case VarGet:
			return VarDef{Name: expr.Name, Expr: val, Local: local}, nil
		case Get:
			return Set{Expr: expr.Expr, Name: expr.Name, Val: val}, nil
		case IterGet:
			return IterSet{Expr: expr.Expr, Index: expr.Index, 
				Line: expr.Line, Val: val}, nil
		default:
			return nil, []error{SyntaxError("Invalid assignment target", equals.Line,
				equals.Location, text, filename)}
		}
	} else if local {
		return nil, []error{SyntaxError("Expected variable name", peek().Line,
			peek().Location, text, filename)}
	}

	return expr, nil
}

func compound_assignment() (Expr, []error) {
	expr, err := tuple()
	if err != nil {
		return nil, err
	}

	if match(ADD_EQUAL, SUB_EQUAL, MUL_EQUAL, DIV_EQUAL, MOD_EQUAL) {
		equals := previous()
		equals.Type -= 1
		val, err := compound_assignment()
		if err != nil {
			return nil, err
		}
		switch expr := expr.(type) {
		case VarGet:
			return VarDef{
				Name: expr.Name,
				Expr: Binary{
					Left:     expr,
					Operator: equals,
					Right:    val,
				},
			}, nil
		case Get:
			return Set{
				Name: expr.Name,
				Expr: expr.Expr,
				Val: Binary{
					Left:     expr,
					Operator: equals,
					Right:    val,
				},
			}, nil
		case IterGet:
			return IterSet{
				Expr: expr.Expr,
				Index: expr.Index,
				Line: expr.Line,
				Val: Binary{
					Left: expr,
					Operator: equals,
					Right: val,
				},
			}, nil
		default:
			return nil, []error{SyntaxError("Invalid assignment target", equals.Line,
				equals.Location, text, filename)}
		}
	}

	return expr, nil
}

func tuple() (Expr, []error) {
	expr, err := logic_or()
	if err != nil {
		return nil, err
	}

	if check(COMMA) {
		exprs := []Expr{expr}
		for match(COMMA) {
			val, err := logic_or()
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, val)
		}

		expr = Tuple{Vals: exprs}
	}

	return expr, nil
}

func logic_or() (Expr, []error) {
	expr, err := logic_and()
	if err != nil {
		return nil, err
	}

	for match(OR) {
		operator := previous()
		right, err := logic_and()
		if err != nil {
			return nil, err
		}
		expr = Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func logic_and() (Expr, []error) {
	expr, err := equality()
	if err != nil {
		return nil, err
	}

	for match(AND) {
		operator := previous()
		right, err := equality()
		if err != nil {
			return nil, err
		}
		expr = Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func equality() (Expr, []error) {
	var expr, err = comparison()
	if err != nil {
		return nil, err
	}

	for match(EQUAL_EQUAL, BANG_EQUAL) {
		var operator = previous()
		var right, err = comparison()
		if err != nil {
			return nil, err
		}
		expr = Truth{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func comparison() (Expr, []error) {
	var expr, err = term()
	if err != nil {
		return nil, err
	}

	for match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		var operator = previous()
		var right, err = term()
		if err != nil {
			return nil, err
		}
		expr = Truth{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func term() (Expr, []error) {
	var expr, err = factor()
	if err != nil {
		return nil, err
	}

	for match(ADD, SUB) {
		var operator = previous()
		var right, err = factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func factor() (Expr, []error) {
	var expr, err = power()
	if err != nil {
		return nil, err
	}

	for match(MUL, DIV, MOD) {
		var operator = previous()
		var right, err = power()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func power() (Expr, []error) {
	expr, err := unary()
	if err != nil {
		return nil, err
	}

	for match(POWER) {
		operator := previous()
		right, err := unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func unary() (Expr, []error) {
	if match(BANG, SUB, ADD) {
		var operator = previous()
		var expr, err = unary()
		if err != nil {
			return nil, err
		}
		return Unary{Operator: operator, Expr: expr}, nil
	}

	return list()
}

func list() (Expr, []error) {
	expr, err := call()
	if err != nil {
		return nil, err
	}

	for true {
		if match(LEFT_SQUARE) {
			left_square := previous()
			grouping++
			val, err := logic_or()
			if err != nil {
				return nil, err
			}
			consume("Unclosed [", left_square.Location, RIGHT_SQUARE)
			expr = IterGet{Expr: expr, Index: val, Line: previous().Line}
		} else {
			break
		}
	}

	return expr, nil
}

func call() (Expr, []error) {
	expr, err := atom()
	if err != nil {
		return nil, err
	}

	for true {
		if match(LEFT_PAREN) {
			left_paren := previous()
			grouping++
			args := []Expr{}
			if !check(RIGHT_PAREN) {
				arg, err := logic_or()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
				for match(COMMA) {
					arg, err := logic_or()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)
				}
			}
			var right_paren Token
			if check(RIGHT_PAREN) {
				index++
				right_paren = previous()
			} else {
				return nil, []error{SyntaxError("Unclosed '('",
					left_paren.Line, left_paren.Location, text, filename)}
			}
			if err != nil {
				return nil, err
			}
			grouping--
			expr = Call{Callee: expr, Args: args, Paren: right_paren}
		} else if match(DOT) {
			name, err := consume("Expected property name",
				peek().Location, IDENTIFIER)
			if err != nil {
				return nil, err
			}
			expr = Get{Expr: expr, Name: name}
		} else {
			break
		}
	}

	return expr, nil
}

func atom() (Expr, []error) {
	if match(FALSE) {
		return Literal{Value: false}, nil
	} else if match(TRUE) {
		return Literal{Value: true}, nil
	} else if match(NULL) {
		return Literal{Value: nil}, nil
	} else if match(LEFT_SQUARE) {
		left_square := previous()
		exprs := []Expr{}
		if !check(RIGHT_SQUARE) {
			grouping++
			expr, err := logic_or()
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, expr)
			if check(COMMA) {
				for match(COMMA) {
					val, err := logic_or()
					if err != nil {
						return nil, err
					}
					exprs = append(exprs, val)
				}

			}
		}
		_, err := consume("Unclosed [", left_square.Location, RIGHT_SQUARE)
		if err != nil {
			return nil, err
		}

		grouping--
		return ListDef{Vals: exprs}, nil
	} else if match(NUM, STR) {
		return Literal{Value: previous().Value}, nil
	} else if match(IDENTIFIER) {
		return VarGet{Name: previous()}, nil
	} else if match(SELF) {
		return Self{Line: previous().Line}, nil
	} else if match(LEFT_PAREN) {
		grouping++
		var left_paren = previous()
		var expr, err = expression()
		if err != nil {
			return nil, err
		}
		consume("'(' was not closed", left_paren.Location, RIGHT_PAREN)
		grouping--
		return Grouping{Expr: expr}, nil
	}

	return nil, []error{SyntaxError("Expected Expression", peek().Line,
		peek().Location, text, filename)}
}

func match(types ...TokenType) bool {
	if slices.Contains(types, peek().Type) {
		advance()
		return true
	}
	return false
}

func advance() Token {
	if !at_end() {
		index++
	}
	if (grouping > 0 && peek().Type == EOL ||
		previous().Type == EOL && peek().Type == EOL) &&
		!at_end() {
		return advance()
	}
	return previous()
}

func at_end() bool {
	temp := index
	for tokens[temp].Type != EOF {
		temp++
		if tokens[temp-1].Type != EOL {
			return false
		}
	}
	return true
}

func peek() Token {
	return tokens[index]
}

func previous() Token {
	if grouping > 0 {
		temp := -1
		for (temp+index) > 0 && tokens[index+temp].Type == EOL {
			temp -= 1
		}
		return tokens[index+temp]
	} else {
		return tokens[index-1]
	}
}

func consume(err string, location int, token_type TokenType) (Token, []error) {
	if match(token_type) {
		return previous(), nil
	}

	// Kinda a bandaid solution. Removing semi colons is a pain in the ass
	if token_type == EOL && tokens[index-1].Type == EOL {
		return previous(), nil
	}

	// I can't return nil since this function returns a struct,
	// but the results of this function are discared if an error is thrown
	return previous(), []error{SyntaxError(err, previous().Line,
		location, text, filename)}
}

func check(token_type TokenType) bool {
	return peek().Type == token_type
}

func synchronize() {
	grouping = 0
	advance()

	for !at_end() {
		if previous().Type == EOL {
			return
		}

		switch peek().Type {
		case CLASS, FN, FOR, WHILE, IF, RETURN:
			return
		}

		advance()
	}
}
