package lisp

import "fmt"

// valid types: string, []Exp, Func, or error
type Exp interface{}

type Func func(...Exp) Exp

var (
	env   Exp
	Nil   Exp = list()
	True  Exp = "t"
	False Exp = Nil
)

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
	return args[0]
}

//
// #2
//

func atom(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	switch t := x.(type) {
	case string:
		return True
	case []Exp:
		if len(t) == 0 {
			return True
		}
		return False
	default:
		return fmt.Errorf("illegal atom call: %T %v", t, t)
	}
}

//
// #3
//

func eq(args ...Exp) Exp {
	out := eq0(args...)
	//fmt.Printf("eq(%q,%q) = %q\n", args[0], args[1], out)
	return out
}

func eq0(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x, y := args[0], args[1]
	switch tx := x.(type) {
	case string:
		switch ty := y.(type) {
		case string: // both atoms
			return boolToExp(tx == ty)
		default:
			return False
		}
	case []Exp:
		switch ty := y.(type) {
		case []Exp: // both lists
			return boolToExp(len(tx) == 0 && len(ty) == 0)
		default:
			return False
		}
	default:
		return fmt.Errorf("bad eq arguments")
	}
}

//
// #4
//

func car(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	switch t := x.(type) {
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
	x := args[0]
	switch t := x.(type) {
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
	x, y := args[0], args[1]
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
	a := args[0]
	fmt.Printf("(display %s)\n", String(a))
	return a
}

//
// #9 (kind of a like "quote" for multiple args)
//

func list(args ...Exp) Exp {
	return args
}
