// compiled lisp stuff
package lisp

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/xoba/turd/cnfg"
)

type EvalFunc func(e Exp) Exp

func Eval(e Exp) Exp {
	if true {
		return CompiledEval(e)
	} else {
		return InterpretedEval(e, CompiledEval)
	}
}

func CompiledEval(e Exp) Exp {
	return UnsanitizeGo(eval([]Exp{SanitizeGo(e), env}...))
}

func InterpretedEval(e Exp, eval EvalFunc) Exp {
	q := func(e Exp) Exp {
		return []Exp{"quote", e}
	}
	return eval([]Exp{eval_label, q(e), q(env)})
}

func TestParse(c cnfg.Config) error {
	test0 := func(s string) error {
		fmt.Printf("testing %s\n", s)
		exp, err := Parse(s)
		if err != nil {
			return err
		}
		fmt.Printf("  -> %v\n", exp)
		return nil
	}
	test := func(s string) {
		if err := test0(s); err != nil {
			log.Fatal(err)
		}
	}
	test2 := func(s string) {
		if err := test0(s); err == nil {
			log.Fatalf("expected error for %q", s)
		}
	}
	if c.Lisp != "" {
		return test0(c.Lisp)
	}

	test2("a b")
	test("(a b)")
	test("(a b c)")
	test("(a (x y) b c)")
	test("(a (x (123) y) b c)")
	test("(a (x (123 z) y) b c)")

	test("'a")
	test("'(a b c)")
	test("'(a b 'c)")
	test("'(a \"this is a test\" 'c)")

	return nil
}

func Run(c cnfg.Config) error {
	var last string

	tests := make(map[string]int)

	test0 := func(msg, input, expect string, f EvalFunc, name string) {
		if msg == "" {
			return
		}
		if msg != last {
			fmt.Println()
		}
		last = msg
		in, err := Parse(input)
		if err != nil {
			log.Fatal(err)
		}
		tests[String(in)]++
		log.Printf("%3d. %-13s %-10s %-20s -> %s", len(tests), name, msg+":", String(in), expect)
		{
			buf, err := Marshal(in)
			if err != nil {
				log.Fatal(err)
			}
			if len(buf) < len(input) { // is asn.1 more compact?
				fmt.Printf("%d/%d asn1 bytes = %s\n",
					len(buf), len(input),
					base64.StdEncoding.EncodeToString(buf),
				)
			}
		}
		res := f(in)
		if expect == "" {
			fmt.Printf("got %s\n", String(res))
			return
		}
		if got := String(res); got != expect {
			log.Fatalf("expected %q, got %q\n", expect, got)
		}
	}

	evals := map[string]EvalFunc{
		"compiled": CompiledEval,
		"interpreted": func(e Exp) Exp {
			return InterpretedEval(e, CompiledEval)
		},
		"interpreted2": func(e Exp) Exp {
			return InterpretedEval(e, func(e Exp) Exp {
				return InterpretedEval(e, CompiledEval)
			})
		},
	}

	// interpreted2 is way too slow! never actually saw it finish:
	delete(evals, "interpreted2")

	if !c.Debug {
		delete(evals, "interpreted2")
		delete(evals, "interpreted")
	}

	test := func(msg, input, expect string) {
		for k, f := range evals {
			test0(msg, input, expect, f, k)
		}
	}

	file := func(f, expect string) {
		buf, err := ioutil.ReadFile(filepath.Join("lisp", "tests", f))
		if err != nil {
			log.Fatal(err)
		}
		test(f, string(buf), expect)
	}

	test("quote1", "(quote a)", "a")
	test("quote1", "(quote (a b c))", "(a b c)")

	test("quote2", "'a", "a")
	test("quote2", "'(a b c)", "(a b c)")

	test("atom", "(atom 'a)", "t")
	test("atom", "(atom '(a b c))", "()")
	test("atom", "(atom '())", "t")
	test("atom", "(atom 'a)", "t")
	test("atom", "(atom '(atom 'a))", "()")

	test("eq", "(eq 'a 'a)", "t")
	test("eq", "(eq 'a 'b)", "()")
	test("eq", "(eq '() '())", "t")

	test("cxr", "(car '())", "()")
	test("cxr", "(cdr '())", "()")
	test("cxr", "(car '(a))", "a")
	test("cxr", "(cdr '(a))", "()")
	test("cxr", "(car '(a b))", "a")
	test("cxr", "(cdr '(a b))", "(b)")
	test("cxr", "(car '(a b c))", "a")
	test("cxr", "(cdr '(a b c))", "(b c)")

	test("cons", "(cons 'a '())", "(a)")
	test("cons", "(cons '() '())", "(())")
	test("cons", "(cons 'a '(b))", "(a b)")
	test("cons", "(cons '() '(b))", "(() b)")
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
	test("cxr", "(caddr '(a b c d e))", "c")
	test("cxr", "(cadddr '(a b c d e))", "d")
	test("cxr", "(caddddr '(a b c d e))", "e")

	test("list", "(cons 'a (cons 'b (cons 'c '())))", "(a b c)")
	test("list", "(list 'a 'b 'c)", "(a b c)")
	test("list", "(car (list 'a 'b 'c))", "a")
	test("list", "(car (cons 'a '(b c)))", "a")
	test("list", "(cdr (cons 'a '(b c)))", "(b c)")

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

	test("display", "(display ())", "()")
	test("display", "(display 'a)", "a")

	if false {
		test("macro", `((macro test (x) (cdr x))
 'a 'b 'c)`, "(a b c)")
		test("printf", "(printf 'a)", "()")
	}

	test("arith", "(plus '1 '2)", "3")
	test("arith", "(plus '1 '-2)", "-1")
	test("arith", "(inc '4)", "5")
	test("arith", "(inc '5)", "6")
	test("arith", "(minus '1 '2)", "-1")
	test("arith", "(minus '1 '-2)", "3")
	test("arith", "(mult '4 '5)", "20")
	test("arith", "(mult '4 '-2)", "-8")
	test("arith", "(mult '234862342873462784637846 '104380123947329857341285)",
		"24514960459692328665578339488420194423149272110",
	)
	test("arith", "(eq '0 (minus '5 '5))", "t")
	test("arith", "(eq '1 (minus '5 '5))", "()")

	test("factorial", "(factorial '0)", "1")
	test("factorial", "(factorial '1)", "1")
	test("factorial", "(factorial '3)", "6")
	test("factorial", "(factorial '10)", "3628800")
	test("factorial", "(factorial '100)", "93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000")

	test("lexpr", " ((lambda (x) (cdr x)) '(a b c))", "(b c)")
	test("lexpr", " ((lambda x (cdr x)) 'a 'b 'c)", "(b c)")
	test("lexpr", "((lambda x x) 'a 'b 'c)", "(a b c)")
	test("lexpr", "(list 'a 'b 'c)", "(a b c)")

	test("list", "(xlist 'a 'b 'c)", "(a b c)")

	test("crypto", "(newkey)", "")

	test("crypto", "(eq (hash 'RSrYRpagDHgCuQ) 'ChUPiFeRkYRliIqlN8CB4Vfjce4/zHoEN9wBRKr2MKY)", "t")
	test("crypto", "(eq (hash 'RSrYRpagDHgCuQ) '0000iFeRkYRliIqlN8CB4Vfjce4/zHoEN9wBRKr2MKY)", "()")
	test("crypto", "(hash (newkey))", "")
	test("crypto", `((lambda (priv content)
   (verify (pub priv) (hash content) (sign priv (hash content)))
    ) (newkey) 'c2RmZgo)
`, "t")
	// content mismatch:
	test("crypto", `((lambda (priv content)
   (verify (pub priv) (hash content) (sign priv (hash 'MTIzCg)))
    ) (newkey) 'c2RmZgo)
`, "()")
	// public key mismatch
	test("crypto", `((lambda (priv content)
   (verify (pub (newkey)) (hash content) (sign priv (hash content)))
    ) (newkey) 'c2RmZgo)
`, "()")

	test("blobs", "(concat 'YWJj 'eHl6)", "YWJjeHl6")

	test("crypto", `(assoc 'hash '((hash ymWy8GC+4f6xaqTYlciUS0+FBq+zO7XRe46fYtMni5Y) (IgSk3rkEcxbx1f44G5diGvc53pwSi6WjlsZnbWpHgWk MEYCIQC1Fnd+LuN4AFJ7lYWBVjEcjO7SvrTAoUtcUct96za1OQIhAKIYsyE/rroV4GuNuNJNGFIcDMsL27VlZgXNcIZQk81Z)))`, "")

	test("lambda", `
((lambda (x y)
   ((lambda (z) (plus z x)) y))
 '3 '4)
`, "7")

	test("time", `(after '2020-11-20T10:00:00.000Z '2020-11-21T10:00:00.000Z)`, `()`)
	test("time", `(after '2020-11-21T10:00:00.000Z '2020-11-21T10:00:00.000Z)`, `()`)
	test("time", `(after '2020-11-22T10:00:00.000Z '2020-11-21T10:00:00.000Z)`, `t`)

	file("crypto.lisp", "YxxZaKeUU5dgl5KQ/emPanGk8OsM/r9XfIjwpGv1MBc")
	file("crypto2.lisp", "/Bmm7Jh2vHGd56htMfK1looDkZfay3hJgLrY+QZRQvM")
	file("block.lisp", "1000")
	file("block1.lisp", "2020-11-22T11:16:18.838Z")
	file("trans.lisp", "")
	file("block2.lisp", "rwSvYFupU7ZH1xqBj/YFOzaHZsyeSye0TGErmVV3mAI")

	test("len", `(length 'x)`, `0`)
	test("len", `(length '())`, `0`)
	test("len", `(length '(a))`, `1`)
	test("len", `(length '(a b c))`, `3`)
	test("len", `(length '(a (1 2 3) c))`, `3`)

	// first-order functions
	test("fof", `((lambda (f) (f '(a))) 'car)`, `a`)

	if false {
		test("iscxr", `(iscxr 'car)`, `()`) // car is axiom
		test("iscxr", `(iscxr 'cdr)`, `()`) // cdr is axiom
		test("iscxr", `(iscxr 'caar)`, `t`)
		test("iscxr", `(iscxr 'cadr)`, `t`)
		test("iscxr", `(iscxr 'cdar)`, `t`)
		test("iscxr", `(iscxr 'cddr)`, `t`)
		test("iscxr", `(iscxr 'caaar)`, `t`)
		test("iscxr", `(iscxr 'caaaar)`, `t`)
		test("iscxr", `(iscxr 'dfdf)`, `()`)
		test("iscxr", `(iscxr '123)`, `()`)
		test("iscxr", `(iscxr 'caxr)`, `()`)
		test("iscxr", `(iscxr 'cxdr)`, `()`)
		test("iscxr", `(iscxr 'cdxr)`, `()`)
		test("iscxr", `(iscxr 'caxaar)`, `()`)
		test("iscxr", `(iscxr 'a)`, `()`)
		test("iscxr", `(iscxr 'd)`, `()`)
		test("iscxr", `(iscxr 'c)`, `()`)
		test("iscxr", `(iscxr 'r)`, `()`)
		test("iscxr", `(iscxr 'ar)`, `()`)
		test("iscxr", `(iscxr 'dr)`, `()`)
		test("iscxr", `(iscxr 'cr)`, `()`)
	}

	test("runes", `(runes 'abc)`, "(a b c)")
	test("runes", `(runes 'caddr)`, "(c a d d r)")

	test("errors", `(err 'abc)`, "error: abc")
	test("errors", `(err (list 'a 'b 'c))`, "error: (a b c)")

	test("try", "(try '(10 0) 'x '((x a) (y b)))", "a")
	test("try", "(try '(10 0) '(car '(a b)) '())", "a")
	test("try", "(try '(10 0) '(cdr '(a b)) '())", "(b)")
	test("try", fmt.Sprintf("(try '(10 0) '(cdr '(a b)) '%s)", String(env)), "(b)")
	test("try", "(try '(10 0) '(blah '(a b)) '())", "error: (max 10)")

	test("next", `(next '(10 6))`, `(10 7)`)
	test("next", `(next '(10 10))`, `error: (max 10)`)

	test("ltest", `((label lambdatest (lambda (x) 
  (list (car x) (cdr x)))) '(a b c))`, "(a (b c))")
	test("ltest", `((label lambdatest
	(lambda (x)
	  ((lambda (first rest) 
	     (list first rest)) (car x) (cdr x))))
 '(a b c))`, "(a (b c))")
	test("ltest", `(lambdatest '(a b c))`, `(a (b c))`)

	test("ltest", `(test '(a b c))`, `(a (b c))`)

	test("", ``, ``)
	test("", ``, ``)
	test("", ``, ``)
	test("", ``, ``)
	test("", ``, ``)
	test("", ``, ``)
	test("", ``, ``)
	test("", ``, ``)
	test("", ``, ``)

	return nil

	// how to handle infinite eval loop with unknown operator?
	test("eval", `(blah 'x)`, ``)

	return nil
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
			return fmt.Errorf("illegal type %T: %v", t, t)
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
		return fmt.Errorf("expected %d, got %d args", n, len(args))
	}
	if err := checkargs(args); err != nil {
		return err
	}
	return nil
}

func boolToExp(v bool) Exp {
	if v {
		return True
	}
	return False
}
