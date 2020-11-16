package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xoba/turd/lisp"
	"github.com/xoba/turd/lisp/exp"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ToExp(e exp.Expression) Exp {
	if err := e.Error(); err != nil {
		panic(err)
	}
	if a := e.Atom(); a != nil {
		return a.String()
	}
	var list []Exp
	for _, x := range e.List() {
		list = append(list, ToExp(x))
	}
	return list
}

func main() {

	{

		var last string
		test := func(msg, input, expect string) {
			if msg == "" {
				return
			}
			if msg != last {
				fmt.Println()
			}
			last = msg
			in, err := lisp.Read(input)
			check(err)
			fmt.Printf("%-10s %-20s -> %s\n", msg+":", in, expect)
			in = lisp.SanitizeGo(in)
			res := eval(ToExp(in), env)
			if got := String(res); got != expect {
				panic(fmt.Errorf("expected %q, got %q\n", expect, got))
			}
		}

		test("quote", "(quote a)", "a")
		test("quote", "(quote (a b c))", "(a b c)")

		test("atom", "(atom 'a)", "t")
		test("atom", "(atom '(a b c))", "()")
		test("atom", "(atom '())", "t")
		test("atom", "(atom 'a)", "t")
		test("atom", "(atom '(atom 'a))", "()")

		test("eq", "(eq 'a 'a)", "t")
		test("eq", "(eq 'a 'b)", "()")
		test("eq", "(eq '() '())", "t")

		test("car", "(car '(a b c))", "a")
		test("cdr", "(cdr '(a b c))", "(b c)")

		test("cons", "(cons 'a '(b c))", "(a b c)")
		test("cons", "(cons 'a (cons 'b (cons 'c '())))", "(a b c)")
		test("cons", "(car (cons 'a '(b c)))", "a")
		test("cons", "(cdr (cons 'a '(b c)))", "(b c)")

		test("cond", "(cond ((eq 'a 'b) 'first) ((atom 'a) 'second))", "second")
		test("cond", "(cond ((eq 'a 'a) 'first) ((atom 'a) 'second))", "first")

		test("lambda", "((lambda (x) (cons x '(b))) 'a)", "(a b)")

		test("label", `(
 (label subst 
	(lambda (x y z)
	  (cond ((atom z) (
			   cond ((eq z y) x)
				('t z)))
		('t (cons (subst x y (car z))
			  (subst x y (cdr z))))))
	)
 'm 'b '(a b (a b c) d))`, "(a m (a m c) d)")
		test("label", "(subst 'm 'b '(a b (a b c) d))", "(a m (a m c) d)")

		test("cxr", "(cadr '((a b) (c d) e))", "(c d)")
		test("cxr", "(caddr '((a b) (c d) e))", "e")
		test("cxr", "(cdar '((a b) (c d) e))", "(b)")

		test("list", "(cons 'a (cons 'b (cons 'c '())))", "(a b c)")
		test("list", "(list 'a 'b 'c)", "(a b c)")

		test("null", "(null 'a)", "()")
		test("null", "(null '())", "t")

		test("and", "(and (atom 'a) (eq 'a 'a))", "t")
		test("and", "(and (atom 'a) (eq 'a 'b))", "()")

		test("not", "(not (eq 'a 'a))", "()")
		test("not", "(not (eq 'a 'b))", "t")

		test("append", "(append '(a b) '(c d))", "(a b c d)")
		test("append", "(append '() '(c d))", "(c d)")

		test("", "", "")
		test("", "", "")
		test("", "", "")
		test("", "", "")
		test("", "", "")
		test("", "", "")
		test("", "", "")
		test("", "", "")
		test("", "", "")
		test("", "", "")
		test("", "", "")

	}

	return

	test := func(msg string, e Exp, expected string) {
		got := String(e)
		fmt.Printf("%s: %s\n", msg, got)
		if got != expected {
			fmt.Printf("*** expected %q, got %q\n", expected, got)
		}
	}
	test("0", testing(list("abc", "2")), "abc")
	//return

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

	fmt.Printf("eval(%s)\n", e)
	test("12", eval(e, a), "x")

	{
		e := list("eq", list("quote", "a"), list("quote", "a"))
		fmt.Printf("testing %s\n", e)
		test("16", eval(e, a), "t")
	}
	{
		e := list("eq", list("quote", "a"), list("quote", "b"))
		fmt.Printf("testing %s\n", e)
		test("16", eval(e, a), "()")
	}
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
		panic(fmt.Errorf("can't stringify type %T %v", t, t))
	}
	return w.String()
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
// #1
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
		panic(fmt.Errorf("illegal atom call: %T %v", x, x))
	}
}

//
// #3
//

func eq(args ...Exp) Exp {
	checklen(2, args)
	s := func(e Exp) string {
		return fmt.Sprintf("%T %s", e, e)
	}
	x, y := args[0], args[1]
	sx, sy := s(x), s(y)
	if sx == sy {
		return True
	}
	return False
}

//
// #4
//

func car(args ...Exp) Exp {
	checklen(1, args)
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
		panic("car needs list")
	}
}

//
// #5
//

func cdr(args ...Exp) Exp {
	checklen(1, args)
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
		panic("cdr needs list")
	}
}

//
// #6
//

func cons(args ...Exp) Exp {
	checklen(2, args)
	x, y := args[0], args[1]
	IsAtom := func(e Exp) bool {
		switch e.(type) {
		case string:
			return true
		case []Exp:
			return false
		default:
			panic("illegal type in cons")
		}
	}
	if IsAtom(y) {
		panic(fmt.Errorf("cons atom %T %v", t, t))
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
	for _, a := range args {
		switch t := a.(type) {
		case []Exp:
			checklen(2, t)
			p, e := t[0], t[1]
			pl, ok := p.(Func)
			if !ok {
				panic("p not lazy")
			}
			v := pl()
			if boolean(v) {
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
// #8
//

func display(args ...Exp) Exp {
	checklen(1, args)
	a := args[0]
	fmt.Printf("(display %s)\n", a)
	return a
}
