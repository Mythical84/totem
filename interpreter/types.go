package interpreter

import (
	"fmt"
	. "main/env"
	. "main/errors"
	. "main/visitor"
	"reflect"
	"strings"
)

type interpreter struct{}

type Callable interface {
	Call(inter interpreter, args []any, line int) (any, error)
	Arity() int
}

type Function struct {
	name     string
	args     []string
	body     Stmt
	instance *ClassInstance
}

func (self Function) Call(inter interpreter, args []any, _ int) (any, error) {
	prev := env
	env = CreateEnv(prev)
	for i := range args {
		env.Define(self.args[i], args[i], true)
	}
	if self.instance != nil {
		env.Define("self", self.instance, true)
	}

	err := inter.execute(self.body)
	env = prev
	switch err := err.(type) {
	case ReturnErrorType:
		if self.instance != nil && err.Value != nil {
			return nil, RuntimeError(fmt.Sprintf("init() should return null, not %T",
				err.Value), err.Line, filename)
		}
		return err.Value, nil
	default:
		return nil, err
	}
}

func (self Function) Arity() int {
	return len(self.args)
}

func (self Function) String() string {
	if self.instance == nil {
		return "<function '" + self.name + "'>"
	} else {
		return "<method '" + self.name + "' of " + self.instance.String() + ">"
	}
}

type Class struct {
	name  string
	env   *Environment
	arity int
}

func (self Class) Call(inter interpreter, args []any, line int) (any, error) {
	prev := env
	env = self.env

	instance := ClassInstance{name: self.name, env: env}
	for key := range env.Keys() {
		variable, _ := env.Get(key, 0, "")
		switch variable := variable.(type) {
		case Function:
			variable.instance = &instance
			env.Define(key, variable, true)
		}
	}

	init, _ := env.Get("init", 0, "")
	switch init := init.(type) {
	case Function:
		_, err := init.Call(inter, args, line)
		if err != nil {
			return nil, err
		}
		env = prev
		return instance, nil
	default:
		env = prev
		return nil, RuntimeError(fmt.Sprintf("'%T' object is not callable", init),
			line, filename)
	}
}

func (self Class) Arity() int {
	return self.arity
}

func (self Class) String() string {
	return "<class '" + self.name + "'>"
}

type ClassInstance struct {
	name string
	env  *Environment
}

func (self ClassInstance) String() string {
	return "<'" + self.name + "' object>"
}

type IterType interface {
	Getvals() []any
	Setval(val any)
}

type TupleType struct {
	vals []any
}

func (self TupleType) Getvals() []any {
	return self.vals
}

func (self *TupleType) Setval(val any) {}

func (self TupleType) String() string {
	var str strings.Builder
	str.WriteString("(")
	for i := range self.vals {
		fmt.Fprint(&str, self.vals[i])
		if i != len(self.vals)-1 {
			str.WriteString(", ")
		}
	}
	str.WriteString(")")
	return str.String()
}

type ListType struct {
	vals []any
}

func (self ListType) Getvals() []any {
	return self.vals
}

func (self *ListType) Setval(val any) {
	self.vals = val.([]any)
}

func (self ListType) String() string {
	var str strings.Builder
	str.WriteString("[")
	for i := range self.vals {
		fmt.Fprint(&str, self.vals[i])
		if i != len(self.vals)-1 {
			str.WriteString(", ")
		}
	}
	str.WriteString("]")
	return str.String()
}

type StrType struct {
	str string
}

func (self StrType) Getvals() []any {
	v := reflect.ValueOf(self.str)
	r := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		r[i] = string([]rune(self.str)[i])
	}
	return r
}

func (self *StrType) Setval(val any) {
	var str strings.Builder
	for _, v := range val.([]any) {
		fmt.Fprintf(&str, "%v", v)
	}
	self.str = str.String()
}

func (self StrType) String() string {
	return self.str
}
