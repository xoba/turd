package main

import (
	"bytes"
	"fmt"
	"strings"
)

type Exp interface{}

type Func func(...Exp) Exp

var env Exp

var (
	Nil   Exp = list()
	t     Exp = "t"
	True  Exp = "t"
	False Exp = Nil
)

func main() {

	fmt.Printf("and = %v\n", env_and)
	fmt.Printf("and = %s\n", String(env_and))

	test := func(msg string, e Exp, expected string) {
		got := String(e)
		fmt.Printf("%s: %s\n", msg, got)
		if got != expected {
			fmt.Printf("*** expected %q, got %q\n", expected, got)
		}
	}
	test("1", quote("howdy"), "howdy")

	test("2", atom(quote("howdy")), "t")
	test("11", atom(list("quote", "a")), "()")
	test("12", atom("x"), "t")

	test("3", apply(atom, apply(quote, "howdy")), "t")
	test("4", atom(list()), "t")
	test("5", atom(list("a")), "()")
	test("6", eq("a", "a"), "t")
	test("7", eq("a", "b"), "()")

	f1 := and("t", "t")

	test("14", f1, "t")
	test("15", and("t", list()), "()")

	lazy := func(e Exp) Func {
		return Func(func(...Exp) Exp {
			return e
		})
	}

	test("8", cond(
		list(lazy(True), lazy(quote("a"))),
		list(lazy(True), lazy(quote("b"))),
	), "a")
	test("9", cond(
		list(lazy(False), lazy(quote("a"))),
		list(lazy(True), lazy(quote("b"))),
	), "b")

	test("10", cons("a", list("b", "c")), "(a b c)")

	e, a := list("quote", "x"), Nil
	test("13", apply(
		atom,
		e,
	), "()")

	test("12", eval(e, a), "x")
}

func String(e Exp) string {
	w := new(bytes.Buffer)
	switch t := e.(type) {
	case string:
		fmt.Fprint(w, t)
	case []Exp:
		var list []string
		for _, e := range t {
			list = append(list, String(e))
		}
		fmt.Fprintf(w, "(%s)", strings.Join(list, " "))
	case Func:
		return String(t())
	default:
		// TODO: why do we get "panic: can't stringify type int 0"?
		panic(fmt.Errorf("can't stringify type %T %v", t, t))
	}
	return w.String()
}

func IsAtom(e Exp) bool {
	switch e.(type) {
	case string:
		return true
	default:
		return false
	}
}

func apply(f Func, args ...Exp) Exp {
	return f(args...)
}

func checklen(n int, args []Exp) {
	if len(args) != n {
		panic(fmt.Errorf("expected %d args, got %d", n, len(args)))
	}
}

func list(args ...Exp) Exp {
	return args
}

func boolean(e Exp) bool {
	return fmt.Sprintf("%v", e) == "t"
}

// ----------------------------------------------------------------------
// AXIOMS
// ----------------------------------------------------------------------

//
// #1 (1-8 from paul graham)
//
func quote(args ...Exp) Exp {
	checklen(1, args)
	return args[0]
}

//
// #2
//

func atom(args ...Exp) Exp {
	checklen(1, args)
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
		panic("illegal atom call")
	}
}

//
// #3
//

func eq(args ...Exp) Exp {
	checklen(2, args)
	s := func(e Exp) string {
		return fmt.Sprintf("%s", e)
	}
	x, y := args[0], args[1]
	if s(x) == s(y) {
		return True
	}
	return False
}

//
// #4
//

func car(args ...Exp) Exp {
	checklen(1, args)
	return args[0]
}

//
// #5
//

func cdr(args ...Exp) Exp {
	checklen(1, args)
	return args[1:]
}

//
// #6
//

func cons(args ...Exp) Exp {
	checklen(2, args)
	x, y := args[0], args[1]
	if IsAtom(y) {
		panic("cons to atom")
	}
	slice := y.([]Exp)
	var out []Exp
	out = append(out, x)
	out = append(out, slice...)
	return out
}

//
// #7
//

func cond(args ...Exp) Exp {
	for i, a := range args {
		switch t := a.(type) {
		case []Exp:
			if len(t) != 2 {
				panic(fmt.Errorf("len[%d] = %d", i, len(t)))
			}
			p, e := t[0], t[1]
			pl, ok := p.(Func)
			if !ok {
				panic("p not lazy")
			}
			if boolean(pl()) {
				el, ok := e.(Func)
				if !ok {
					panic("e not lazy")
				}
				return el()
			}
		default:
			panic(fmt.Errorf("cond %T", t))
		}
	}
	panic(fmt.Errorf("cond fallthrough with %d args", len(args)))
}

//
// #8 (from chaitin)
//

func display(args ...Exp) Exp {
	checklen(1, args)
	a := args[0]
	fmt.Printf("(display %s)\n", a)
	return a
}
