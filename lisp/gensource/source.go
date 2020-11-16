package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/xoba/turd/lisp"
	"github.com/xoba/turd/lisp/exp"
)

// valid types: string, []Exp, Func, or error
type Exp interface{}

type Func func(...Exp) Exp

var env Exp

var (
	Nil   Exp = list()
	True  Exp = "t"
	False Exp = Nil
)

func ToExp(e exp.Expression) Exp {
	if err := e.Error(); err != nil {
		return err
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

func Eval(e Exp) Exp {
	return eval([]Exp{e, env}...)
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {

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
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s %-20s -> %s\n", msg+":", in, expect)
		in = lisp.SanitizeGo(in)
		res := Eval(ToExp(in))
		if got := String(res); got != expect {
			log.Fatalf("expected %q, got %q\n", expect, got)
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

	test("pair", "(pair '(x y z) '(a b c))", "((x a) (y b) (z c))")

	test("assoc", "(assoc 'x '((x a) (y b)))", "a")
	test("assoc", "(assoc 'x '((x c) (y b)))", "c")
	test("assoc", "(assoc 'y '((x c) (y b)))", "b")

	test("eval", "(eval 'x '((x a) (y b)))", "a")
	test("eval", "(eval '(eq 'a 'a) '())", "t")
	test("eval", "(eval '(cons x '(b c)) '((x a) (y b)))", "(a b c)")
	test("eval", "(eval '(cond ((atom x) 'atom) ('t 'list)) '((x '(a b))))", "list")
	test("eval", "(eval '(f '(b c)) '((f (lambda (x) (cons 'a x)))))", "(a b c)")
	test("eval", "(eval '((label firstatom (lambda (x) (cond ((atom x) x) ('t (firstatom (car x)))))) y) '((y ((a b) (c d)))))", "a")
	test("eval", "(eval '((lambda (x y) (cons x (cdr y))) 'a '(b c d)) '())", "(a c d)")
	test("", "", "")
	test("", "", "")
	test("", "", "")

	return nil
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

func checkargs(args []Exp) error {
	for _, a := range args {
		switch t := a.(type) {
		case string:
		case []Exp:
		case Func:
		case error:
			return t
		default:
			return fmt.Errorf("illegal type: %T %v", t, t)
		}
	}
	return nil
}

func apply(f Func, args ...Exp) Exp {
	if err := checkargs(args); err != nil {
		return err
	}
	return f(args...)
}

func checklen(n int, args []Exp) error {
	if len(args) != n {
		return fmt.Errorf("expected %d args, got %d", n, len(args))
	}
	if err := checkargs(args); err != nil {
		return err
	}
	return nil
}

func expToBool(e Exp) bool {
	switch t := e.(type) {
	case string:
		return t == "t"
	default:
		return false
	}
}

func boolToExp(v bool) Exp {
	if v {
		return True
	}
	return False
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
	return fmt.Errorf("cond fallthrough with %d args", len(args))
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
