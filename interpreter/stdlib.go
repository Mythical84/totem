package interpreter

import (
	"fmt"
	. "main/env"
	. "main/errors"
	"strconv"
	"strings"
	"time"
)

func CreateGlobals() *Environment {
	globals := CreateEnv(nil, "")

	globals.Define("clock", Clock{name: "clock"}, false)
	globals.Define("input", Input{name: "input"}, false)
	globals.Define("int", Int{name: "int"}, false)
	globals.Define("append", Append{name: "append"}, false)
	globals.Define("chr", Chr{name: "chr"}, false)
	globals.Define("len", Len{name: "len"}, false)
	globals.Define("ord", Ord{name: "ord"}, false)
	globals.Define("print", Print{name: "print"}, false)
	globals.Define("println", Println{name: "println"}, false)

	return globals
}

func CreateBuiltin(name string) (*Environment, bool) {
	switch name {
	case "Math":
		return CreateMath(), true
	default:
		return nil, false
	}
}

type Base struct{ name string }

func (self Base) String() string {
	return "<builtin-function '" + self.name + "'>"
}

type Clock Base

func (self Clock) Call(inter interpreter, args []any, _ int) (any, error) {
	return float64(time.Now().Unix()), nil
}

func (self Clock) Arity() int {
	return 0
}

type Input Base

func (self Input) Call(inter interpreter, args []any, _ int) (any, error) {
	fmt.Print(args[0])
	var buf string
	_, err := fmt.Scanln(&buf)
	if err != nil {
		return StrType{str: ""}, nil 
	}
	return &StrType{str: buf}, nil
}

func (self Input) Arity() int {
	return 1
}

type Int Base

func (self Int) Call(inter interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return val, nil
	case string:
		i, _ := strconv.ParseFloat(val, 64)
		return i, nil
	case bool:
		if val {
			return 1, nil
		} else {
			return 0, nil
		}
	case nil:
		return 0, nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Int) Arity() int {
	return 1
}

type Append Base

func (self Append) Call(inter interpreter, args []any, line int) (any, error) {
	switch expr := args[0].(type) {
	case IterType:
		expr.Setval(append(expr.Getvals(), args[1]))
		return nil, nil
	default:
		return nil, RuntimeError("Object not subscriptable", line, filename)
	}
}

func (self Append) Arity() int {
	return 2
}

type Chr Base

func (self Chr) Call(inter interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return string(rune(val)), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Chr) Arity() int {
	return 1
}

type Len Base

func (self Len) Call(inter interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case IterType:
		return float64(len(val.Getvals())), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Len) Arity() int {
	return 1
}

type Ord Base

func (self Ord) Call(inter interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case StrType:
		if len(val.Getvals()) > 1 {
			return nil, RuntimeError("Input string to long", line, filename)
		}
		return float64([]rune(val.Getvals()[0].(string))[0]), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Ord) Arity() int {
	return 1
}

type Print Base

func (self Print) Call(inter interpreter, args []any, _ int) (any, error) {
	if len(args) == 0 {
		return nil, nil
	}
	var str strings.Builder
	for _, val := range args[:len(args)-1] {
		fmt.Fprintf(&str, "%v ", val)
	}
	fmt.Fprintf(&str, "%v", args[len(args)-1])
	fmt.Print(str.String())

	return nil, nil
}

func (self Print) Arity() int {
	return -1
}

type Println Base

func (self Println) Call(inter interpreter, args []any, _ int) (any, error) {
	if len(args) == 0 {
		fmt.Println("")
		return nil, nil
	}
	var str strings.Builder
	for _, val := range args[:len(args)-1] {
		fmt.Fprintf(&str, "%v ", val)
	}
	fmt.Fprintf(&str, "%v", args[len(args)-1])
	fmt.Println(str.String())

	return nil, nil
}

func (self Println) Arity() int {
	return -1
}

type F Base

func (self F) call(inter interpreter, args []any, _ int) (error, any) {
	return nil, nil
}

func (self F) Arity() int {
	return -1
}


