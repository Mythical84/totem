package interpreter

import (
	"fmt"
	. "main/env"
	. "main/errors"
	. "main/lexer"
	. "main/visitor"
	"math"
	"os"
	"reflect"
	"slices"
)

var filename string
var env *Environment

func Interpret(stmts []Stmt, file string) {
	filename = file
	env = CreateEnv(CreateGlobals())
	inter := interpreter{}
	for _, stmt := range stmts {
		err := inter.execute(stmt)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func (self interpreter) execute(stmt Stmt) error {
	return stmt.Accept(self)
}

func (self interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(self)
}

// Declarations
func (self interpreter) VisitFunc(func_stmt Func) error {
	env.Define(
		func_stmt.Name.Value.(string),
		Function{
			name: func_stmt.Name.Value.(string),
			args: func_stmt.Args,
			body: func_stmt.Body,
		},
		false,
	)
	return nil
}

// TODO: Change so that inheritance doesn't copy the entire env
// instead, set the parent class env as the parent to the subclass env
// TODO: Add proper support for super instead of simply calling
// the parent init function
func (self interpreter) VisitClassDecl(class_decl ClassDecl) error {
	prev := env
	temp := CreateEnv(prev)
	env = temp

	if class_decl.Superclass != nil {
		val, err := self.evaluate(class_decl.Superclass)
		if err != nil {
			return err
		}
		switch val := val.(type) {
		case Class:
			for key := range val.env.Keys() {
				v, _ := val.env.Get(key, 0, "")
				if key == "init" {
					key = "super"
				}
				env.Define(key, v, true)
			}
		default:
			return RuntimeError("Superclass must be a class",
				class_decl.Name.Line, filename)
		}
	}
	env.Define("init", Function{name: "init", args: []string{},
		body: Block{Body: []Stmt{}}, instance: nil}, true)

	for _, stmt := range class_decl.Body {
		err := self.execute(stmt)
		if err != nil {
			env = prev
			return err
		}
	}

	arity := 0
	val, _ := env.Get("init", 0, "")
	switch val := val.(type) {
	case Function:
		arity = val.Arity()
	}

	class := Class{
		name:  class_decl.Name.Value.(string),
		env:   env,
		arity: arity,
	}
	env = prev

	env.Define(class_decl.Name.Value.(string), class, false)
	return nil
}

// Statements
func (self interpreter) VisitExpression(expression Expression) error {
	_, err := self.evaluate(expression.Expr)
	return err
}

func (self interpreter) VisitBlock(block Block) error {
	prev := env
	env = CreateEnv(env)

	for _, stmt := range block.Body {
		err := self.execute(stmt)
		if err != nil {
			env = prev
			return err
		}
	}

	env = prev

	return nil
}

func (self interpreter) VisitIf(if_stmt If) error {
	condition, err := self.evaluate(if_stmt.Condition)
	if err != nil {
		return err
	}
	if isTruthy(condition) {
		err = self.execute(if_stmt.IfBody)
		if err != nil {
			return err
		}
	} else if if_stmt.ElseBody != nil {
		err = self.execute(if_stmt.ElseBody)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self interpreter) VisitWhile(while_stmt While) error {
	condition, err := self.evaluate(while_stmt.Condition)
	if err != nil {
		return err
	}
	for isTruthy(condition) {
		err = self.execute(while_stmt.Body)
		switch err := err.(type) {
		case BreakErrorType:
			break
		case ContinueErrorType:
			continue
		case nil:
		default:
			return err
		}
		condition, err = self.evaluate(while_stmt.Condition)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self interpreter) VisitLoop(loop Loop) error {
	loop: for true {
		err := self.execute(loop.Body)
		switch err := err.(type) {
		case BreakErrorType:
			break loop
		case ContinueErrorType:
			continue loop
		case nil:
		default:
			return err
		}
	}
	return nil
}

func (self interpreter) VisitReturn(return_stmt Return) error {
	val, err := self.evaluate(return_stmt.Expr)
	if err != nil {
		return err
	}
	return ReturnErrorType{Value: val, Line: return_stmt.Line, Filename: filename}
}

func (self interpreter) VisitBreak(break_stmt Break) error {
	return BreakErrorType{Line: break_stmt.Line, Filename: filename}
}

func (self interpreter) VisitContinue(continue_stmt Continue) error {
	return ContinueErrorType{Line: continue_stmt.Line, Filename: filename}
}

// TODO: Implement once I have classes
func (self interpreter) VisitTryCatch(try_catch TryCatch) error {
	err := self.execute(try_catch.TryBody)
	if err != nil {
		switch err := err.(type) {
		case ReturnErrorType, ContinueErrorType, BreakErrorType:
			return err
		default:
			self.execute(try_catch.CatchBody)
			return nil
		}
	} else {
		return nil
	}
}

// Expressions
func (self interpreter) VisitBinary(binary Binary) (any, error) {
	left, err := self.evaluate(binary.Left)
	if err != nil {
		return nil, err
	}
	right, err := self.evaluate(binary.Right)
	if err != nil {
		return nil, err
	}

	switch binary.Operator.Type {
	case SUB, MOD, DIV, POWER:
		if typeName(left) != "float64" || typeName(right) != "float64" {
			RuntimeError("Token must be a number", binary.Operator.Line, filename)
		}
	}

	switch binary.Operator.Type {
	case SUB:
		return left.(float64) - right.(float64), nil
	case POWER:
		return math.Pow(left.(float64), right.(float64)), nil
	case DIV:
		if right.(float64) == 0 {
			RuntimeError("Division by 0", binary.Operator.Line, filename)
		}
		return left.(float64) / right.(float64), nil
	case MOD:
		if right.(float64) == 0 {
			RuntimeError("Division by 0", binary.Operator.Line, filename)
		}
		return math.Mod(left.(float64), right.(float64)), nil
	case MUL:
		switch left := left.(type) {
		case IterType:
			if right, ok := right.(float64); ok {
				temp := left
				temp.Setval(slices.Repeat(left.Getvals(), int(right)))
				return temp, nil
			}
		case float64:
		if right, ok := right.(IterType); ok {
				temp := right
				temp.Setval(slices.Repeat(right.Getvals(), int(left)))
				return temp, nil
			}
		}
		return nil, RuntimeError("Unsupported operand types for *: "+
			typeName(left)+" and "+typeName(right), binary.Operator.Line, filename)
	case ADD:
		if typeName(left) == "float64" && typeName(right) == "float64" {
			return left.(float64) + right.(float64), nil
		}
		return fmt.Sprintf("%v%v", left, right), nil
	}

	panic("Unknown binary operator")
}

func (self interpreter) VisitTruth(truth Truth) (any, error) {
	left, err := self.evaluate(truth.Left)
	if err != nil {
		return nil, err
	}
	right, err := self.evaluate(truth.Right)
	if err != nil {
		return nil, err
	}

	switch truth.Operator.Type {
	case GREATER_EQUAL, GREATER, LESS_EQUAL, LESS:
		if typeName(left) != "float64" || typeName(right) != "float64" {
			RuntimeError("Token must be a number", truth.Operator.Line, filename)
		}
	}

	switch truth.Operator.Type {
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64), nil
	case LESS_EQUAL:
		return left.(float64) <= right.(float64), nil
	case LESS:
		return left.(float64) < right.(float64), nil
	case GREATER:
		return left.(float64) > right.(float64), nil
	}

	panic("Unknown truth operator")
}

func (self interpreter) VisitLogical(logical Logical) (any, error) {
	left, err := self.evaluate(logical.Left)
	if err != nil {
		return nil, err
	}

	switch logical.Operator.Type {
	case OR:
		if isTruthy(left) {
			return left, nil
		}
	case AND:
		if !isTruthy(left) {
			return left, nil
		}
	}
	return self.evaluate(logical.Right)
}

func (self interpreter) VisitUnary(unary Unary) (any, error) {
	value, err := self.evaluate(unary.Expr)
	if err != nil {
		return nil, err
	}
	switch unary.Operator.Type {
	case ADD:
		if typeName(value) != "float64" {
			RuntimeError("Unsupported operand types for unary + "+
				fmt.Sprintf("%v", value), unary.Operator.Line, filename)
		}
		return value, nil
	case SUB:
		if typeName(value) != "float64" {
			RuntimeError("Unsupported operand types for unary - "+
				fmt.Sprintf("%v", value), unary.Operator.Line, filename)
		}
		return -value.(float64), nil
	case BANG:
		return !isTruthy(value), nil
	}

	panic("Unkown unary operator")
}

func (self interpreter) VisitTuple(tuple Tuple) (any, error) {
	vals := []any{}
	for _, val := range tuple.Vals {
		val, err := self.evaluate(val)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return &TupleType{vals: vals}, nil
}

func (self interpreter) VisitTupleDef(tuple_def TupleDef) (any, error) {
	for i, str := range tuple_def.Vars {
		val, err := self.evaluate(tuple_def.Vals[i])
		if err != nil {
			return nil, err
		}
		env.Define(str, val, tuple_def.Local)
	}

	return tuple_def.Vals, nil
}

func (self interpreter) VisitListDef(list_def ListDef) (any, error) {
	vals := []any{}
	for _, val := range list_def.Vals {
		val, err := self.evaluate(val)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return &ListType{vals: vals}, nil
}

func (self interpreter) VisitIterGet(iter_get IterGet) (any, error) {
	val, err := self.evaluate(iter_get.Expr)
	if err != nil {
		return nil, err
	}

	index, err := self.evaluate(iter_get.Index)
	if err != nil {
		return nil, err
	}

	if typeName(index) != "float64" {
		return nil, RuntimeError("Invalid index type", iter_get.Line, filename)
	}

	switch val := val.(type) {
	case IterType:
		if index.(float64) < 0 || int(index.(float64)) > len(val.Getvals()) {
			return nil, RuntimeError("Index out of bounds", iter_get.Line, filename)
		} else {
			v := val.Getvals()[int(index.(float64))]
			switch v := v.(type) {
			case string:
				return &StrType{str: v}, nil
			default:
				return v, nil
			}
		}
	default:
		return nil, RuntimeError("Object is not subscriptable",
			iter_get.Line, filename)
	}
}

func (self interpreter) VisitIterSet(iter_set IterSet) (any, error) {
	expr, err := self.evaluate(iter_set.Expr)
	if err != nil {
		return nil, err
	}
	index, err := self.evaluate(iter_set.Index)

	if typeName(index) != "float64" {
		return nil, RuntimeError("Invalid index type", iter_set.Line, filename)
	}

	val, err := self.evaluate(iter_set.Val)
	if err != nil {
		return nil, err
	}

	switch expr := expr.(type) {
	case IterType:
		e := expr.Getvals()
		e[int(index.(float64))] = val
		expr.Setval(e)
		return val, nil
	default:
		return nil, RuntimeError("Object not subscriptable", 
			iter_set.Line, filename)
	
	}
}

func (self interpreter) VisitCall(call Call) (any, error) {
	callee, err := self.evaluate(call.Callee)
	if err != nil {
		return nil, err
	}
	switch callee := callee.(type) {
	case Callable:
		args := []any{}
		for _, arg := range call.Args {
			val, err := self.evaluate(arg)
			if err != nil {
				return nil, err
			}
			args = append(args, val)
		}
		if callee.Arity() != -1 && len(args) != callee.Arity() {
			return nil, RuntimeError(fmt.Sprintf("Expected %d arguments, but got %d",
				callee.Arity(), len(args)), call.Paren.Line, filename)
		}
		return callee.Call(self, args, call.Paren.Line)
	default:
		return nil, RuntimeError(fmt.Sprintf("'%T' object not callable", callee),
			call.Paren.Line, filename)
	}
}

func (self interpreter) VisitGet(get Get) (any, error) {
	instance, err := self.evaluate(get.Expr)
	if err != nil {
		return nil, err
	}

	switch instance := instance.(type) {
	case ClassInstance:
		return instance.env.Get(get.Name.Value.(string), get.Name.Line, filename)
	default:
		return nil, RuntimeError("Only instances have properties",
			get.Name.Line, filename)
	}
}

func (self interpreter) VisitSet(set Set) (any, error) {
	instance, err := self.evaluate(set.Expr)
	if err != nil {
		return nil, err
	}

	switch instance := instance.(type) {
	case ClassInstance:
		val, err := self.evaluate(set.Val)
		if err != nil {
			return nil, err
		}
		instance.env.Define(set.Name.Value.(string), val, true)
		return val, nil
	default:
		return nil, RuntimeError("Only instances have properties",
			set.Name.Line, filename)
	}
}

func (self interpreter) VisitVarGet(var_get VarGet) (any, error) {
	return env.Get(var_get.Name.Value.(string), var_get.Name.Line, filename)
}

func (self interpreter) VisitVarDef(var_def VarDef) (any, error) {
	val, err := self.evaluate(var_def.Expr)
	if err != nil {
		return nil, err
	}
	env.Define(var_def.Name.Value.(string), val, var_def.Local)
	return val, nil
}

func (self interpreter) VisitSelf(self_stmt Self) (any, error) {
	val, err := env.Get("self", self_stmt.Line, filename)
	if err != nil {
		return nil, err
	}
	switch val := val.(type) {
	case *ClassInstance:
		return *val, nil
	default:
		return val, nil
	}
}

func (self interpreter) VisitGrouping(grouping Grouping) (any, error) {
	return self.evaluate(grouping.Expr)
}

func (self interpreter) VisitLiteral(literal Literal) (any, error) {
	switch v := literal.Value.(type) {
	case string:
		return &StrType{str: v}, nil
	default:
		return literal.Value, nil
	}
}

func isTruthy(value any) bool {
	switch value {
	case nil, false, 0:
		return false
	default:
		return true
	}
}

func typeName(value any) string {
	return fmt.Sprintf("%T", value)
}

func isEqual(left any, right any) bool {
	if typeName(left) != typeName(right) {
		return false
	}
	switch left := left.(type) {
	case IterType:
		switch right := right.(type) {
		case IterType:
			return reflect.DeepEqual(left.Getvals(), right.Getvals())
		}
	}
	return left == right
}
