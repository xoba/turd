package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {

	show := func(e Exp) {
		fmt.Println(String(e))
	}
	show(quote())
	show(quote("howdy"))
	show(atom(quote("howdy")))
	show(apply(atom, apply(quote, "howdy")))
	show(atom(list()))
	show(atom(list("a")))
	show(eq(quote("a"), quote("a")))
	show(eq(quote("b"), quote("a")))

	lazy := func(e Exp) Lazy {
		return Lazy(func() Exp {
			return e
		})
	}

	show(cond(
		list(lazy(True), lazy(quote("a"))),
		list(lazy(True), lazy(quote("b"))),
	))
	show(cond(
		list(lazy(False), lazy(quote("a"))),
		list(lazy(True), lazy(quote("b"))),
	))

	show(cons(quote("a"), list(quote("b"), quote("c"))))
}

func cons(args ...Exp) Exp {
	checklen(2, args)
	x, y := args[0], args[1]
	if IsAtom(y) {
		panic("cons to atom")
	}
	slice := y.([]Exp)
	out := make([]Exp, len(slice)+1)
	out[0] = x
	for i, e := range slice {
		out[i+1] = e
	}
	return out
}

type Exp interface{}

func String(e Exp) string {
	w := new(bytes.Buffer)
	switch t := e.(type) {
	case string:
		fmt.Fprint(w, t)
	case []Exp:
		list := make([]string, len(t))
		for i, e := range t {
			list[i] = String(e)
		}
		fmt.Fprintf(w, "(%s)", strings.Join(list, " "))
	default:
		panic(fmt.Errorf("exp type %T", t))
	}
	return w.String()
}

type Func func(...Exp) Exp

type Lazy func() Exp

var (
	Nil    Exp  = list()
	True   Exp  = quote("t")
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
	switch len(args) {
	case 0:
		return "quote"
	case 1:
		return args[0]
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

func cond(args ...Exp) Exp {
	fmt.Printf("cond(%v)\n", args)
	for i, a := range args {
		switch t := a.(type) {
		case []Exp:
			if len(t) != 2 {
				panic(fmt.Errorf("len[%d] = %d", i, len(t)))
			}
			p, e := t[0], t[1]
			pl, ok := p.(Lazy)
			if !ok {
				panic("p not lazy")
			}
			if boolean(pl()) {
				el, ok := e.(Lazy)
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
