package scr

import (
	"fmt"
	"log"
	"strings"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/scr/exp"
)

func Lisp(config cnfg.Config) error {

	defer func() {
		fmt.Printf("test coverage = %v\n", coverage)
	}()

	if config.Debug {
		return TestCond()
	}

	m := make(map[string]bool)

	a := exp.NewList()

	check := func(e error) {
		if e != nil {
			log.Fatal(e)
		}
	}

	wrap := func(i interface{}, e error) error {
		if e == nil {
			return nil
		}
		return fmt.Errorf("#%v. %w", i, e)
	}

	var count int
	test := func(i interface{}, in, expect string) {
		if in == "" {
			return
		}
		if m[in] {
			panic("duplication: " + in)
		}
		count++
		m[in] = true
		e := Read(in)
		fmt.Printf("%2d.%v: %-55s -> %s\n", count, i, e, expect)
		check(wrap(i, e.Error()))
		x := Eval(e, a)
		check(wrap(i, x.Error()))
		if got := x.String(); got != expect {
			check(wrap(i, fmt.Errorf("expected %q, got %q", expect, got)))
		}
	}

	define := func(name, lambda string) {
		if name == "" {
			return
		}
		e := Read(fmt.Sprintf("(label %s %s)", name, lambda))
		check(wrap(name, e.Error()))
		a = exp.NewList(exp.NewList(exp.NewString(name), e), a)
	}

	test2 := func(x, y string) {
		test(0, x, y)
	}

	if true {
		const (
			Null = `
(label null
       (lambda (x) (eq x '())))
`
			Append = `
(label append
       (lambda (x y)
	 (cond ((null x) y)
	       ('t (cons (car x) (append (cdr x)))))))
`
		)
		test("funcs", "("+Null+" '())", "t")
		test("funcs", "("+Null+" 'a)", "()")

		Append2 := strings.Replace(Append, "null", Null, -1)
		test("funcs", "("+Append2+" '(a b) '(c d))", "(a b c d)")
		test("funcs", "("+Append2+" '() '(c d))", "(c d)")

		define("null", Null)

		test("funcs", "("+Append+" '(a b) '(c d))", "(a b c d)")
		test("funcs", "("+Append+" '() '(c d))", "(c d)")
		//return nil

	}

	if true {
		// TODO: this group works at start of test suite, but not in middle or end!
		const null = "(lambda (x) (eq x '()))"
		define("null", null)
		test("funcs", "("+null+" '())", "t")
		test("funcs", "("+null+" 'a)", "()")
		test("funcs", "(null 'a)", "()")
		test("funcs", "(null '())", "t")
		const (
			def = "(lambda (x y) (cond ((" + null + " x) y) ('t (cons (car x) (append (cdr x) y)))))"
		)
		test("funcs", "("+def+" '(c d))", "(a b c d)")
		return nil
		define("append", def)
		test("funcs", `(append '(a b) '(c d))`, "(a b c d)")
		test("funcs", `(append '() '(c d))`, "(c d)")
	}

	// other:
	test2("(quote z)", "z")
	test2("(atom '(1 2 3))", "()")
	test2("(eq 'b 'a)", "()")
	test2("(eq 'b '())", "()")
	test2("(car '(x))", "x")
	test2("(cdr '(x))", "()")

	// 1
	test(1, "(quote a)", "a")
	test(1, "'a", "a")
	test(1, "(quote (a b c))", "(a b c)")

	// 2
	test(2, "(atom 'a)", "t")
	test(2, "(atom '(a b c))", "()")
	test(2, "(atom '())", "t")
	test(2, "(atom (atom 'a))", "t")
	test(2, "(atom '(atom 'a))", "()")

	// 3
	test(3, "(eq 'a 'a)", "t")
	test(3, "(eq 'a 'b)", "()")
	test(3, "(eq '() '())", "t")

	// 4
	test(4, "(car '(a b c))", "a")

	// 5
	test(5, "(cdr '(a b c))", "(b c)")

	// 6
	test(6, "(cons 'a '(b c))", "(a b c)")
	test(6, "(cons 'a (cons 'b (cons 'c '())))", "(a b c)")
	test(6, "(car (cons 'a '(b c)))", "a")
	test(6, "(cdr (cons 'a '(b c)))", "(b c)")

	// 7
	test(7, `

(cond ((eq 'a 'b) 'first) 
      ((atom 'a) 'second))

`, "second")

	// lambda
	test("λ", "((lambda (x) (cons x '(b))) 'a)", "(a b)")
	test("λ", "((lambda (x y) (cons x (cdr y))) 'z '(a b c))", "(z b c)")
	test("λ", "((lambda (f) (f '(b c))) '(lambda (x) (cons 'a x)))", "(a b c)")

	test("label", `

((label subst 
	(lambda (x y z)
	  (cond ((atom z)
		 (cond ((eq z y) x)
		       ('t z)))
		('t (cons (subst x y (car z))
			  (subst x y (cdr z)))))))
 'm 'b '(a b (a b c) d))

`, "(a m (a m c) d)")

	test("funcs", "(cadr '((a b) (c d) e))", "(c d)")
	test("funcs", "(caddr '((a b) (c d) e))", "e")
	test("funcs", "(cdar '((a b) (c d) e))", "(b)")

	test("list", "(list 'a 'b 'c)", "(a b c)")

	define("and", `

(lambda (x y)
   (cond (x (cond (y 't) ('t ())))
	 ('t '())))
`)

	test("funcs", "(and (atom 'a) (eq 'a 'a))", "t")
	test("funcs", "(and (atom 'a) 't)", "t")
	test("funcs", "(and (atom 'a) '())", "()")

	define("not", "(lambda (x) (cond (x '()) ('t 't)))")

	test("funcs", `(not (eq 'a 'a))`, "()")
	test("funcs", `(not (eq 'a 'b))`, "t")

	define("pair", `
(lambda (x y) (cond ((and (null x) (null y)) '())
		    ((and (not (atom x)) (not (atom y)))
		     (cons (list (car x) (car y))
			   (pair (cdr x) (cdr y))))))`)
	define("assoc", `(lambda (x y)
  (cond ((eq (caar y) x) (cadar y))
	('t (assoc x (cdr y)))))
`)

	define("eval", `(lambda (e a)
  (cond
   ((atom e) (assoc e a))
   ((atom (car e))
    (cond
     ((eq (car e) 'quote)  (cadr e))
     ((eq (car e) 'atom)   (atom   (eval (cadr e)  a)))
     ((eq (car e) 'eq)     (eq     (eval (cadr e)  a)
				   (eval (caddr e) a)))
     ((eq (car e) 'car)    (car    (eval (cadr e)  a)))
     ((eq (car e) 'cdr)    (cdr    (eval (cadr e)  a)))
     ((eq (car e) 'cons)   (cons   (eval (cadr e)  a)
				   (eval (caddr e) a)))
     ((eq (car e) 'cond)   (evcon (cdr e) a))
     ('t                   (eval  (cons (assoc (car e) a)
					(cdr e))
				  a))))
   ((eq (caar e) 'label)
    (eval (cons (caddar e) (cdr e))
	  (cons (list (cadar e) (car e)) a)))
   ((eq (caar e) 'lambda)
    (eval (caddar e)
	  (append (pair (cadar e) (evlis (cdr e) a))
		  a)))))
`)

	define("evcon", `(lambda(c a)
  (cond ((eval (caar c) a)
	 (eval (cadar c) a))
	('t (evcon (cdr c) a))))
`)

	define("evlis", `(lambda (m a)
  (cond ((null m) '())
	('t (cons (eval  (car m) a)
		  (evlis (cdr m) a)))))
`)

	if false {
		// make sure this errors out:
		test("funcs", `(xyz 'a)`, "")
	}

	test("funcs", `(pair '(x y z) '(a b c))`, "((x a) (y b) (z c))")
	test("funcs", `(assoc x '((x a) (y b)))`, "a")
	test("funcs", `(assoc 'x '((x new) (x a) (y b)))`, "new")
	test("funcs", `(eval 'x '((x a) (y b)))`, "a")
	test("funcs", `(eval '(eq 'a 'a) '())`, "t")
	test("funcs", `(eval '(cons x '(b c)) '((x a) (y b)))`, "(a b c)")
	test("funcs", `(eval '(cond ((atom x) 'atom) ('t 'list)) '((x '(a b))))`, "list")
	test("funcs", `(eval '(f '(b c)) '((f (lambda (x) (cons 'a x)))))`, "(a b c)")
	test("funcs", `(eval '((label firstatom (lambda (x)
			   (cond ((atom x) x)
				 ('t (firstatom (car x))))))
	y)
      '((y ((a b) (c d)))))`, "a")
	test("funcs", `(eval '((lambda (x y) (cons x (cdr y)))
	'a
	'(b c d))
      '())
`, "(a c d)")
	test("funcs", ``, "")
	test("funcs", ``, "")

	test(0, "", "")
	test(0, "", "")
	test(0, "", "")
	test(0, "", "")

	return nil
}

// vars for test coverage:
var (
	coverage map[string]int
	reg      func(string)
)

func init() {
	coverage = make(map[string]int)
	reg = func(x string) {
		coverage[x]++
	}
}

func Assoc(x, y exp.Expression) exp.Expression {
	reg("assoc")
	if !IsAtom(x) {
		return exp.NewError(fmt.Errorf("needs an atom to assoc"))
	}
	switch {
	case IsList(y) && IsEmpty(y):
		return x
	case !IsList(y):
		return exp.NewError(fmt.Errorf("needs a list to assoc"))
	}
	caar := Car(Car(y))
	eq := Eq(caar, x)
	if eq.String() == "t" {
		return Car(Cdr(Car(y)))
	}
	cdr := Cdr(y)
	return Assoc(x, cdr)
}

func Eq(x, y exp.Expression) exp.Expression {
	reg("eq")
	r := func(v bool) exp.Expression {
		if v {
			return exp.NewString("t")
		}
		return exp.NewList()
	}
	switch {
	case IsAtom(x) && IsAtom(y):
		xa := x.Atom()
		ya := y.Atom()
		return r(AtomsEqual(xa, ya))
	case IsList(x) && IsList(y):
		return r(IsEmpty(x) && IsEmpty(y))
	default:
		return r(false)
	}
}

func cxr(s string) (MonadFunc, error) {
	runes := []rune(s)
	n := len(runes) - 1
	if runes[0] != 'c' || runes[n] != 'r' {
		return nil, fmt.Errorf("not CxR: %q", s)
	}
	var list []MonadFunc
	for i := 1; i < n; i++ {
		var f MonadFunc
		switch runes[i] {
		case 'a':
			f = Car
		case 'd':
			f = Cdr
		default:
			return nil, fmt.Errorf("not CxR: %q", s)
		}
		list = append(list, f)
	}
	return Compose(list...), nil
}

type efunc func(...exp.Expression) exp.Expression
type onefunc func(exp.Expression) exp.Expression
type twofunc func(exp.Expression, exp.Expression) exp.Expression

func two(f twofunc) efunc {
	return func(args ...exp.Expression) exp.Expression {
		if n := len(args); n != 2 {
			return exp.Errorf("needs two args, got %d", n)
		}
		return f(args[0], args[1])
	}

}

func one(f onefunc) efunc {
	return func(args ...exp.Expression) exp.Expression {
		if n := len(args); n != 1 {
			return exp.Errorf("needs two args, got %d", n)
		}
		return f(args[0])
	}
}

func apply(f efunc, args ...exp.Expression) exp.Expression {
	for _, arg := range args {
		if arg.Error() != nil {
			return arg
		}
	}
	return f(args...)
}

func x(s string) efunc {
	f, err := cxr(s)
	if err != nil {
		panic(err)
	}
	return one(func(a exp.Expression) exp.Expression {
		reg(s)
		return f(a)
	})
}

var (
	caar   = x("caar")
	cadar  = x("cadar")
	caddar = x("caddar")
	caddr  = x("caddr")
	cadr   = x("cadr")
	car    = x("car")
	cdr    = x("cdr")
)

func Cond(args ...exp.Expression) exp.Expression {
	for _, a := range args {
		if Boolean(car(a)) {
			return cadr(a)
		}
	}
	return exp.Errorf("cond fallthrough")
}

func TestCond() error {
	lazy := func(name string, e exp.Expression) exp.Expression {
		var done bool
		return exp.NewLazy(func() exp.Expression {
			fmt.Printf("%v: lazily evaluating %s -> %s\n", done, name, e)
			done = true
			return e
		})
	}
	clause := two(func(p, e exp.Expression) exp.Expression {
		return exp.NewList(
			lazy("p", p),
			lazy("e", e),
		)
	})
	q := func(s string) exp.Expression {
		return exp.NewString(s)
	}
	t, f := True(), False()
	first := clause(f, q("first"))
	second := clause(f, q("second"))
	third := clause(t, q("third"))
	result := apply(Cond, first, second, third)
	fmt.Printf("result = %s\n", result)
	return result.Error()
}

// TODO: compile this from lisp source
func Eval(e, a exp.Expression) exp.Expression {

	reg("eval")

	//fmt.Printf("eval(%q, %q)\n", e, a)

	list := func(args ...exp.Expression) exp.Expression {
		return exp.NewList(args...)
	}

	if e.Error() != nil {
		return e
	}

	if IsAtom(e) {
		return Assoc(e, a)
	}

	evlis := two(Evlis)
	atom := one(Atom)
	eval := two(Eval)
	eq := two(Eq)
	cons := two(Cons)

	evcon := two(Evcon)
	assoc := two(Assoc)

	q := func(s string) exp.Expression {
		return exp.NewString(s)
	}

	is := func(s string) exp.Expression {
		return apply(eq, apply(car, e), q(s))
	}

	cxxr := func(f string) exp.Expression {
		return exp.NewList( // 'cdr
			exp.NewLazy(func() exp.Expression {
				return is(f)
			}),
			exp.NewLazy(func() exp.Expression {
				return apply(x(f), apply(eval, apply(cadr, e), a))
			}),
		)
	}

	return apply(Cond,
		exp.NewList( // ((atom e) (assoc e a))
			exp.NewLazy(func() exp.Expression {
				return apply(atom, e)
			}),
			exp.NewLazy(func() exp.Expression {
				return apply(assoc, e, a)
			}),
		),
		exp.NewList( // ((atom (car e))
			exp.NewLazy(func() exp.Expression {
				return apply(atom, apply(car, e))
			}),
			exp.NewLazy(func() exp.Expression {
				return apply(Cond,
					exp.NewList( // 'quote
						exp.NewLazy(func() exp.Expression {
							return is("quote")
						}),
						exp.NewLazy(func() exp.Expression {
							return apply(cadr, e)
						}),
					),
					exp.NewList( // 'atom
						exp.NewLazy(func() exp.Expression {
							return is("atom")
						}),
						exp.NewLazy(func() exp.Expression {
							return apply(atom, apply(eval, apply(cadr, e), a))
						}),
					),
					exp.NewList( // 'eq
						exp.NewLazy(func() exp.Expression {
							return is("eq")
						}),
						exp.NewLazy(func() exp.Expression {
							return apply(eq,
								apply(eval, apply(cadr, e), a),
								apply(eval, apply(caddr, e), a))
						}),
					),
					cxxr("car"),
					cxxr("cdr"),
					cxxr("cadr"),
					cxxr("caddr"),
					cxxr("cdar"),
					exp.NewList( // 'list
						exp.NewLazy(func() exp.Expression {
							return is("list")
						}),
						exp.NewLazy(func() exp.Expression {
							return apply(evlis, apply(cdr, e), a)
						}),
					),
					exp.NewList( // 'cons
						exp.NewLazy(func() exp.Expression {
							return is("cons")
						}),
						exp.NewLazy(func() exp.Expression {
							return apply(cons, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
						}),
					),
					exp.NewList( // 'cond
						exp.NewLazy(func() exp.Expression {
							return is("cond")
						}),
						exp.NewLazy(func() exp.Expression {
							return apply(evcon, apply(cdr, e), a)
						}),
					),
					exp.NewList( // 't
						exp.NewLazy(func() exp.Expression {
							return True()
						}),
						exp.NewLazy(func() exp.Expression {
							return apply(eval, apply(cons, apply(assoc, apply(car, e), a), apply(cdr, e)), a)
						}),
					),
				)
			}),
		),
		exp.NewList( // ((eq (caar e) 'label))
			exp.NewLazy(func() exp.Expression {
				return apply(eq, apply(caar, e), q("label"))
			}),
			exp.NewLazy(func() exp.Expression {
				return apply(eval,
					apply(cons, apply(caddar, e), apply(cdr, e)),
					apply(cons, apply(list, apply(cadar, e), apply(car, e)), a))
			}),
		),
		exp.NewList( // ((eq (caar e) 'lambda))
			exp.NewLazy(func() exp.Expression {
				return apply(eq, apply(caar, e), q("lambda"))
			}),
			exp.NewLazy(func() exp.Expression {
				return apply(eval,
					apply(caddar, e),
					apply(two(Append), apply(two(Pair), apply(cadar, e), apply(evlis, apply(cdr, e), a)), a))
			}),
		),
	)

	if x := car(e); IsAtom(x) {
		if f, err := cxr(x.String()); err == nil {
			return f(Eval(cadr(e), a))
		}
		switch x.String() {
		case "list":
			return apply(evlis, apply(cdr, e), a)
		case "quote":
			return apply(cadr, e)
		case "atom":
			return apply(atom, apply(eval, apply(cadr, e), a))
		case "eq":
			return apply(eq, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
		case "cons":
			return apply(cons, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
		case "cond":
			return apply(evcon, apply(cdr, e), a)
		default:
			return apply(eval, apply(cons, apply(assoc, apply(car, e), a), apply(cdr, e)), a)
		}
	}

	if x := caar(e); x.String() == "label" {
		return Eval(
			Cons(caddar(e), cdr(e)),
			Cons(list(cadar(e), car(e)), a),
		)
	}

	if x := caar(e); x.String() == "lambda" {
		e2 := caddar(e)
		a2 := Append(Pair(cadar(e), Evlis(cdr(e), a)),
			a,
		)
		return Eval(
			e2, a2,
		)
	}

	return exp.NewError(fmt.Errorf("eval can't handle (%s %s)", e, a))
}

func Pair(x, y exp.Expression) exp.Expression {
	reg("pair")
	if x.String() == "()" && y.String() == "()" {
		return exp.NewList()
	}
	if !IsAtom(x) && !IsAtom(y) {
		carx := Car(x)
		cary := Car(y)
		cdrx := Cdr(x)
		cdry := Cdr(y)
		list := exp.NewList(carx, cary)
		pair := Pair(cdrx, cdry)
		return Cons(list, pair)
	}
	return exp.NewError(fmt.Errorf("illegal pair state"))
}

type TwoFunc func(a, b exp.Expression) exp.Expression

func w2(name string, f TwoFunc) TwoFunc {
	return func(a, b exp.Expression) exp.Expression {
		out := f(a, b)
		fmt.Printf("%s(%q, %q) = %q\n", name, a, b, out)
		return out
	}
}

func Evlis(m, a exp.Expression) exp.Expression {
	reg("evlis")
	if m.String() == "()" {
		return exp.NewList()
	}
	car := Car(m)
	cdr := Cdr(m)
	eval := Eval(car, a)
	evlis := Evlis(cdr, a)
	return Cons(eval, evlis)
}

func Append(x, y exp.Expression) exp.Expression {
	reg("append")
	if x.String() == "()" {
		return y
	}
	car := Car(x)
	cdr := Cdr(x)
	tail := Append(cdr, y)
	return Cons(car, tail)
}

func Cons(x, y exp.Expression) exp.Expression {
	reg("cons")
	if !IsList(y) {
		return exp.NewError(fmt.Errorf("second arg not a list: %s", y))
	}
	var args []exp.Expression
	add := func(e exp.Expression) {
		args = append(args, e)
	}
	add(x)
	list := y.List()
	for _, e := range list {
		add(e)
	}
	return exp.NewList(args...)
}

func Evcon(c, a exp.Expression) exp.Expression {

	if c == nil || a == nil {
		return exp.NewError(fmt.Errorf("nil arguments"))
	}
	if !(IsList(c) && IsList(a)) {
		return exp.NewError(fmt.Errorf("needs lists"))
	}
	list := c.List()
	for _, arg := range list {
		car := Car(arg)
		r := Eval(car, a)
		if r.String() == "t" {
			cdr := Cdr(arg)
			cadr := Car(cdr)
			return Eval(cadr, a)
		}
	}
	return exp.NewError(fmt.Errorf("no condition satisfied"))
}

func Read(s string) exp.Expression {
	n, err := parse(s)
	if err != nil {
		return exp.NewError(err)
	}
	return n.Expression()
}

func Car(e exp.Expression) exp.Expression {
	if !IsList(e) {
		return exp.Errorf("can only car a list: %q", e)
	}
	list := e.List()
	if len(list) == 0 {
		return Nil()
	}
	return list[0]
}

func Cdr(e exp.Expression) exp.Expression {
	if !IsList(e) {
		return exp.Errorf("can only cdr a list: %q", e)
	}
	list := e.List()
	if len(list) == 0 {
		return Nil()
	}
	return exp.NewList(list[1:]...)
}

func Nil() exp.Expression {
	return exp.NewList()
}

func True() exp.Expression {
	return exp.NewString("t")
}

func False() exp.Expression {
	return Nil()
}

func Boolean(e exp.Expression) bool {
	return e.String() == "t"
}

func Atom(e exp.Expression) exp.Expression {
	reg("atom")
	if IsAtom(e) {
		return exp.NewString("t")
	}
	return Nil()
}

func IsAtom(e exp.Expression) bool {
	reg("isatom")
	if e.Atom() != nil {
		return true
	}
	if len(e.List()) == 0 {
		return true
	}
	return false
}

func IsList(e exp.Expression) bool {
	return e.Atom() == nil
}

func IsEmpty(e exp.Expression) bool {
	return len(e.List()) == 0
}

func AtomsEqual(a, b exp.Atom) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil:
		return false
	case b == nil:
		return false
	}
	return a.String() == b.String()
}
