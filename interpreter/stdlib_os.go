package interpreter

import (
	. "main/env"
	. "main/errors"
	"os"
	"os/exec"
	"time"
)


func CreateOs() *Environment {
	var temp = CreateEnv(nil, "")
	
	temp.Define("clock", Clock{name: "clock"}, false)
	temp.Define("time", Time{name: "time"}, false)
	temp.Define("execute", Execute{name: "execute"}, false)
	temp.Define("exit", Exit{name: "exit"}, false)

	return temp
}

type Clock Base

func (self Clock) Call(_ interpreter, args []any, line int) (any, error) {
	return time.Since(start).Seconds(), nil
}

func (self Clock) Arity() int {
	return 0
}

type Time Base

func (self Time) Call(_ interpreter, args []any, line int) (any, error) {
	return float64(time.Now().Unix()), nil
}

func (self Time) Arity() int {
	return 0
}

type Execute Base

func (self Execute) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case *StrType:
		a := []string{}
		for _, val2 := range args[1:len(args)-1] {
			switch val2 := val2.(type) {
			case *StrType:
				a = append(a, val2.str)
			default:
				return nil, RuntimeError("Invalid input type", line, filename)
			}
		}
		cmd := exec.Command(val.str, a...)
		stdout, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		return &StrType{str: string(stdout)}, nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Execute) Arity() int {
	return -1
}

type Exit Base

func (self Exit) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		os.Exit(int(val))
		return nil, nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Exit) Arity() int {
	return 1
}
