package env

import (
	"fmt"
	"iter"
	. "main/errors"
	"maps"
)

type Environment struct {
	values map[string]any
	Parent *Environment
}

func CreateEnv(parent *Environment) *Environment {
	return &Environment{
		Parent: parent,
		values: map[string]any{},
	}
}

func (self Environment) Define(name string, value any, local bool) {
	if !local && self.Parent != nil {
		if _, err := self.Parent.Get(name, 0, ""); err == nil {
			self.Parent.Define(name, value, local)
		} else {
			self.values[name] = value
		}
	} else {
		self.values[name] = value
	}
}

func (self Environment) Get(name string, line int, file string) (any, error) {
	val, ok := self.values[name]
	if ok {
		return val, nil
	}

	if self.Parent != nil {
		return self.Parent.Get(name, line, file)
	}

	return nil, RuntimeError("Undefined Variable: "+name, line, file)
}

func (self Environment) Keys() iter.Seq[string] {
	return maps.Keys(self.values)
}

func (self Environment) PrintMap() {
	fmt.Printf("%v\n", self.values)
}
