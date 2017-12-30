package picol

import (
	"fmt"
	"strconv"
	"strings"
)

func ArityErr(i *Interp, name string, argv []interface{}) error {
	return fmt.Errorf("Wrong number of args for %s %#v", name, argv)
}

func CommandMath(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 3 {
		return "", ArityErr(i, argv[0].(string), argv)
	}
	a, _ := strconv.Atoi(argv[1].(string))
	b, _ := strconv.Atoi(argv[2].(string))
	var c int
	switch {
	case argv[0] == "+":
		c = a + b
	case argv[0] == "-":
		c = a - b
	case argv[0] == "*":
		c = a * b
	case argv[0] == "/":
		c = a / b
	case argv[0] == ">":
		if a > b {
			c = 1
		}
	case argv[0] == ">=":
		if a >= b {
			c = 1
		}
	case argv[0] == "<":
		if a < b {
			c = 1
		}
	case argv[0] == "<=":
		if a <= b {
			c = 1
		}
	case argv[0] == "==":
		if a == b {
			c = 1
		}
	case argv[0] == "!=":
		if a != b {
			c = 1
		}
	default: // FIXME I hate warnings
		c = 0
	}
	return fmt.Sprintf("%d", c), nil
}

func CommandSet(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 3 {
		return "", ArityErr(i, argv[0].(string), argv)
	}
	i.SetVar(argv[1].(string), argv[2])

	return "", nil
}

func CommandUnset(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 2 {
		return "", ArityErr(i, argv[0].(string), argv)
	}
	i.UnsetVar(argv[1].(string))
	return "", nil
}

func CommandIf(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 3 && len(argv) != 5 {
		return "", ArityErr(i, argv[0].(string), argv)
	}

	result, err := i.Eval(argv[1].(string))
	if err != nil {
		return "", err
	}

	if r, _ := strconv.Atoi(result.(string)); r != 0 {
		return i.Eval(argv[2].(string))
	} else if len(argv) == 5 {
		return i.Eval(argv[4].(string))
	}

	return result, nil
}

func CommandWhile(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 3 {
		return "", ArityErr(i, argv[0].(string), argv)
	}

	for {
		result, err := i.Eval(argv[1].(string))
		if err != nil {
			return "", err
		}
		if r, _ := strconv.Atoi(result.(string)); r != 0 {
			result, err := i.Eval(argv[2].(string))
			switch err {
			case PICOL_CONTINUE, nil:
				//pass
			case PICOL_BREAK:
				return result, nil
			default:
				return result, err
			}
		} else {
			return result, nil
		}
	}
}

func CommandRetCodes(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 1 {
		return "", ArityErr(i, argv[0].(string), argv)
	}
	switch argv[0].(string) {
	case "break":
		return "", PICOL_BREAK
	case "continue":
		return "", PICOL_CONTINUE
	}
	return "", nil
}

func CommandCallProc(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	var x []string

	if pd, ok := pd.([]string); ok {
		x = pd
	} else {
		return "", nil
	}

	i.callframe = &CallFrame{vars: make(map[string]interface{}), parent: i.callframe}
	defer func() { i.callframe = i.callframe.parent }() // remove the called proc callframe

	arity := 0
	for _, arg := range strings.Split(x[0], " ") {
		if len(arg) == 0 {
			continue
		}
		arity++
		i.SetVar(arg, argv[arity])
	}

	if arity != len(argv)-1 {
		return "", fmt.Errorf("Proc '%s' called with wrong arg num", argv[0])
	}

	body := x[1]
	result, err := i.Eval(body)
	if err == PICOL_RETURN {
		err = nil
	}
	return result, err
}

func CommandProc(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 4 {
		return "", ArityErr(i, argv[0].(string), argv)
	}
	return "", i.RegisterCommand(argv[1].(string), CommandCallProc, []string{argv[2].(string), argv[3].(string)})
}

func CommandReturn(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 1 && len(argv) != 2 {
		return "", ArityErr(i, argv[0].(string), argv)
	}

	// return type need be restricted to string
	var r interface{}
	if len(argv) == 2 {
		r = argv[1]
	}
	return r, PICOL_RETURN
}

func CommandError(i *Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 1 && len(argv) != 2 {
		return "", ArityErr(i, argv[0].(string), argv)
	}
	return "", fmt.Errorf(argv[1].(string))
}

func (i *Interp) RegisterCoreCommands() {
	name := [...]string{"+", "-", "*", "/", ">", ">=", "<", "<=", "==", "!="}
	for _, n := range name {
		i.RegisterCommand(n, CommandMath, nil)
	}
	i.RegisterCommand("set", CommandSet, nil)
	i.RegisterCommand("unset", CommandUnset, nil)
	i.RegisterCommand("if", CommandIf, nil)
	i.RegisterCommand("while", CommandWhile, nil)
	i.RegisterCommand("break", CommandRetCodes, nil)
	i.RegisterCommand("continue", CommandRetCodes, nil)
	i.RegisterCommand("proc", CommandProc, nil)
	i.RegisterCommand("return", CommandReturn, nil)
	i.RegisterCommand("error", CommandError, nil)
}
