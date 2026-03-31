package interpreter

import (
	"fmt"
	. "main/env"
	. "main/errors"
	"math"
	"math/rand"
)

var r = rand.New(rand.NewSource(rand.Int63()))
func CreateMath() *Environment {
	var temp = CreateEnv(nil, "")

	// Functions
	temp.Define("abs", Abs{name: "abs"}, false)
	temp.Define("acos", Acos{name: "acos"}, false)
	temp.Define("asin", Asin{name: "asin"}, false)
	temp.Define("atan", Atan{name: "atan"}, false)
	temp.Define("atan2", Atan2{name: "atan2"}, false)
	temp.Define("ceil", Ceil{name: "ceil"}, false)
	temp.Define("cos", Cos{name: "cos"}, false)
	temp.Define("deg", Deg{name: "deg"}, false)
	temp.Define("floor", Floor{name: "floor"}, false)
	temp.Define("log", Log{name: "log"}, false)
	temp.Define("max", Max{name: "max"}, false)
	temp.Define("min", Min{name: "min"}, false)
	temp.Define("rad", Rad{name: "rad"}, false)
	temp.Define("random", Random{name: "random"}, false)
	temp.Define("seed", Seed{name: "seed"}, false)
	temp.Define("sin", Sin{name: "sin"}, false)
	temp.Define("sqrt", Sqrt{name: "sqrt"}, false)
	temp.Define("tan", Tan{name: "tan"}, false)
	
	// Constants
	temp.Define("pi", math.Pi, false)
	temp.Define("e", math.E, false)
	temp.Define("phi", math.Phi, false)
	temp.Define("maxint", math.MaxFloat64, false)
	temp.Define("minint", -math.MaxFloat64, false)
	temp.Define("posinf", math.Inf(1), false)
	temp.Define("neginf", math.Inf(-1), false)

	return temp
}

type Abs Base

func (self Abs) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return math.Abs(val), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Abs) Arity() int {
	return 1
}

type Acos Base

func (self Acos) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return math.Acos(val), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Acos) Arity() int {
	return 1
}

type Asin Base

func (self Asin) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return math.Asin(val), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Asin) Arity() int {
	return 1
}

type Atan Base

func (self Atan) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return math.Atan(val), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Atan) Arity() int {
	return 1
}

type Atan2 Base

func (self Atan2) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		switch val2 := args[1].(type) {
		case float64:
			return math.Atan2(val, val2), nil
		default:
			return nil, RuntimeError("Invalid input type", line, filename)
		}
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Atan2) Arity() int {
	return 2
}

type Ceil Base

func (self Ceil) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return float64(math.Ceil(val)), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Ceil) Arity() int {
	return 1
}

type Cos Base

func (self Cos) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return math.Cos(val), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Cos) Arity() int {
	return 1
}

type Deg Base

func (self Deg) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return (val/math.Pi) * 180, nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Deg) Arity() int {
	return 1
}

type Floor Base

func (self Floor) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return float64(math.Floor(val)), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Floor) Arity() int {
	return 1
}

type Log Base

func (self Log) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		switch val2 := args[1].(type) {
		case float64:
			return math.Log(val2) / math.Log(val), nil
		default:
			return nil, RuntimeError("Invalid input type", line, filename)
		}
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Log) Arity() int {
	return 2
}

type Max Base

func (self Max) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		switch val1 := args[1].(type) {
		case float64:
			return math.Max(val, val1), nil
		default:
			return nil, RuntimeError("Invalid input type", line, filename)
		}
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Max) Arity() int {
	return 2
}

type Min Base

func (self Min) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		switch val1 := args[1].(type) {
		case float64:
			return math.Min(val, val1), nil
		default:
			return nil, RuntimeError("Invalid input type", line, filename)
		}
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Min) Arity() int {
	return 2
}

type Rad Base

func (self Rad) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return (val*math.Pi) / 180, nil 
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Rad) Arity() int {
	return 1
}

type Random Base

func (self Random) Call(_ interpreter, args []any, line int) (any, error) {
	if len(args) == 1 || len(args) > 2 {
		return nil, RuntimeError(fmt.Sprintf("Expected 0 or 2 arguments, but got %d",
			len(args)), line, filename)
	}
	if len(args) == 0 {
		return r.Float64(), nil
	} else {
		switch val := args[0].(type) {
		case float64:
			switch val2 := args[1].(type) {
			case float64:
				if val2 < val {
					return nil, RuntimeError("Second input cannot be less than" +
						"first input", line, filename)
				}
				return float64(r.Intn(int(val2 - val)+1) + int(val)), nil
			default:
				return nil, RuntimeError("Invalid input type", line, filename)
			}
		default:
			return nil, RuntimeError("Invalid input type", line, filename)
		}
	}
}

func (self Random) Arity() int {
	return -1
}

type Seed Base

func (self Seed) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		r.Seed(int64(val))
		return nil, nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Seed) Arity() int {
	return 1
}

type Sin Base

func (self Sin) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return math.Sin(val), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Sin) Arity() int {
	return 1
}

type Sqrt Base

func (self Sqrt) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return math.Sqrt(val), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Sqrt) Arity() int {
	return 1
}

type Tan Base

func (self Tan) Call(_ interpreter, args []any, line int) (any, error) {
	switch val := args[0].(type) {
	case float64:
		return math.Tan(val), nil
	default:
		return nil, RuntimeError("Invalid input type", line, filename)
	}
}

func (self Tan) Arity() int {
	return 1
}

