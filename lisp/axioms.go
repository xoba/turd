package lisp

import (
	"fmt"
	"strings"
)

// valid types: string, []Exp, Func, or error
type Exp interface{}

type Func func(...Exp) Exp

var (
	env   Exp
	Nil   Exp = list()
	True  Exp = "t"
	False Exp = Nil
)

func String(e Exp) string {
	switch t := e.(type) {
	case string:
		return t
	case []Exp:
		if len(t) == 2 && t[0] == "quote" {
			return fmt.Sprintf("'%s", String(t[1]))
		}
		var list []string
		for _, e := range t {
			list = append(list, String(e))
		}
		return fmt.Sprintf("(%s)", strings.Join(list, " "))
	case Func:
		return String(t())
	case error:
		return fmt.Sprintf("error: %v", t)
	default:
		panic(fmt.Errorf("can't stringify type %T %v", t, t))
	}
}

func manifest(a Exp) Exp {
	if f, ok := a.(Func); ok {
		a = f()
	}
	return a
}

func one(args ...Exp) Exp {
	return manifest(args[0])
}

func two(args ...Exp) (Exp, Exp) {
	x, y := args[0], args[1]
	return manifest(x), manifest(y)
}

// ----------------------------------------------------------------------
// AXIOMS
// ----------------------------------------------------------------------

//
// axiom #1
//
func quote(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	return args[0]
}

//
// axiom #2
//
func atom(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	switch t := one(args...).(type) {
	case string:
		return True
	case []Exp:
		return boolToExp(len(t) == 0)
	default:
		return fmt.Errorf("illegal atom call: %T %v", t, t)
	}
}

//
// axiom #3
//
func eq(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x, y := two(args...)
	switch x := x.(type) {
	case string:
		switch y := y.(type) {
		case string: // both atoms
			return boolToExp(x == y)
		case []Exp:
			return False
		default:
			return fmt.Errorf("bad second argument to eq: %T", y)
		}
	case []Exp:
		switch y := y.(type) {
		case string:
			return False
		case []Exp: // both lists
			return boolToExp(len(x) == 0 && len(y) == 0)
		default:
			return fmt.Errorf("bad second argument to eq: %T", y)
		}
	default:
		return fmt.Errorf("bad first argument to eq: %T", x)
	}
}

//
// axiom #4
//
func car(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	switch t := one(args...).(type) {
	case []Exp:
		switch len(t) {
		case 0:
			return Nil
		default:
			return t[0]
		}
	default:
		return fmt.Errorf("car needs list, got %T %v", t, t)
	}
}

//
// axiom #5
//
func cdr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	switch t := one(args...).(type) {
	case []Exp:
		switch len(t) {
		case 0:
			return Nil
		default:
			return t[1:]
		}
	default:
		return fmt.Errorf("cdr needs list, got %T %v", t, t)
	}
}

//
// axiom #6
//
func cons(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := manifest(args[1])
	switch y.(type) {
	case []Exp:
		var out []Exp
		out = append(out, x)
		out = append(out, y.([]Exp)...)
		return out
	default:
		return fmt.Errorf("cons needs a list")
	}
}

//
// axiom #7
//
func cond(args ...Exp) Exp {
	if err := checkargs(args); err != nil {
		return err
	}
	for _, a := range args {
		switch t := a.(type) {
		case []Exp:
			if err := checklen(2, t); err != nil {
				return err
			}
			p, ok := t[0].(Func)
			if !ok {
				return fmt.Errorf("cond predicate not lazy")
			}
			if expToBool(p()) {
				e, ok := t[1].(Func)
				if !ok {
					return fmt.Errorf("cond exp not lazy")
				}
				return e()
			}
		default:
			return fmt.Errorf("illegal cond arg type %T", t)
		}
	}
	return fmt.Errorf("cond fallthrough")
}

//
// axiom #8
//
func display(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	a := one(args...)
	fmt.Printf("(display %s)\n", String(a))
	return a
}

//
// axiom #9
//
func printf(args ...Exp) Exp {
	if err := checkargs(args); err != nil {
		return err
	}
	if len(args) == 0 {
		return fmt.Errorf("printf needs at least one arg")
	}
	format, ok := car(manifest(args[0])).(string)
	if !ok {
		return fmt.Errorf("format needs to be a string")
	}
	var list []interface{}
	for _, a := range args[1:] {
		list = append(list, String(a))
	}
	fmt.Printf(format, list...)
	return Nil
}

//
// axiom #10 (kind of a like "quote" for multiple args)
//
func list(args ...Exp) Exp {
	return args
}
