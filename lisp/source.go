// compiled lisp stuff
package lisp

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/xoba/turd/cnfg"
)

func Eval(e Exp) Exp {
	e = SanitizeGo(e)
	return UnsanitizeGo(eval([]Exp{e, env}...))
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
	if c.Lisp != "" {
		return test0(c.Lisp)
	}

	test("a b")
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

func Run(cnfg.Config) error {
	var last string
	test := func(msg, input, expect string) {
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
		fmt.Printf("%-10s %-20s -> %s\n", msg+":", String(in), expect)
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
		res := Eval(in)
		if expect == "" {
			fmt.Printf("got %s\n", String(res))
			return
		}
		if got := String(res); got != expect {
			log.Fatalf("expected %q, got %q\n", expect, got)
		}
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

	test("display", "(display ())", "()")
	test("display", "(display 'a)", "a")

	if false {
		test("macro", `((macro test (x) (cdr x))
 'a 'b 'c)`, "(a b c)")
		test("printf", "(printf 'a)", "()")
	}

	test("arith", "(plus '1 '2)", "3")
	test("arith", "(plus '1 '-2)", "-1")
	test("arith", "(minus '1 '2)", "-1")
	test("arith", "(minus '1 '-2)", "3")
	test("arith", "(mult '4 '5)", "20")
	test("arith", "(mult '4 '-2)", "-8")
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

	test("crypto", "(hash 'MHcCAQEEIBTMVA5sze5UsF4PMb5xNKndc7YKVIg5AbyjoBWiPWnfoAoGCCqGSM49AwEHoUQDQgAECu2rhbYGyKK5wT5zjgFbVlCzMQZe6LeinrX0xvw+dt7qRRTJERKvKiCuTmLKt/O3SAZGVnozSGoBCGuKaSx/mw==)", "")
	test("crypto", "(hash (newkey))", "")
	test("crypto", `((lambda (private content)
   (verify (pub private) (hash content) (sign private (hash content)))
    ) (newkey) 'c2RmZgo=)
`, "t")
	test("crypto", `((lambda (private content)
   (verify (pub private) (hash content) (sign private (hash 'MTIzCg==)))
    ) (newkey) 'c2RmZgo=)
`, "()")
	test("", "", "")
	test("", "", "")
	test("", "", "")

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
			return fmt.Errorf("illegal exp type: %T %v", t, t)
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
