package main

import (
	"bytes"
	"fmt"
	"strings"
)

var env Exp

func main() {

	show := func(msg string, e Exp) {
		fmt.Printf("%s: %s\n", msg, StringLazy(e, true))
	}
	show("1", quote("howdy"))

	show("2", atom(quote("howdy")))
	show("11", atom(list("quote", "a")))
	show("12", atom("x"))

	show("3", apply(atom, apply(quote, "howdy")))
	show("4", atom(list()))
	show("5", atom(list("a")))
	show("6", eq("a", "a"))
	show("7", eq("a", "b"))

	f1 := and("t", "t")

	show("14", f1)
	show("15", and("t", list()))

	lazy := func(e Exp) Lazy {
		return Lazy(func() Exp {
			return e
		})
	}

	show("8", cond(
		list(lazy(True), lazy(quote("a"))),
		list(lazy(True), lazy(quote("b"))),
	))
	show("9", cond(
		list(lazy(False), lazy(quote("a"))),
		list(lazy(True), lazy(quote("b"))),
	))

	show("10", cons("a", list("b", "c")))

	e, a := list("quote", "x"), Nil
	show("13", apply(
		atom,
		e,
	))
	show("12", eval(e, a))
}

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

type Exp interface{}

func String(e Exp) string {
	return StringLazy(e, false)
}

func StringLazy(e Exp, evalLazy bool) string {
	show := func(f func(...Exp) Exp) string {
		if evalLazy {
			v := f()
			return String(v)
		}
		return String(e)
	}
	w := new(bytes.Buffer)
	switch t := e.(type) {
	case string:
		fmt.Fprint(w, t)
	case Lazy:
		return show(func(...Exp) Exp {
			return t()
		})
	case func() Exp:
		return show(func(...Exp) Exp {
			return t()
		})
	case Func:
		return show(func(...Exp) Exp {
			return t()
		})
	case []Exp:
		var list []string
		for e := range t {
			list = append(list, String(e))
		}
		fmt.Fprintf(w, "(%s)", strings.Join(list, " "))
	default:
		panic(fmt.Errorf("can't stringify type %T %v", t, t))
	}
	return w.String()
}

type Func func(...Exp) Exp

type Lazy func() Exp

var (
	Nil    Exp  = list()
	True   Exp  = "t"
	False  Exp  = Nil
	t      Func = nil
	lambda Exp  = "lambda"
	label  Exp  = "label"
)

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

func quote(args ...Exp) Exp {
	checklen(1, args)
	return args[0]

	switch len(args) {
	case 0:
		return "quote"
	case 1:
		return list("quote", args[0])
	default:
		panic("quote")
	}
}

func checklen(n int, args []Exp) {
	if len(args) != n {
		panic(fmt.Errorf("len = %d vs %d", len(args), n))
	}
}

func car(args ...Exp) Exp {
	checklen(1, args)
	return args[0]
}

func cdr(args ...Exp) Exp {
	checklen(1, args)
	return args[1:]
}

func atom(args ...Exp) Exp {
	checklen(1, args)
	if len(args) != 1 {
		panic("args")
	}
	switch t := args[0].(type) {
	case string:
		return True
	case []Exp:
		if len(t) == 0 {
			return True
		}
		return False
	default:
		return False
	}
}

func eq(args ...Exp) Exp {
	checklen(2, args)
	s := func(e Exp) string {
		return fmt.Sprintf("%s", e)
	}
	if s(args[0]) == s(args[1]) {
		return True
	}
	return False
}

func list(args ...Exp) Exp {
	return args
}

func boolean(e Exp) bool {
	return fmt.Sprintf("%v", e) == "t"
}

func debug(name string, args ...Exp) {
	return
	fmt.Printf("%s%s\n", name, String(list(args...)))
}

func cond(args ...Exp) Exp {
	debug("cond", args...)
	for i, a := range args {
		switch t := a.(type) {
		case []Exp:
			if len(t) != 2 {
				panic(fmt.Errorf("len[%d] = %d", i, len(t)))
			}
			p, e := t[0], t[1]
			fmt.Printf("p,e = %s, %s\n", String(p), String(e))
			pl, ok := p.(func() Exp)
			if !ok {
				panic("p not lazy")
			}
			if boolean(pl()) {
				el, ok := e.(func() Exp)
				if !ok {
					panic("e not lazy")
				}
				return el()
			}
		default:
			panic(fmt.Errorf("cond %T", t))
		}
	}
	panic("cond")
}
