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
	default:
		panic(fmt.Errorf("can't stringify type %T %v", t, t))
	}
}

func one(args ...Exp) Exp {
	x := args[0]
	if f, ok := x.(Func); ok {
		x = f()
	}
	return x
}

func two(args ...Exp) (Exp, Exp) {
	x, y := args[0], args[1]
	if f, ok := x.(Func); ok {
		x = f()
	}
	if f, ok := y.(Func); ok {
		x = f()
	}
	return x, y
}

// ----------------------------------------------------------------------
// AXIOMS
// ----------------------------------------------------------------------

//
// #1
//
func quote(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	return one(args...)
}

//
// #2
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
// #3
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
// #4
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
// #5
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
// #6
//

func cons(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x, y := two(args...)
	switch y.(type) {
	case []Exp:
	default:
		return fmt.Errorf("cons needs a list")
	}
	var out []Exp
	out = append(out, x)
	out = append(out, y.([]Exp)...)
	return out
}

//
// #7
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
			p, e := t[0], t[1]
			pl, ok := p.(Func)
			if !ok {
				return fmt.Errorf("p not lazy")
			}
			v := pl()
			if expToBool(v) {
				el, ok := e.(Func)
				if !ok {
					return fmt.Errorf("e not lazy")
				}
				return el()
			}
		default:
			return fmt.Errorf("cond %T", t)
		}
	}
	return fmt.Errorf("cond fallthrough")
}

//
// #8
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
// #9 (kind of a like "quote" for multiple args)
//

func list(args ...Exp) Exp {
	return args
}
