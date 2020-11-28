// THIS FILE IS AUTOGENERATED, DO NOT EDIT!

package lisp

import "fmt"

func init() {
	return
	fmt.Println("gen.go: THIS FILE IS AUTOGENERATED, DO NOT EDIT!")
}

var (
	L = list
	A = apply
)

func parse_env(s string) Exp {
	e, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return e
}

//
// and (compiled)
//

var and_label = parse_env("(label and (lambda (x y) (cond (x (cond (y 't) ('t ()))) ('t '()))))")

func and(args ...Exp) Exp {
	x := args[0]
	y := args[1]
	return A(
		cond,
		L(
			x,
			Func(func(...Exp) Exp {
				return A(cond, L(
					y,
					"t",
				), L(
					"t",
					Nil,
				))
			}),
		),
		L(
			"t",
			Nil,
		),
	)
}

//
// append (compiled)
//

var append_go_sanitized_label = parse_env("(label append_go_sanitized (lambda (x y) (cond ((null x) y) ('t (cons (car x) (append_go_sanitized (cdr x) y))))))")

func append_go_sanitized(args ...Exp) Exp {
	x := args[0]
	y := args[1]
	return A(
		cond,
		L(
			Func(func(...Exp) Exp {
				return A(null, x)
			}),
			y,
		),
		L(
			"t",
			Func(func(...Exp) Exp {
				return A(cons, A(car, x), A(append_go_sanitized, A(cdr, x), y))
			}),
		),
	)
}

//
// assoc (compiled)
//

var assoc_label = parse_env("(label assoc (lambda (x y) (cond ((eq (caar y) x) (cadar y)) ('t (assoc x (cdr y))))))")

func assoc(args ...Exp) Exp {
	x := args[0]
	y := args[1]
	return A(
		cond,
		L(
			Func(func(...Exp) Exp {
				return A(eq, A(caar, y), x)
			}),
			Func(func(...Exp) Exp {
				return A(cadar, y)
			}),
		),
		L(
			"t",
			Func(func(...Exp) Exp {
				return A(assoc, x, A(cdr, y))
			}),
		),
	)
}

//
// caar (compiled)
//

var caar_label = parse_env("(label caar (lambda (x) (car (car x))))")

func caar(args ...Exp) Exp {
	x := args[0]
	return A(
		car,
		A(
			car,
			x,
		),
	)
}

//
// cadar (compiled)
//

var cadar_label = parse_env("(label cadar (lambda (x) (car (cdr (car x)))))")

func cadar(args ...Exp) Exp {
	x := args[0]
	return A(
		car,
		A(
			cdr,
			A(
				car,
				x,
			),
		),
	)
}

//
// caddar (compiled)
//

var caddar_label = parse_env("(label caddar (lambda (x) (car (cdr (cdr (car x))))))")

func caddar(args ...Exp) Exp {
	x := args[0]
	return A(
		car,
		A(
			cdr,
			A(
				cdr,
				A(
					car,
					x,
				),
			),
		),
	)
}

//
// cadddar (compiled)
//

var cadddar_label = parse_env("(label cadddar (lambda (x) (car (cdr (cdr (cdr (car x)))))))")

func cadddar(args ...Exp) Exp {
	x := args[0]
	return A(
		car,
		A(
			cdr,
			A(
				cdr,
				A(
					cdr,
					A(
						car,
						x,
					),
				),
			),
		),
	)
}

//
// caddddar (compiled)
//

var caddddar_label = parse_env("(label caddddar (lambda (x) (car (cdr (cdr (cdr (cdr (car x))))))))")

func caddddar(args ...Exp) Exp {
	x := args[0]
	return A(
		car,
		A(
			cdr,
			A(
				cdr,
				A(
					cdr,
					A(
						cdr,
						A(
							car,
							x,
						),
					),
				),
			),
		),
	)
}

//
// caddddr (interpreted)
//

var caddddr_label = parse_env("(label caddddr (lambda (x) (car (cdr (cdr (cdr (cdr x)))))))")

//
// cadddr (compiled)
//

var cadddr_label = parse_env("(label cadddr (lambda (x) (car (cdr (cdr (cdr x))))))")

func cadddr(args ...Exp) Exp {
	x := args[0]
	return A(
		car,
		A(
			cdr,
			A(
				cdr,
				A(
					cdr,
					x,
				),
			),
		),
	)
}

//
// caddr (compiled)
//

var caddr_label = parse_env("(label caddr (lambda (x) (car (cdr (cdr x)))))")

func caddr(args ...Exp) Exp {
	x := args[0]
	return A(
		car,
		A(
			cdr,
			A(
				cdr,
				x,
			),
		),
	)
}

//
// cadr (compiled)
//

var cadr_label = parse_env("(label cadr (lambda (x) (car (cdr x))))")

func cadr(args ...Exp) Exp {
	x := args[0]
	return A(
		car,
		A(
			cdr,
			x,
		),
	)
}

//
// cdar (interpreted)
//

var cdar_label = parse_env("(label cdar (lambda (x) (cdr (car x))))")

//
// cdddar (compiled)
//

var cdddar_label = parse_env("(label cdddar (lambda (x) (cdr (cdr (cdr (car x))))))")

func cdddar(args ...Exp) Exp {
	x := args[0]
	return A(
		cdr,
		A(
			cdr,
			A(
				cdr,
				A(
					car,
					x,
				),
			),
		),
	)
}

//
// eval (compiled)
//

var eval_label = parse_env("(label eval (lambda (e a) (cond ((atom e) (assoc e a)) ((atom (car e)) (cond ((eq (car e) 'test) (test (eval (cadr e) a))) ((eq (car e) 'test2) (test2 (eval (cadr e) a))) ((eq (car e) 'quote) (cadr e)) ((eq (car e) 'atom) (atom (eval (cadr e) a))) ((eq (car e) 'eq) (eq (eval (cadr e) a) (eval (caddr e) a))) ((eq (car e) 'car) (car (eval (cadr e) a))) ((eq (car e) 'cdr) (cdr (eval (cadr e) a))) ((eq (car e) 'cons) (cons (eval (cadr e) a) (eval (caddr e) a))) ((eq (car e) 'cond) (evcon (cdr e) a)) ((eq (car e) 'plus) (plus (eval (cadr e) a) (eval (caddr e) a))) ((eq (car e) 'inc) (plus (eval (cadr e) a) '1)) ((eq (car e) 'minus) (minus (eval (cadr e) a) (eval (caddr e) a))) ((eq (car e) 'mult) (mult (eval (cadr e) a) (eval (caddr e) a))) ((eq (car e) 'exp) (exp (eval (cadr e) a) (eval (caddr e) a) (eval (cadddr e) a))) ((eq (car e) 'after) (after (eval (cadr e) a) (eval (caddr e) a))) ((eq (car e) 'concat) (concat (eval (cadr e) a) (eval (caddr e) a))) ((eq (car e) 'hash) (hash (eval (cadr e) a))) ((eq (car e) 'newkey) (newkey)) ((eq (car e) 'pub) (pub (eval (cadr e) a))) ((eq (car e) 'sign) (sign (eval (cadr e) a) (eval (caddr e) a))) ((eq (car e) 'verify) (verify (eval (cadr e) a) (eval (caddr e) a) (eval (cadddr e) a))) ((eq (car e) 'display) (display (eval (cadr e) a))) ((eq (car e) 'runes) (runes (eval (cadr e) a))) ((eq (car e) 'err) (err (eval (cadr e) a))) ((eq (car e) 'list) (evlis (cdr e) a)) ('t (eval (cons (assoc (car e) a) (cdr e)) a)))) ((eq (caar e) 'macro) (eval (display (eval (cadddar e) (pair (caddar e) (cdr e)))) a)) ((eq (caar e) 'label) (eval (cons (caddar e) (cdr e)) (cons (list (cadar e) (car e)) a))) ((eq (caar e) 'lambda) (cond ((atom (cadar e)) (eval (caddar e) (cons (list (cadar e) (evlis (cdr e) a)) a))) ('t (eval (caddar e) (append_go_sanitized (pair (cadar e) (evlis (cdr e) a)) a))))))))")

func eval(args ...Exp) Exp {
	e := args[0]
	a := args[1]
	return A(
		cond,
		L(
			Func(func(...Exp) Exp {
				return A(atom, e)
			}),
			Func(func(...Exp) Exp {
				return A(assoc, e, a)
			}),
		),
		L(
			Func(func(...Exp) Exp {
				return A(atom, A(car, e))
			}),
			Func(func(...Exp) Exp {
				return A(cond, L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "test")
					}),
					Func(func(...Exp) Exp {
						return A(test, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "test2")
					}),
					Func(func(...Exp) Exp {
						return A(test2, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "quote")
					}),
					Func(func(...Exp) Exp {
						return A(cadr, e)
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "atom")
					}),
					Func(func(...Exp) Exp {
						return A(atom, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "eq")
					}),
					Func(func(...Exp) Exp {
						return A(eq, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "car")
					}),
					Func(func(...Exp) Exp {
						return A(car, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "cdr")
					}),
					Func(func(...Exp) Exp {
						return A(cdr, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "cons")
					}),
					Func(func(...Exp) Exp {
						return A(cons, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "cond")
					}),
					Func(func(...Exp) Exp {
						return A(evcon, A(cdr, e), a)
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "plus")
					}),
					Func(func(...Exp) Exp {
						return A(plus, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "inc")
					}),
					Func(func(...Exp) Exp {
						return A(plus, A(eval, A(cadr, e), a), "1")
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "minus")
					}),
					Func(func(...Exp) Exp {
						return A(minus, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "mult")
					}),
					Func(func(...Exp) Exp {
						return A(mult, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "exp")
					}),
					Func(func(...Exp) Exp {
						return A(exp, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a), A(eval, A(cadddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "after")
					}),
					Func(func(...Exp) Exp {
						return A(after, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "concat")
					}),
					Func(func(...Exp) Exp {
						return A(concat, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "hash")
					}),
					Func(func(...Exp) Exp {
						return A(hash, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "newkey")
					}),
					Func(func(...Exp) Exp {
						return A(newkey)
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "pub")
					}),
					Func(func(...Exp) Exp {
						return A(pub, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "sign")
					}),
					Func(func(...Exp) Exp {
						return A(sign, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "verify")
					}),
					Func(func(...Exp) Exp {
						return A(verify, A(eval, A(cadr, e), a), A(eval, A(caddr, e), a), A(eval, A(cadddr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "display")
					}),
					Func(func(...Exp) Exp {
						return A(display, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "runes")
					}),
					Func(func(...Exp) Exp {
						return A(runes, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "err")
					}),
					Func(func(...Exp) Exp {
						return A(err, A(eval, A(cadr, e), a))
					}),
				), L(
					Func(func(...Exp) Exp {
						return A(eq, A(car, e), "list")
					}),
					Func(func(...Exp) Exp {
						return A(evlis, A(cdr, e), a)
					}),
				), L(
					"t",
					Func(func(...Exp) Exp {
						return A(eval, A(cons, A(assoc, A(car, e), a), A(cdr, e)), a)
					}),
				))
			}),
		),
		L(
			Func(func(...Exp) Exp {
				return A(eq, A(caar, e), "macro")
			}),
			Func(func(...Exp) Exp {
				return A(eval, A(display, A(eval, A(cadddar, e), A(pair, A(caddar, e), A(cdr, e)))), a)
			}),
		),
		L(
			Func(func(...Exp) Exp {
				return A(eq, A(caar, e), "label")
			}),
			Func(func(...Exp) Exp {
				return A(eval, A(cons, A(caddar, e), A(cdr, e)), A(cons, A(list, A(cadar, e), A(car, e)), a))
			}),
		),
		L(
			Func(func(...Exp) Exp {
				return A(eq, A(caar, e), "lambda")
			}),
			Func(func(...Exp) Exp {
				return A(cond, L(
					Func(func(...Exp) Exp {
						return A(atom, A(cadar, e))
					}),
					Func(func(...Exp) Exp {
						return A(eval, A(caddar, e), A(cons, A(list, A(cadar, e), A(evlis, A(cdr, e), a)), a))
					}),
				), L(
					"t",
					Func(func(...Exp) Exp {
						return A(eval, A(caddar, e), A(append_go_sanitized, A(pair, A(cadar, e), A(evlis, A(cdr, e), a)), a))
					}),
				))
			}),
		),
	)
}

//
// evcon (compiled)
//

var evcon_label = parse_env("(label evcon (lambda (c a) (cond ((eval (caar c) a) (eval (cadar c) a)) ('t (evcon (cdr c) a)))))")

func evcon(args ...Exp) Exp {
	c := args[0]
	a := args[1]
	return A(
		cond,
		L(
			Func(func(...Exp) Exp {
				return A(eval, A(caar, c), a)
			}),
			Func(func(...Exp) Exp {
				return A(eval, A(cadar, c), a)
			}),
		),
		L(
			"t",
			Func(func(...Exp) Exp {
				return A(evcon, A(cdr, c), a)
			}),
		),
	)
}

//
// evlis (compiled)
//

var evlis_label = parse_env("(label evlis (lambda (m a) (cond ((null m) '()) ('t (cons (eval (car m) a) (evlis (cdr m) a))))))")

func evlis(args ...Exp) Exp {
	m := args[0]
	a := args[1]
	return A(
		cond,
		L(
			Func(func(...Exp) Exp {
				return A(null, m)
			}),
			Nil,
		),
		L(
			"t",
			Func(func(...Exp) Exp {
				return A(cons, A(eval, A(car, m), a), A(evlis, A(cdr, m), a))
			}),
		),
	)
}

//
// factorial (interpreted)
//

var factorial_label = parse_env("(label factorial (lambda (n) (cond ((eq '0 n) '1) ('t (mult n (factorial (minus n '1)))))))")

//
// inc (interpreted)
//

var inc_label = parse_env("(label inc (lambda (x) (plus '1 x)))")

//
// lambdatest (interpreted)
//

var lambdatest_label = parse_env("(label lambdatest (lambda (x) (list (car x) (cdr x))))")

//
// length (interpreted)
//

var length_label = parse_env("(label length (lambda (x) (cond ((atom x) '0) ('t (plus '1 (length (cdr x)))))))")

//
// next (interpreted)
//

var next_label = parse_env("(label next (lambda (t) (cond ((eq (car t) (cadr t)) (err (list 'max (car t)))) ('t (list (car t) (inc (cadr t)))))))")

//
// not (compiled)
//

var not_label = parse_env("(label not (lambda (x) (cond (x '()) ('t 't))))")

func not(args ...Exp) Exp {
	x := args[0]
	return A(
		cond,
		L(
			x,
			Nil,
		),
		L(
			"t",
			"t",
		),
	)
}

//
// null (compiled)
//

var null_label = parse_env("(label null (lambda (x) (eq x '())))")

func null(args ...Exp) Exp {
	x := args[0]
	return A(
		eq,
		x,
		Nil,
	)
}

//
// pair (compiled)
//

var pair_label = parse_env("(label pair (lambda (x y) (cond ((and (null x) (null y)) '()) ((and (not (atom x)) (not (atom y))) (cons (list (car x) (car y)) (pair (cdr x) (cdr y)))))))")

func pair(args ...Exp) Exp {
	x := args[0]
	y := args[1]
	return A(
		cond,
		L(
			Func(func(...Exp) Exp {
				return A(and, A(null, x), A(null, y))
			}),
			Nil,
		),
		L(
			Func(func(...Exp) Exp {
				return A(and, A(not, A(atom, x)), A(not, A(atom, y)))
			}),
			Func(func(...Exp) Exp {
				return A(cons, A(list, A(car, x), A(car, y)), A(pair, A(cdr, x), A(cdr, y)))
			}),
		),
	)
}

//
// subst (interpreted)
//

var subst_label = parse_env("(label subst (lambda (x y z) (cond ((atom z) (cond ((eq z y) x) ('t z))) ('t (cons (subst x y (car z)) (subst x y (cdr z)))))))")

//
// tassoc (interpreted)
//

var tassoc_label = parse_env("(label tassoc (lambda (t x y) (cond ((eq (caar y) x) (cadar y)) ('t (tassoc (next t) x (cdr y))))))")

//
// test (compiled)
//

var test_label = parse_env("(label test (lambda (x) ((lambda (first rest) (list first rest)) (car x) (cdr x))))")

func test(args ...Exp) Exp {
	x := args[0]
	return func() Exp {
		var lambda func(...Exp) Exp
		lambda = func(args ...Exp) Exp {
			first := args[0]
			rest := args[1]
			return A(list, first, rest)
		}
		return lambda(A(car, x), A(cdr, x))
	}()

}

//
// test2 (compiled)
//

var test2_label = parse_env("(label test2 (lambda (x) ((label f (lambda (first rest) (list first rest))) (car x) (cdr x))))")

func test2(args ...Exp) Exp {
	x := args[0]
	return func() Exp {
		var f func(...Exp) Exp
		f = func(args ...Exp) Exp {
			first := args[0]
			rest := args[1]
			return A(list, first, rest)
		}
		return f(A(car, x), A(cdr, x))
	}()

}

//
// test3 (compiled)
//

var test3_label = parse_env("(label test3 (lambda (x) ((label fx (lambda (first rest) (cond ((eq first '0) (list first rest)) ('t (fx (minus first '1) rest))))) (car x) (cdr x))))")

func test3(args ...Exp) Exp {
	x := args[0]
	return func() Exp {
		var fx func(...Exp) Exp
		fx = func(args ...Exp) Exp {
			first := args[0]
			rest := args[1]
			return A(cond, L(
				Func(func(...Exp) Exp {
					return A(eq, first, "0")
				}),
				Func(func(...Exp) Exp {
					return A(list, first, rest)
				}),
			), L(
				"t",
				Func(func(...Exp) Exp {
					return A(fx, A(minus, first, "1"), rest)
				}),
			))
		}
		return fx(A(car, x), A(cdr, x))
	}()

}

//
// test4 (interpreted)
//

var test4_label = parse_env("(label test4 (lambda (x) ((label f (lambda (first rest) (cond ((eq first '0) (list first rest)) ('t (f (minus first '1) rest))))) (car x) (cdr x))))")

//
// try (interpreted)
//

var try_label = parse_env("(label try (lambda (t e a) (cond ((atom e) (assoc e a)) ((atom (car e)) (cond ((eq (car e) 'quote) (cadr e)) ((eq (car e) 'atom) (atom (try (next t) (cadr e) a))) ((eq (car e) 'eq) (eq (try (next t) (cadr e) a) (try (next t) (caddr e) a))) ((eq (car e) 'car) (car (try (next t) (cadr e) a))) ((eq (car e) 'cdr) (cdr (try (next t) (cadr e) a))) ((eq (car e) 'cons) (cons (try (next t) (cadr e) a) (try (next t) (caddr e) a))) ((eq (car e) 'cond) (evcon (cdr e) a)) ((eq (car e) 'list) (evlis (cdr e) a)) ('t (try (next t) (cons (tassoc (next t) (car e) a) (cdr e)) a)))) ((eq (caar e) 'label) (try (next t) (cons (caddar e) (cdr e)) (cons (list (cadar e) (car e)) a))) ((eq (caar e) 'lambda) (try (next t) (caddar e) (append_go_sanitized (pair (cadar e) (evlis (cdr e) a)) a))))))")

//
// xlist (interpreted)
//

var xlist_label = parse_env("(label xlist (lambda x x))")

func init() {
	env = L(
		L("and", and_label),
		L("append_go_sanitized", append_go_sanitized_label),
		L("assoc", assoc_label),
		L("caar", caar_label),
		L("cadar", cadar_label),
		L("caddar", caddar_label),
		L("cadddar", cadddar_label),
		L("caddddar", caddddar_label),
		L("caddddr", caddddr_label),
		L("cadddr", cadddr_label),
		L("caddr", caddr_label),
		L("cadr", cadr_label),
		L("cdar", cdar_label),
		L("cdddar", cdddar_label),
		L("eval", eval_label),
		L("evcon", evcon_label),
		L("evlis", evlis_label),
		L("factorial", factorial_label),
		L("inc", inc_label),
		L("lambdatest", lambdatest_label),
		L("length", length_label),
		L("next", next_label),
		L("not", not_label),
		L("null", null_label),
		L("pair", pair_label),
		L("subst", subst_label),
		L("tassoc", tassoc_label),
		L("test", test_label),
		L("test2", test2_label),
		L("test3", test3_label),
		L("test4", test4_label),
		L("try", try_label),
		L("xlist", xlist_label),
	)
}
