package scr

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/scr/exp"
)

func RunTests() error {

	type test struct {
		name string
		f    func() error
	}
	var tests []test

	single := func(file string) {
		tests = append(tests, test{
			name: file,
			f: func() error {
				return RunTest(file)
			},
		})
	}

	if true {
		if err := RunTest("tests/append.lisp"); err != nil {
			return err
		}
		if err := RunTest("tests/append_label.lisp"); err != nil {
			return err
		}
	} else {
		files, err := loadLisp("tests")
		if err != nil {
			return err
		}
		for _, file := range files {
			single(file)
		}
	}

	for _, t := range tests {
		if err := t.f(); err != nil {
			return fmt.Errorf("can't run %q: %w", t.name, err)
		}
	}
	return nil
}

func debugging() error {
	buf, err := ioutil.ReadFile("test.lisp")
	if err != nil {
		return err
	}
	e, err := Read(string(buf))
	if err != nil {
		return err
	}

	r := Eval(e, Nil())
	fmt.Println(r)

	return nil

}

func loadLisp(dir string) ([]string, error) {
	type file struct {
		name string
		size int
	}
	var files []file
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	hashes := make(map[string]bool)
	hash := func(b []byte) string {
		h := md5.New()
		h.Write(b)
		return fmt.Sprintf("%x", h.Sum(nil))
	}
	space := regexp.MustCompile(`\s+`)
	for _, fi := range list {
		name := filepath.Join(dir, fi.Name())
		if filepath.Ext(name) != ".lisp" {
			continue
		}
		buf, err := ioutil.ReadFile(name)
		if err != nil {
			return nil, err
		}
		if hashes[hash(buf)] {
			fmt.Printf("removing %q\n", name)
			os.Remove(name)
			continue
		}
		hashes[hash(buf)] = true
		content := space.ReplaceAllString(string(buf), " ")
		files = append(files, file{
			name: name,
			size: len(content),
		})
	}
	// start with smaller files:
	sort.Slice(files, func(i, j int) bool {
		return files[i].size < files[j].size
	})
	var out []string
	for _, f := range files {
		out = append(out, f.name)
	}
	return out, nil
}

func LoadEnv(files ...string) (exp.Expression, error) {
	a := Nil()
	for _, file := range files {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		def, err := Read(string(buf))
		if err != nil {
			return nil, err
		}
		name := Car(Cdr(def))
		args := Car(Cdr(Cdr(def)))
		body := Car(Cdr(Cdr(Cdr(def))))
		label := exp.NewList(
			q("label"),
			name,
			exp.NewList(
				q("lambda"),
				args,
				body,
			),
		)
		env := exp.NewList(name, label)
		var list []exp.Expression
		list = append(list, env)
		list = append(list, a.List()...)
		a = exp.NewList(list...)
	}
	return a, nil
}

func LoadAllDefs() (exp.Expression, error) {
	files, err := loadLisp("defs")
	if err != nil {
		return nil, err
	}
	return LoadEnv(files...)
}

func singleTest(file string, env exp.Expression) error {
	fmt.Printf("run(%q)\n", file)
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	test, err := Read(string(buf))
	if err != nil {
		return err
	}
	in := Car(test)
	out := Car(Cdr(test))
	fmt.Printf("%-16s: %s -> %s\n", filepath.Base(file), in, out)
	e := Eval(in, env)
	if e.String() != out.String() {
		return fmt.Errorf("%s failed; expected %s, got %s", file, out, e)
	}
	return nil
}

func RunTest(file string) error {
	a, err := LoadAllDefs()
	if err != nil {
		return err
	}
	return singleTest(file, a)
}

func Lisp(config cnfg.Config) error {

	return RunTests()

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

	names := make(map[string]bool)
	test := func(i interface{}, in, expect string) {

		if in == "" {
			return
		}
		if m[in] {
			panic("duplication: " + in)
		}
		count++
		m[in] = true

		var name string
		switch t := i.(type) {
		case int:
			name = fmt.Sprintf("%03d-%03d", t, count)
		case string:
			name = fmt.Sprintf("%s-%03d", t, count)
		default:
			panic(t)
		}

		if names[name] {
			panic("dup name " + name)
		}
		names[name] = true

		f, err := os.Create(fmt.Sprintf("tests/%s.lisp", name))
		check(err)
		defer f.Close()
		fmt.Fprintf(f, fmt.Sprintf(`(%s
%s
)
`, in, expect))

		return

		e, err := Read(in)
		check(wrap(i, err))

		ex, err := Read(expect)
		check(wrap(i, err))

		fmt.Printf("%2d.%v: %-55s -> %s\n", count, i, e, ex)
		check(wrap(i, e.Error()))
		x := Eval(e, a)
		check(wrap(i, x.Error()))
		if got := x.String(); got != expect {
			check(wrap(i, fmt.Errorf("expected %q, got %q", expect, got)))
		}
	}

	define := func(name, lambda string) {
		return
		if name == "" {
			return
		}
		lambdaE, err := Read(lambda)
		check(err)

		e, err := Read(fmt.Sprintf("(label %s %s)", name, lambda))
		check(err)

		check(wrap(name, e.Error()))
		a = exp.NewList(exp.NewList(exp.NewString(name), e), a)

		args := Car(Cdr(lambdaE))
		body := Car(Cdr(Cdr(lambdaE)))
		f, err := os.Create(fmt.Sprintf("defs/%s.lisp", name))
		check(err)
		defer f.Close()
		fmt.Fprintf(f, fmt.Sprintf(`(defun %s %s 
%s)
`, name, args, body))
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
	}

	if true {
		define("not", "(lambda (x) (cond (x '()) ('t 't)))")

		define("and", `

(lambda (x y)
   (cond (x (cond (y 't) ('t ())))
	 ('t '())))
`)

		define("pair", `
(lambda (x y) (cond ((and (null x) (null y)) '())
		    ((and (not (atom x)) (not (atom y)))
		     (cons (list (car x) (car y))
			   (pair (cdr x) (cdr y))))))`)
		test("funcs", `(pair '(x y z) '(a b c))`, "((x a) (y b) (z c))")
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

	test("funcs", "(and (atom 'a) (eq 'a 'a))", "t")
	test("funcs", "(and (atom 'a) 't)", "t")
	test("funcs", "(and (atom 'a) '())", "()")

	test("funcs", `(not (eq 'a 'a))`, "()")
	test("funcs", `(not (eq 'a 'b))`, "t")

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

func Assoc(x, y exp.Expression) exp.Expression {
	eq := two(Eq)
	assoc := two(Assoc)
	return apply(Cond,
		exp.NewList(
			exp.NewLazy(func() exp.Expression {
				return apply(eq, apply(caar, y), x)
			}),
			exp.NewLazy(func() exp.Expression {
				return apply(cadar, y)
			}),
		),
		exp.NewList(
			True(),
			exp.NewLazy(func() exp.Expression {
				return apply(assoc, x, apply(cdr, y))
			}),
		),
	)
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

	q = func(s string) exp.Expression {
		return exp.NewString(s)
	}
)

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
	t, f := True(), False()
	first := clause(f, q("first"))
	second := clause(f, q("second"))
	third := clause(t, q("third"))
	result := apply(Cond, first, second, third)
	fmt.Printf("result = %s\n", result)
	return result.Error()
}

func List(x, y exp.Expression) exp.Expression {
	cons := two(Cons)
	return apply(cons, x, apply(cons, y, Nil()))
}

// TODO: compile this from lisp source
func Eval(e, a exp.Expression) exp.Expression {

	debug("EVAL(%q, %q)", e, a)

	list := two(List)

	if e.Error() != nil {
		return e
	}
	if a.Error() != nil {
		return a
	}

	evlis := two(Evlis)
	atom := one(Atom)
	eval := two(Eval)
	eq := two(Eq)
	cons := two(Cons)
	appnd := two(Append)
	pair := two(Pair)
	evcon := two(Evcon)
	assoc := two(Assoc)

	nl := exp.NewList
	z := exp.NewLazy
	y := apply

	is2 := func(s string) exp.Expression {
		is := func(s string) exp.Expression {
			return apply(eq, apply(car, e), q(s))
		}
		return z(func() exp.Expression {
			return is(s)
		})
	}

	is3 := func(s string) exp.Expression {
		return z(func() exp.Expression {
			return y(eq, y(caar, e), q(s))
		})
	}

	cxxr := func(f string) exp.Expression {
		return nl(
			is2(f),
			z(func() exp.Expression {
				return y(x(f), y(eval, y(cadr, e), a))
			}),
		)
	}

	return y(Cond,
		nl(
			y(atom, e),
			z(func() exp.Expression {
				return y(assoc, e, a)
			}),
		),
		nl(
			z(func() exp.Expression {
				return y(atom, y(car, e))
			}),
			z(func() exp.Expression {
				return y(Cond,
					nl(
						is2("quote"),
						z(func() exp.Expression {
							return y(cadr, e)
						}),
					),
					nl(
						is2("atom"),
						z(func() exp.Expression {
							return y(atom, y(eval, y(cadr, e), a))
						}),
					),
					nl(
						is2("eq"),
						z(func() exp.Expression {
							return y(eq,
								y(eval, y(cadr, e), a),
								y(eval, y(caddr, e), a))
						}),
					),
					cxxr("car"),
					cxxr("cdr"),
					nl(
						is2("list"),
						z(func() exp.Expression {
							return y(evlis, y(cdr, e), a)
						}),
					),
					nl(
						is2("cons"),
						z(func() exp.Expression {
							return y(cons,
								y(eval, y(cadr, e), a),
								y(eval, y(caddr, e), a))
						}),
					),
					nl(
						is2("cond"),
						z(func() exp.Expression {
							return y(evcon, y(cdr, e), a)
						}),
					),
					nl(
						True(),
						z(func() exp.Expression {
							return y(eval,
								y(cons,
									y(assoc, y(car, e), a),
									y(cdr, e)),
								a)
						}),
					),
				)
			}),
		),
		nl(
			is3("label"),
			z(func() exp.Expression {
				return y(eval,
					y(cons, y(caddar, e), y(cdr, e)),
					y(cons, y(list, y(cadar, e), y(car, e)), a))
			}),
		),
		nl(
			is3("lambda"),
			z(func() exp.Expression {
				return y(eval,
					y(caddar, e),
					y(appnd, y(pair, y(cadar, e), y(evlis, y(cdr, e), a)),
						a))
			}),
		),
	)

	return exp.Errorf("eval can't handle (%s %s)", e, a)
}

func IsAtom(e exp.Expression) bool {
	if e.Atom() != nil {
		return true
	}
	if len(e.List()) == 0 {
		return true
	}
	return false
}

func And(x, y exp.Expression) exp.Expression {
	return apply(Cond,
		exp.NewList(
			x,
			apply(Cond,
				exp.NewList(
					y,
					True(),
				),
			),
		),
		exp.NewList(
			True(),
			Nil(),
		),
	)
}

func Not(x exp.Expression) exp.Expression {
	return apply(Cond,
		exp.NewList(x, Nil()),
		exp.NewList(True(), True()),
	)
}

func debug(format string, args ...interface{}) {
	x := fmt.Sprintf(format, args...)
	fmt.Println(x)
	return

	hash := md5.New()
	hash.Write([]byte(x))
	fmt.Printf("%x: %s\n", hash.Sum(nil)[:2], x)
}

func Pair(x, y exp.Expression) exp.Expression {
	debug("PAIR(%q, %q)", x, y)
	if false {
		and := two(And)
		not := one(Not)
		atom := one(Atom)
		null := one(Null)
		cons := two(Cons)
		list := two(List)
		pair := two(Pair)
		return apply(Cond,
			exp.NewList(
				apply(and, apply(null, x), apply(null, y)),
				Nil(),
			),
			exp.NewList(
				exp.NewLazy(func() exp.Expression {
					return apply(and,
						apply(not, apply(atom, x)),
						apply(not, apply(atom, y)))
				}),
				exp.NewLazy(func() exp.Expression {
					return apply(cons,
						apply(list, apply(car, x), apply(car, y)),
						apply(pair, apply(cdr, x), apply(cdr, y)))

				}),
			),
			exp.NewList(
				True(),
				exp.Errorf("illegal"),
			),
		)
	}

	if x.String() == "()" && y.String() == "()" {
		return exp.NewList()
	}
	if !IsAtom(x) && !IsAtom(y) {
		return Cons(exp.NewList(Car(x), Car(y)), Pair(Cdr(x), Cdr(y)))
	}
	return exp.Errorf("illegal pair state")
}

type TwoFunc func(a, b exp.Expression) exp.Expression

func w2(name string, f TwoFunc) TwoFunc {
	return func(a, b exp.Expression) exp.Expression {
		out := f(a, b)
		fmt.Printf("%s(%q, %q) = %q\n", name, a, b, out)
		return out
	}
}

func Append(x, y exp.Expression) exp.Expression {
	null := one(Null)
	cons := two(Cons)
	append := two(Append)
	return apply(Cond,
		exp.NewList(
			apply(null, x),
			y,
		),
		exp.NewList(
			True(),
			exp.NewLazy(func() exp.Expression {
				return apply(cons,
					apply(car, x),
					apply(append, apply(cdr, x), y),
				)
			}),
		),
	)

	if x.String() == "()" {
		return y
	}
	car := Car(x)
	cdr := Cdr(x)
	tail := Append(cdr, y)
	return Cons(car, tail)
}

func Evlis(m, a exp.Expression) exp.Expression {
	null := one(Null)
	cons := two(Cons)
	eval := two(Eval)
	evlis := two(Evlis)
	return apply(Cond,
		exp.NewList(
			exp.NewLazy(func() exp.Expression {
				return apply(null, m)
			}),
			Nil(),
		),
		exp.NewList(
			True(),
			exp.NewLazy(func() exp.Expression {
				return apply(
					cons,
					apply(eval, apply(car, m), a),
					apply(evlis, apply(cdr, m), a),
				)
			}),
		),
		exp.NewList(),
	)
}

func Null(x exp.Expression) exp.Expression {
	eq := two(Eq)
	return apply(eq, x, Nil())
}

func Evcon(c, a exp.Expression) exp.Expression {
	eval := two(Eval)
	evcon := two(Evcon)
	return apply(Cond,
		exp.NewList(
			exp.NewLazy(func() exp.Expression {
				return apply(eval, apply(caar, c), a)
			}),
			exp.NewLazy(func() exp.Expression {
				return apply(eval, apply(cadar, c), a)
			}),
		),
		exp.NewList(
			True(),
			exp.NewLazy(func() exp.Expression {
				return apply(evcon, apply(cdr, c), a)
			}),
		),
	)
	return exp.Errorf("evcon fallthrough")
}

func Read(s string) (exp.Expression, error) {
	n, err := parse(s)
	if err != nil {
		return nil, err
	}
	return n.Expression()
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
	return e.String() == True().String()
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
