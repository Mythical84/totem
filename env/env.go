package env

import (
	"fmt"
	"iter"
	. "main/errors"
	"maps"
)

type Environment struct {
	local map[string] any
	values map[string]any
	Parent *Environment
	file string
}

func CreateEnv(parent *Environment, file string) *Environment {
	return &Environment{
		Parent: parent,
		values: map[string]any{},
		local: map[string]any{},
		file: file,
	}
}

func (self Environment) Define(name string, value any, local bool) {
	if local {
		self.local[name] = value
		return
	}
	if self.Parent != nil {
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
	if file == self.file {
		val, ok := self.local[name]
		if ok {
			return val, nil
		}
	}
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
