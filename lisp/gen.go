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
// caadr (compiled)
//

var caadr_label = parse_env("(label caadr (lambda (x) (car (car (cdr x)))))")

func caadr(args ...Exp) Exp {
	x := args[0]
	return A(
		car,
		A(
			car,
			A(
				cdr,
				x,
			),
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
// caddddr (compiled)
//

var caddddr_label = parse_env("(label caddddr (lambda (x) (car (cdr (cdr (cdr (cdr x)))))))")

func caddddr(args ...Exp) Exp {
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
						x,
					),
				),
			),
		),
	)
}

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
// cdar (compiled)
//

var cdar_label = parse_env("(label cdar (lambda (x) (cdr (car x))))")

func cdar(args ...Exp) Exp {
	x := args[0]
	return A(
		cdr,
		A(
			car,
			x,
		),
	)
}

//
// cddar (compiled)
//

var cddar_label = parse_env("(label cddar (lambda (x) (cdr (cdr (car x)))))")

func cddar(args ...Exp) Exp {
	x := args[0]
	return A(
		cdr,
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
// cddr (compiled)
//

var cddr_label = parse_env("(label cddr (lambda (x) (cdr (cdr x))))")

func cddr(args ...Exp) Exp {
	x := args[0]
	return A(
		cdr,
		A(
			cdr,
			x,
		),
	)
}

//
// eval (compiled)
//

var eval_label = parse_env("(label eval (lambda (e a) (cond ((atom e) (assoc e a)) ((atom (car e)) ((λ (op first rest) ((λ (second third) (cond ((eq op 'funcall) (eval (cons (eval first a) rest) a)) ((eq op 'quote) (cadr e)) ((eq op 'cond) (evcon (cdr e) a)) ((eq op 'list) (evlis (cdr e) a)) ((eq op 'add) (add (eval first a) (eval second a))) ((eq op 'after) (after (eval first a) (eval second a))) ((eq op 'and) (and (eval first a) (eval second a))) ((eq op 'append_go_sanitized) (append_go_sanitized (eval first a) (eval second a))) ((eq op 'assoc) (assoc (eval first a) (eval second a))) ((eq op 'atom) (atom (eval first a))) ((eq op 'caadr) (caadr (eval first a))) ((eq op 'caar) (caar (eval first a))) ((eq op 'cadar) (cadar (eval first a))) ((eq op 'caddar) (caddar (eval first a))) ((eq op 'cadddar) (cadddar (eval first a))) ((eq op 'caddddar) (caddddar (eval first a))) ((eq op 'caddddr) (caddddr (eval first a))) ((eq op 'cadddr) (cadddr (eval first a))) ((eq op 'caddr) (caddr (eval first a))) ((eq op 'cadr) (cadr (eval first a))) ((eq op 'car) (car (eval first a))) ((eq op 'cdar) (cdar (eval first a))) ((eq op 'cddar) (cddar (eval first a))) ((eq op 'cdddar) (cdddar (eval first a))) ((eq op 'cddr) (cddr (eval first a))) ((eq op 'cdr) (cdr (eval first a))) ((eq op 'concat) (concat (eval first a) (eval second a))) ((eq op 'cons) (cons (eval first a) (eval second a))) ((eq op 'display) (display (eval first a))) ((eq op 'eq) (eq (eval first a) (eval second a))) ((eq op 'err) (err (eval first a))) ((eq op 'eval) (eval (eval first a) (eval second a))) ((eq op 'evcon) (evcon (eval first a) (eval second a))) ((eq op 'evlis) (evlis (eval first a) (eval second a))) ((eq op 'exp) (exp (eval first a) (eval second a) (eval third a))) ((eq op 'hash) (hash (eval first a))) ((eq op 'hashed) (hashed (eval first a))) ((eq op 'inc) (inc (eval first a))) ((eq op 'length) (length (eval first a))) ((eq op 'mul) (mul (eval first a) (eval second a))) ((eq op 'newkey) (newkey)) ((eq op 'next) (next (eval first a))) ((eq op 'not) (not (eval first a))) ((eq op 'null) (null (eval first a))) ((eq op 'or) (or (eval first a) (eval second a))) ((eq op 'pair) (pair (eval first a) (eval second a))) ((eq op 'pub) (pub (eval first a))) ((eq op 'runes) (runes (eval first a))) ((eq op 'sign) (sign (eval first a) (eval second a))) ((eq op 'sub) (sub (eval first a) (eval second a))) ((eq op 'tassoc) (tassoc (eval first a) (eval second a) (eval third a))) ((eq op 'test1) (test1 (eval first a))) ((eq op 'test2) (test2 (eval first a))) ((eq op 'test3) (test3 (eval first a))) ((eq op 'tevcon) (tevcon (eval first a) (eval second a) (eval third a))) ((eq op 'tevlis) (tevlis (eval first a) (eval second a) (eval third a))) ((eq op 'verify) (verify (eval first a) (eval second a) (eval third a))) ('t (eval (cons (assoc op a) (cdr e)) a)))) (car rest) (cadr rest))) (car e) (cadr e) (cddr e))) ((eq (caar e) 'macro) (eval (eval (cadddar e) (pair (caddar e) (cdr e))) a)) ((eq (caar e) 'label) (eval (cons (caddar e) (cdr e)) (cons (list (cadar e) (car e)) a))) ((or (eq (caar e) 'lambda) (eq (caar e) 'λ)) (cond ((atom (cadar e)) (eval (caddar e) (cons (list (cadar e) (evlis (cdr e) a)) a))) ('t (eval (caddar e) (append_go_sanitized (pair (cadar e) (evlis (cdr e) a)) a))))))))")

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
				return func() Exp {
					var λ func(...Exp) Exp
					λ = func(args ...Exp) Exp {
						op := args[0]
						first := args[1]
						rest := args[2]
						return func() Exp {
							var λ func(...Exp) Exp
							λ = func(args ...Exp) Exp {
								second := args[0]
								third := args[1]
								return func() Exp {
									if f, ok := map_49f1753b18[String(op)]; ok {
										return f(a, e, first, op, rest, second, third)
									}
									return A(eval, A(cons, A(assoc, op, a), A(cdr, e)), a)
								}()

							}
							return λ(A(car, rest), A(cadr, rest))
						}()

					}
					return λ(A(car, e), A(cadr, e), A(cddr, e))
				}()

			}),
		),
		L(
			Func(func(...Exp) Exp {
				return A(eq, A(caar, e), "macro")
			}),
			Func(func(...Exp) Exp {
				return A(eval, A(eval, A(cadddar, e), A(pair, A(caddar, e), A(cdr, e))), a)
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
				return A(or, A(eq, A(caar, e), "lambda"), A(eq, A(caar, e), "λ"))
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

var factorial_label = parse_env("(label factorial (lambda (n) (cond ((eq '0 n) '1) ('t (mul n (factorial (sub n '1)))))))")

//
// inc (compiled)
//

var inc_label = parse_env("(label inc (lambda (x) (add '1 x)))")

func inc(args ...Exp) Exp {
	x := args[0]
	return A(
		add,
		"1",
		x,
	)
}

//
// lambdatest (interpreted)
//

var lambdatest_label = parse_env("(label lambdatest (lambda (x) (list (car x) (cdr x))))")

//
// length (compiled)
//

var length_label = parse_env("(label length (lambda (x) (cond ((atom x) '0) ('t (add '1 (length (cdr x)))))))")

func length(args ...Exp) Exp {
	x := args[0]
	return A(
		cond,
		L(
			Func(func(...Exp) Exp {
				return A(atom, x)
			}),
			"0",
		),
		L(
			"t",
			Func(func(...Exp) Exp {
				return A(add, "1", A(length, A(cdr, x)))
			}),
		),
	)
}

//
// mapcar (interpreted)
//

var mapcar_label = parse_env("(label mapcar (lambda (op arglist) (cond ((eq arglist '()) ()) ('t (cons (funcall op (car arglist)) (mapcar op (cdr arglist)))))))")

//
// next (compiled)
//

var next_label = parse_env("(label next (lambda (t) ((lambda (max current) (cond ((eq max current) (err (list 'max max))) ('t (list max (inc current))))) (car t) (cadr t))))")

func next(args ...Exp) Exp {
	t := args[0]
	return func() Exp {
		var λ func(...Exp) Exp
		λ = func(args ...Exp) Exp {
			max := args[0]
			current := args[1]
			return A(cond, L(
				Func(func(...Exp) Exp {
					return A(eq, max, current)
				}),
				Func(func(...Exp) Exp {
					return A(err, A(list, "max", max))
				}),
			), L(
				"t",
				Func(func(...Exp) Exp {
					return A(list, max, A(inc, current))
				}),
			))
		}
		return λ(A(car, t), A(cadr, t))
	}()

}

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
// or (compiled)
//

var or_label = parse_env("(label or (lambda (x y) (cond (x 't) (y 't) ('t '()))))")

func or(args ...Exp) Exp {
	x := args[0]
	y := args[1]
	return A(
		cond,
		L(
			x,
			"t",
		),
		L(
			y,
			"t",
		),
		L(
			"t",
			Nil,
		),
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
// s (interpreted)
//

var s_label = parse_env("(label s (lambda (f x) (f f x)))")

//
// subst (interpreted)
//

var subst_label = parse_env("(label subst (lambda (x y z) (cond ((atom z) (cond ((eq z y) x) ('t z))) ('t (cons (subst x y (car z)) (subst x y (cdr z)))))))")

//
// tassoc (compiled)
//

var tassoc_label = parse_env("(label tassoc (lambda (t x y) (cond ((eq (caar y) x) (cadar y)) ('t (tassoc (next t) x (cdr y))))))")

func tassoc(args ...Exp) Exp {
	t := args[0]
	x := args[1]
	y := args[2]
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
				return A(tassoc, A(next, t), x, A(cdr, y))
			}),
		),
	)
}

//
// test1 (compiled)
//

var test1_label = parse_env("(label test1 (lambda (x) ((lambda (first rest) (list first rest)) (car x) (cdr x))))")

func test1(args ...Exp) Exp {
	x := args[0]
	return func() Exp {
		var λ func(...Exp) Exp
		λ = func(args ...Exp) Exp {
			first := args[0]
			rest := args[1]
			return A(list, first, rest)
		}
		return λ(A(car, x), A(cdr, x))
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

var test3_label = parse_env("(label test3 (lambda (x) ((label fx (lambda (first rest) (cond ((eq first '0) (list first rest)) ('t (fx (sub first '1) rest))))) (car x) (cdr x))))")

func test3(args ...Exp) Exp {
	x := args[0]
	return func() Exp {
		var fx func(...Exp) Exp
		fx = func(args ...Exp) Exp {
			first := args[0]
			rest := args[1]
			return func() Exp {
				if f, ok := map_405e36fa33[String(first)]; ok {
					return f(first, rest, x)
				}
				return A(fx, A(sub, first, "1"), rest)
			}()

		}
		return fx(A(car, x), A(cdr, x))
	}()

}

//
// test4 (interpreted)
//

var test4_label = parse_env("(label test4 (lambda (x) ((label f (lambda (first rest) (cond ((eq first '0) (list first rest)) ('t (f (sub first '1) rest))))) (car x) (cdr x))))")

//
// teval (compiled)
//

var teval_label = parse_env("(label teval (lambda (x y z) (list x y z)))")

func teval(args ...Exp) Exp {
	x := args[0]
	y := args[1]
	z := args[2]
	return A(
		list,
		x,
		y,
		z,
	)
}

//
// tevcon (compiled)
//

var tevcon_label = parse_env("(label tevcon (lambda (t c a) (cond ((teval (next t) (caar c) a) (teval (next t) (cadar c) a)) ('t (tevcon (next t) (cdr c) a)))))")

func tevcon(args ...Exp) Exp {
	t := args[0]
	c := args[1]
	a := args[2]
	return A(
		cond,
		L(
			Func(func(...Exp) Exp {
				return A(teval, A(next, t), A(caar, c), a)
			}),
			Func(func(...Exp) Exp {
				return A(teval, A(next, t), A(cadar, c), a)
			}),
		),
		L(
			"t",
			Func(func(...Exp) Exp {
				return A(tevcon, A(next, t), A(cdr, c), a)
			}),
		),
	)
}

//
// tevlis (compiled)
//

var tevlis_label = parse_env("(label tevlis (lambda (t m a) (cond ((null m) '()) ('t (cons (teval (next t) (car m) a) (tevlis (next t) (cdr m) a))))))")

func tevlis(args ...Exp) Exp {
	t := args[0]
	m := args[1]
	a := args[2]
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
				return A(cons, A(teval, A(next, t), A(car, m), a), A(tevlis, A(next, t), A(cdr, m), a))
			}),
		),
	)
}

//
// try (interpreted)
//

var try_label = parse_env("(label try (lambda (t e a) (cond ((atom e) (assoc e a)) ((atom (car e)) (cond ((eq (car e) 'quote) (cadr e)) ((eq (car e) 'atom) (atom (try (next t) (cadr e) a))) ((eq (car e) 'eq) (eq (try (next t) (cadr e) a) (try (next t) (caddr e) a))) ((eq (car e) 'car) (car (try (next t) (cadr e) a))) ((eq (car e) 'cdr) (cdr (try (next t) (cadr e) a))) ((eq (car e) 'cons) (cons (try (next t) (cadr e) a) (try (next t) (caddr e) a))) ((eq (car e) 'cond) (evcon (cdr e) a)) ((eq (car e) 'list) (evlis (cdr e) a)) ('t (try (next t) (cons (tassoc (next t) (car e) a) (cdr e)) a)))) ((eq (caar e) 'label) (try (next t) (cons (caddar e) (cdr e)) (cons (list (cadar e) (car e)) a))) ((eq (caar e) 'lambda) (try (next t) (caddar e) (append_go_sanitized (pair (cadar e) (evlis (cdr e) a)) a))))))")

//
// xlist (interpreted)
//

var xlist_label = parse_env("(label xlist (lambda x x))")

//
// y (interpreted)
//

var y_label = parse_env("(label y (lambda (f) ((lambda (x) (f (x x))) (lambda (x) (f (x x))))))")

// cases:

func F_0_94ffb94c37(first, rest, x Exp) Exp {
	return A(list, first, rest)
}

func F_add_c697c7bfbf(a, e, first, op, rest, second, third Exp) Exp {
	return A(add, A(eval, first, a), A(eval, second, a))
}

func F_after_d1d1801dce(a, e, first, op, rest, second, third Exp) Exp {
	return A(after, A(eval, first, a), A(eval, second, a))
}

func F_and_b62b5eac11(a, e, first, op, rest, second, third Exp) Exp {
	return A(and, A(eval, first, a), A(eval, second, a))
}

func F_appendπgoπsanitized_5ca651c73e(a, e, first, op, rest, second, third Exp) Exp {
	return A(append_go_sanitized, A(eval, first, a), A(eval, second, a))
}

func F_assoc_58df95b48f(a, e, first, op, rest, second, third Exp) Exp {
	return A(assoc, A(eval, first, a), A(eval, second, a))
}

func F_atom_57e6d5c9b3(a, e, first, op, rest, second, third Exp) Exp {
	return A(atom, A(eval, first, a))
}

func F_caadr_9a51460d6b(a, e, first, op, rest, second, third Exp) Exp {
	return A(caadr, A(eval, first, a))
}

func F_caar_18f57ec7ac(a, e, first, op, rest, second, third Exp) Exp {
	return A(caar, A(eval, first, a))
}

func F_cadar_08bb39827e(a, e, first, op, rest, second, third Exp) Exp {
	return A(cadar, A(eval, first, a))
}

func F_caddar_bb34829244(a, e, first, op, rest, second, third Exp) Exp {
	return A(caddar, A(eval, first, a))
}

func F_cadddar_26ece2d994(a, e, first, op, rest, second, third Exp) Exp {
	return A(cadddar, A(eval, first, a))
}

func F_caddddar_4b02e05a2f(a, e, first, op, rest, second, third Exp) Exp {
	return A(caddddar, A(eval, first, a))
}

func F_caddddr_d7840de5c0(a, e, first, op, rest, second, third Exp) Exp {
	return A(caddddr, A(eval, first, a))
}

func F_cadddr_b15fa6eba5(a, e, first, op, rest, second, third Exp) Exp {
	return A(cadddr, A(eval, first, a))
}

func F_caddr_54b1e7324a(a, e, first, op, rest, second, third Exp) Exp {
	return A(caddr, A(eval, first, a))
}

func F_cadr_aa1c6639fb(a, e, first, op, rest, second, third Exp) Exp {
	return A(cadr, A(eval, first, a))
}

func F_car_6e7ecda3ef(a, e, first, op, rest, second, third Exp) Exp {
	return A(car, A(eval, first, a))
}

func F_cdar_f21811a8cc(a, e, first, op, rest, second, third Exp) Exp {
	return A(cdar, A(eval, first, a))
}

func F_cddar_d9b0da6fb9(a, e, first, op, rest, second, third Exp) Exp {
	return A(cddar, A(eval, first, a))
}

func F_cdddar_89f98ebd5c(a, e, first, op, rest, second, third Exp) Exp {
	return A(cdddar, A(eval, first, a))
}

func F_cddr_b3ffda7668(a, e, first, op, rest, second, third Exp) Exp {
	return A(cddr, A(eval, first, a))
}

func F_cdr_6772863567(a, e, first, op, rest, second, third Exp) Exp {
	return A(cdr, A(eval, first, a))
}

func F_concat_832e46a008(a, e, first, op, rest, second, third Exp) Exp {
	return A(concat, A(eval, first, a), A(eval, second, a))
}

func F_cond_26e96b4be1(a, e, first, op, rest, second, third Exp) Exp {
	return A(evcon, A(cdr, e), a)
}

func F_cons_d4b73be861(a, e, first, op, rest, second, third Exp) Exp {
	return A(cons, A(eval, first, a), A(eval, second, a))
}

func F_display_f82db8af96(a, e, first, op, rest, second, third Exp) Exp {
	return A(display, A(eval, first, a))
}

func F_eq_9d693d8748(a, e, first, op, rest, second, third Exp) Exp {
	return A(eq, A(eval, first, a), A(eval, second, a))
}

func F_err_9d3d836b35(a, e, first, op, rest, second, third Exp) Exp {
	return A(err, A(eval, first, a))
}

func F_eval_00e5f48f55(a, e, first, op, rest, second, third Exp) Exp {
	return A(eval, A(eval, first, a), A(eval, second, a))
}

func F_evcon_848079d5b8(a, e, first, op, rest, second, third Exp) Exp {
	return A(evcon, A(eval, first, a), A(eval, second, a))
}

func F_evlis_b9bd2f7eb3(a, e, first, op, rest, second, third Exp) Exp {
	return A(evlis, A(eval, first, a), A(eval, second, a))
}

func F_exp_0cf8749970(a, e, first, op, rest, second, third Exp) Exp {
	return A(exp, A(eval, first, a), A(eval, second, a), A(eval, third, a))
}

func F_funcall_38c9a26765(a, e, first, op, rest, second, third Exp) Exp {
	return A(eval, A(cons, A(eval, first, a), rest), a)
}

func F_hash_890341d522(a, e, first, op, rest, second, third Exp) Exp {
	return A(hash, A(eval, first, a))
}

func F_hashed_2106668384(a, e, first, op, rest, second, third Exp) Exp {
	return A(hashed, A(eval, first, a))
}

func F_inc_878d5d4d19(a, e, first, op, rest, second, third Exp) Exp {
	return A(inc, A(eval, first, a))
}

func F_length_39633e2406(a, e, first, op, rest, second, third Exp) Exp {
	return A(length, A(eval, first, a))
}

func F_list_4d554d264a(a, e, first, op, rest, second, third Exp) Exp {
	return A(evlis, A(cdr, e), a)
}

func F_mul_3841a9b191(a, e, first, op, rest, second, third Exp) Exp {
	return A(mul, A(eval, first, a), A(eval, second, a))
}

func F_newkey_623d29c094(a, e, first, op, rest, second, third Exp) Exp {
	return A(newkey)
}

func F_next_5f5d06563c(a, e, first, op, rest, second, third Exp) Exp {
	return A(next, A(eval, first, a))
}

func F_not_d06984edf2(a, e, first, op, rest, second, third Exp) Exp {
	return A(not, A(eval, first, a))
}

func F_null_0fce507675(a, e, first, op, rest, second, third Exp) Exp {
	return A(null, A(eval, first, a))
}

func F_or_c6b938191a(a, e, first, op, rest, second, third Exp) Exp {
	return A(or, A(eval, first, a), A(eval, second, a))
}

func F_pair_743e252bfd(a, e, first, op, rest, second, third Exp) Exp {
	return A(pair, A(eval, first, a), A(eval, second, a))
}

func F_pub_fd1fdb63c7(a, e, first, op, rest, second, third Exp) Exp {
	return A(pub, A(eval, first, a))
}

func F_quote_a7f3dfeaaf(a, e, first, op, rest, second, third Exp) Exp {
	return A(cadr, e)
}

func F_runes_1f36010f59(a, e, first, op, rest, second, third Exp) Exp {
	return A(runes, A(eval, first, a))
}

func F_sign_11c1de489d(a, e, first, op, rest, second, third Exp) Exp {
	return A(sign, A(eval, first, a), A(eval, second, a))
}

func F_sub_246a160bc3(a, e, first, op, rest, second, third Exp) Exp {
	return A(sub, A(eval, first, a), A(eval, second, a))
}

func F_tassoc_038d29f713(a, e, first, op, rest, second, third Exp) Exp {
	return A(tassoc, A(eval, first, a), A(eval, second, a), A(eval, third, a))
}

func F_test1_b2a90c1647(a, e, first, op, rest, second, third Exp) Exp {
	return A(test1, A(eval, first, a))
}

func F_test2_97cace0d47(a, e, first, op, rest, second, third Exp) Exp {
	return A(test2, A(eval, first, a))
}

func F_test3_54e98f673d(a, e, first, op, rest, second, third Exp) Exp {
	return A(test3, A(eval, first, a))
}

func F_tevcon_79bdb10a3c(a, e, first, op, rest, second, third Exp) Exp {
	return A(tevcon, A(eval, first, a), A(eval, second, a), A(eval, third, a))
}

func F_tevlis_e770634921(a, e, first, op, rest, second, third Exp) Exp {
	return A(tevlis, A(eval, first, a), A(eval, second, a), A(eval, third, a))
}

func F_verify_5199556588(a, e, first, op, rest, second, third Exp) Exp {
	return A(verify, A(eval, first, a), A(eval, second, a), A(eval, third, a))
}

func F_ππ_7069b7d1ee(arglist, op Exp) Exp {
	return Nil
}

var map_283c6e3b83 = make(map[string]func(arglist, op Exp) Exp)

func init() {
	map_283c6e3b83 = map[string]func(arglist, op Exp) Exp{
		"()": F_ππ_7069b7d1ee,
	}
}

var map_405e36fa33 = make(map[string]func(first, rest, x Exp) Exp)

func init() {
	map_405e36fa33 = map[string]func(first, rest, x Exp) Exp{
		"0": F_0_94ffb94c37,
	}
}

var map_49f1753b18 = make(map[string]func(a, e, first, op, rest, second, third Exp) Exp)

func init() {
	map_49f1753b18 = map[string]func(a, e, first, op, rest, second, third Exp) Exp{
		"funcall":             F_funcall_38c9a26765,
		"quote":               F_quote_a7f3dfeaaf,
		"cond":                F_cond_26e96b4be1,
		"list":                F_list_4d554d264a,
		"add":                 F_add_c697c7bfbf,
		"after":               F_after_d1d1801dce,
		"and":                 F_and_b62b5eac11,
		"append_go_sanitized": F_appendπgoπsanitized_5ca651c73e,
		"assoc":               F_assoc_58df95b48f,
		"atom":                F_atom_57e6d5c9b3,
		"caadr":               F_caadr_9a51460d6b,
		"caar":                F_caar_18f57ec7ac,
		"cadar":               F_cadar_08bb39827e,
		"caddar":              F_caddar_bb34829244,
		"cadddar":             F_cadddar_26ece2d994,
		"caddddar":            F_caddddar_4b02e05a2f,
		"caddddr":             F_caddddr_d7840de5c0,
		"cadddr":              F_cadddr_b15fa6eba5,
		"caddr":               F_caddr_54b1e7324a,
		"cadr":                F_cadr_aa1c6639fb,
		"car":                 F_car_6e7ecda3ef,
		"cdar":                F_cdar_f21811a8cc,
		"cddar":               F_cddar_d9b0da6fb9,
		"cdddar":              F_cdddar_89f98ebd5c,
		"cddr":                F_cddr_b3ffda7668,
		"cdr":                 F_cdr_6772863567,
		"concat":              F_concat_832e46a008,
		"cons":                F_cons_d4b73be861,
		"display":             F_display_f82db8af96,
		"eq":                  F_eq_9d693d8748,
		"err":                 F_err_9d3d836b35,
		"eval":                F_eval_00e5f48f55,
		"evcon":               F_evcon_848079d5b8,
		"evlis":               F_evlis_b9bd2f7eb3,
		"exp":                 F_exp_0cf8749970,
		"hash":                F_hash_890341d522,
		"hashed":              F_hashed_2106668384,
		"inc":                 F_inc_878d5d4d19,
		"length":              F_length_39633e2406,
		"mul":                 F_mul_3841a9b191,
		"newkey":              F_newkey_623d29c094,
		"next":                F_next_5f5d06563c,
		"not":                 F_not_d06984edf2,
		"null":                F_null_0fce507675,
		"or":                  F_or_c6b938191a,
		"pair":                F_pair_743e252bfd,
		"pub":                 F_pub_fd1fdb63c7,
		"runes":               F_runes_1f36010f59,
		"sign":                F_sign_11c1de489d,
		"sub":                 F_sub_246a160bc3,
		"tassoc":              F_tassoc_038d29f713,
		"test1":               F_test1_b2a90c1647,
		"test2":               F_test2_97cace0d47,
		"test3":               F_test3_54e98f673d,
		"tevcon":              F_tevcon_79bdb10a3c,
		"tevlis":              F_tevlis_e770634921,
		"verify":              F_verify_5199556588,
	}
}

func init() {
	env = L(
		L("and", and_label),
		L("append_go_sanitized", append_go_sanitized_label),
		L("assoc", assoc_label),
		L("caadr", caadr_label),
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
		L("cddar", cddar_label),
		L("cdddar", cdddar_label),
		L("cddr", cddr_label),
		L("eval", eval_label),
		L("evcon", evcon_label),
		L("evlis", evlis_label),
		L("factorial", factorial_label),
		L("inc", inc_label),
		L("lambdatest", lambdatest_label),
		L("length", length_label),
		L("mapcar", mapcar_label),
		L("next", next_label),
		L("not", not_label),
		L("null", null_label),
		L("or", or_label),
		L("pair", pair_label),
		L("s", s_label),
		L("subst", subst_label),
		L("tassoc", tassoc_label),
		L("test1", test1_label),
		L("test2", test2_label),
		L("test3", test3_label),
		L("test4", test4_label),
		L("teval", teval_label),
		L("tevcon", tevcon_label),
		L("tevlis", tevlis_label),
		L("try", try_label),
		L("xlist", xlist_label),
		L("y", y_label),
	)
}
