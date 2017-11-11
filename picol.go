package picol

import (
	"errors"
	"fmt"
	"strings"
)

var (
	PICOL_RETURN   = errors.New("RETURN")
	PICOL_BREAK    = errors.New("BREAK")
	PICOL_CONTINUE = errors.New("CONTINUE")
	UNDEFINED_VAR  = errors.New("Undefined variable")
	CMD_EXISTS     = errors.New("Command with the given name is already registered")
)

// Not sure if this explicit type for interface it required, I would just go
// with inteface{} in all places, but lets this be there for now.
type Var interface{}

// by default, all args are assumed to be strings, I am changing them to
// interface{}, so we could pass custom data type objects to the commands, the
// return type should also be interface{} instead of the default string
type CmdFunc func(i *Interp, argv []interface{}, privdata interface{}) (interface{}, error)

type Cmd struct {
	fn       CmdFunc
	privdata interface{}
}

type CallFrame struct {
	vars   map[string]interface{}
	parent *CallFrame
}

type Interp struct {
	level     int
	callframe *CallFrame // always points to current procedure call frame
	commands  map[string]Cmd
}

func InitInterp() *Interp {
	return &Interp{
		level:     0,
		callframe: &CallFrame{vars: make(map[string]interface{})},
		commands:  make(map[string]Cmd),
	}
}

// Looks up for variables by iterating through the call frames, starting with
// the current to the root call frame, till a variable name match is found.
func (i *Interp) Var(name string) (interface{}, error) {
	for frame := i.callframe; frame != nil; frame = frame.parent {
		v, ok := frame.vars[name]
		if ok {
			return v, nil
		}
	}
	return "", UNDEFINED_VAR
}

func (i *Interp) SetVar(name string, val interface{}) {
	i.callframe.vars[name] = val
}

func (i *Interp) UnsetVar(name string) {
	delete(i.callframe.vars, name)
}

func (i *Interp) Command(name string) *Cmd {
	v, ok := i.commands[name]
	if !ok {
		return nil
	}
	return &v
}

func (i *Interp) RegisterCommand(name string, fn CmdFunc, privdata interface{}) error {
	c := i.Command(name)
	if c != nil {
		return CMD_EXISTS
	}

	i.commands[name] = Cmd{fn, privdata}
	return nil
}

/* EVAL! */
// Ideally Eval should take interface{}, which would help us either evaluate a
// string or an already parsed lexical tokens. For now we are starting with
// string and we could optimize on this at a later stage when we rewrite our
// parser module.
func (i *Interp) Eval(t string) (interface{}, error) {
	p := InitParser(t)
	var result interface{}
	var err error

	argv := []interface{}{}
	var ir interface{}

	for {
		prevtype := p.Type
		// XXX
		t = p.GetToken()
		if p.Type == PT_EOF {
			break
		}

		switch p.Type {
		case PT_VAR:
			v, err := i.Var(t)
			if err != nil {
				return "", UNDEFINED_VAR
			}
			ir = v
		case PT_CMD:
			result, err = i.Eval(t)
			if err != nil {
				// error in evaluating the argument
				return result, err
			} else {
				ir = result
			}
		case PT_ESC, PT_STR:
			// when the token is a simple string and requires no further
			// processing like variable substitution or command substitution
			ir = t
		case PT_SEP:
			prevtype = p.Type
			continue
		}

		// We have a complete command + args. Call it!
		if p.Type == PT_EOL {
			prevtype = p.Type
			if len(argv) != 0 {
				c := i.Command(argv[0].(string))
				if c == nil {
					return "", fmt.Errorf("No such command '%s'", argv[0])
				}
				result, err = c.fn(i, argv, c.privdata)
				if err != nil {
					return result, err
				}
			}
			// Prepare for the next command
			argv = []interface{}{}
			continue
		}

		// We have a new token, append to the previous or as new arg?
		if prevtype == PT_SEP || prevtype == PT_EOL {
			argv = append(argv, ir)
		} else { // String interpolation: variable substitution, command substitution
			argv[len(argv)-1] = strings.Join([]string{argv[len(argv)-1].(string), ir.(string)}, "")
		}
		prevtype = p.Type
	}
	return result, nil
}
