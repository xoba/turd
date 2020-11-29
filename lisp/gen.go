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

var eval_label = parse_env("(label eval (lambda (e a) (cond ((atom e) (assoc e a)) ((atom (car e)) ((lambda (op first second third) (cond ((eq op 'test1) (test1 (eval first a))) ((eq op 'test2) (test2 (eval first a))) ((eq op 'quote) (cadr e)) ((eq op 'atom) (atom (eval first a))) ((eq op 'eq) (eq (eval first a) (eval second a))) ((eq op 'car) (car (eval first a))) ((eq op 'cdr) (cdr (eval first a))) ((eq op 'cons) (cons (eval first a) (eval second a))) ((eq op 'cond) (evcon (cdr e) a)) ((eq op 'add) (add (eval first a) (eval second a))) ((eq op 'inc) (plus (eval first a) '1)) ((eq op 'sub) (sub (eval first a) (eval second a))) ((eq op 'mul) (mul (eval first a) (eval second a))) ((eq op 'exp) (exp (eval first a) (eval second a) (eval third a))) ((eq op 'after) (after (eval first a) (eval second a))) ((eq op 'concat) (concat (eval first a) (eval second a))) ((eq op 'hash) (hash (eval first a))) ((eq op 'newkey) (newkey)) ((eq op 'pub) (pub (eval first a))) ((eq op 'sign) (sign (eval first a) (eval second a))) ((eq op 'verify) (verify (eval first a) (eval second a) (eval third a))) ((eq op 'display) (display (eval first a))) ((eq op 'runes) (runes (eval (cadr e) a))) ((eq op 'err) (err (eval (cadr e) a))) ((eq op 'list) (evlis (cdr e) a)) ('t (eval (cons (assoc op a) (cdr e)) a)))) (car e) (cadr e) (caddr e) (cadddr e))) ((eq (caar e) 'macro) (eval (eval (cadddar e) (pair (caddar e) (cdr e))) a)) ((eq (caar e) 'label) (eval (cons (caddar e) (cdr e)) (cons (list (cadar e) (car e)) a))) ((eq (caar e) 'lambda) (cond ((atom (cadar e)) (eval (caddar e) (cons (list (cadar e) (evlis (cdr e) a)) a))) ('t (eval (caddar e) (append_go_sanitized (pair (cadar e) (evlis (cdr e) a)) a))))))))")

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
					var lambda func(...Exp) Exp
					lambda = func(args ...Exp) Exp {
						op := args[0]
						first := args[1]
						second := args[2]
						third := args[3]
						return func() Exp {
							if f, ok := map_25[String(op)]; ok {
								return f(e, a, op, first, second, third)
							}
							return A(eval, A(cons, A(assoc, op, a), A(cdr, e)), a)
						}()

					}
					return lambda(A(car, e), A(cadr, e), A(caddr, e), A(cadddr, e))
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

var factorial_label = parse_env("(label factorial (lambda (n) (cond ((eq '0 n) '1) ('t (mul n (factorial (sub n '1)))))))")

//
// inc (interpreted)
//

var inc_label = parse_env("(label inc (lambda (x) (add '1 x)))")

//
// lambdatest (interpreted)
//

var lambdatest_label = parse_env("(label lambdatest (lambda (x) (list (car x) (cdr x))))")

//
// length (interpreted)
//

var length_label = parse_env("(label length (lambda (x) (cond ((atom x) '0) ('t (add '1 (length (cdr x)))))))")

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
// s (interpreted)
//

var s_label = parse_env("(label s (lambda (f x) (f f x)))")

//
// subst (interpreted)
//

var subst_label = parse_env("(label subst (lambda (x y z) (cond ((atom z) (cond ((eq z y) x) ('t z))) ('t (cons (subst x y (car z)) (subst x y (cdr z)))))))")

//
// tassoc (interpreted)
//

var tassoc_label = parse_env("(label tassoc (lambda (t x y) (cond ((eq (caar y) x) (cadar y)) ('t (tassoc (next t) x (cdr y))))))")

//
// test1 (compiled)
//

var test1_label = parse_env("(label test1 (lambda (x) ((lambda (first rest) (list first rest)) (car x) (cdr x))))")

func test1(args ...Exp) Exp {
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

var test3_label = parse_env("(label test3 (lambda (x) ((label fx (lambda (first rest) (cond ((eq first '0) (list first rest)) ('t (fx (sub first '1) rest))))) (car x) (cdr x))))")

func test3(args ...Exp) Exp {
	x := args[0]
	return func() Exp {
		var fx func(...Exp) Exp
		fx = func(args ...Exp) Exp {
			first := args[0]
			rest := args[1]
			return func() Exp {
				if f, ok := map_27[String(first)]; ok {
					return f(x, first, rest)
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

func F_0_60a3caf220a5e9e8d986965c79b20f11(x, first, rest Exp) Exp {
	return A(list, first, rest)
}

func F_add_cb8e93847499008be09106dc15141a43(e, a, op, first, second, third Exp) Exp {
	return A(add, A(eval, first, a), A(eval, second, a))
}

func F_after_515cbbd8f947c37b7d1dc2951e924a4d(e, a, op, first, second, third Exp) Exp {
	return A(after, A(eval, first, a), A(eval, second, a))
}

func F_atom_8f3c75e470915a3502aaee0ca1577fcc(e, a, op, first, second, third Exp) Exp {
	return A(atom, A(eval, first, a))
}

func F_car_9ccdcb1ba9d132a44726ab42f27ff4e9(e, a, op, first, second, third Exp) Exp {
	return A(car, A(eval, first, a))
}

func F_cdr_07b8b4fee0b4a65b5f5e3e5580aaa311(e, a, op, first, second, third Exp) Exp {
	return A(cdr, A(eval, first, a))
}

func F_concat_51c24458711a1d6106cab433642d0c9c(e, a, op, first, second, third Exp) Exp {
	return A(concat, A(eval, first, a), A(eval, second, a))
}

func F_cond_7b60153e8c9796298806c07f80f3c12e(e, a, op, first, second, third Exp) Exp {
	return A(evcon, A(cdr, e), a)
}

func F_cons_e5efa02a9e5b367e867a09002212f851(e, a, op, first, second, third Exp) Exp {
	return A(cons, A(eval, first, a), A(eval, second, a))
}

func F_display_ad275a4d1320cdde3e44710b8a67ddef(e, a, op, first, second, third Exp) Exp {
	return A(display, A(eval, first, a))
}

func F_eq_d0e57c13aee9aef9ae17b85b67503bf0(e, a, op, first, second, third Exp) Exp {
	return A(eq, A(eval, first, a), A(eval, second, a))
}

func F_err_6bbc47e9027043d4689777651239c2ff(e, a, op, first, second, third Exp) Exp {
	return A(err, A(eval, A(cadr, e), a))
}

func F_exp_af5bcba2f722aebdc27bbc820ffab22f(e, a, op, first, second, third Exp) Exp {
	return A(exp, A(eval, first, a), A(eval, second, a), A(eval, third, a))
}

func F_hash_a2c0339d4437ecca89bf201adb7d1163(e, a, op, first, second, third Exp) Exp {
	return A(hash, A(eval, first, a))
}

func F_inc_c7587605642eb29d0d026da3b2833f89(e, a, op, first, second, third Exp) Exp {
	return A(plus, A(eval, first, a), "1")
}

func F_list_69bc1c031779ce3278497999af69288e(e, a, op, first, second, third Exp) Exp {
	return A(evlis, A(cdr, e), a)
}

func F_mul_71b5055f1965765adf59dc035730df6b(e, a, op, first, second, third Exp) Exp {
	return A(mul, A(eval, first, a), A(eval, second, a))
}

func F_newkey_eb320fe7889a92ac8dbdacd07152a23e(e, a, op, first, second, third Exp) Exp {
	return A(newkey)
}

func F_pub_48401e110dae546272407c5eff0a4a24(e, a, op, first, second, third Exp) Exp {
	return A(pub, A(eval, first, a))
}

func F_quote_9d5418c8b7809b2da600bfc812226bc4(e, a, op, first, second, third Exp) Exp {
	return A(cadr, e)
}

func F_runes_9497265e8a0baaf9c7e9ac783fd5c02b(e, a, op, first, second, third Exp) Exp {
	return A(runes, A(eval, A(cadr, e), a))
}

func F_sign_5a1e6f12da842fac062579e3ff554e4b(e, a, op, first, second, third Exp) Exp {
	return A(sign, A(eval, first, a), A(eval, second, a))
}

func F_sub_c810a170d3009592f391f345d73c221a(e, a, op, first, second, third Exp) Exp {
	return A(sub, A(eval, first, a), A(eval, second, a))
}

func F_test1_1bdc61bba4f3c7d50dc11a03e1c223bb(e, a, op, first, second, third Exp) Exp {
	return A(test1, A(eval, first, a))
}

func F_test2_5feee42882cd3faa875e5e551b346d74(e, a, op, first, second, third Exp) Exp {
	return A(test2, A(eval, first, a))
}

func F_verify_66b0d5b8e697b8cf42702e0edd6a8d16(e, a, op, first, second, third Exp) Exp {
	return A(verify, A(eval, first, a), A(eval, second, a), A(eval, third, a))
}

var map_25 = make(map[string]func(e, a, op, first, second, third Exp) Exp)

func init() {
	map_25 = map[string]func(e, a, op, first, second, third Exp) Exp{
		"test1":   F_test1_1bdc61bba4f3c7d50dc11a03e1c223bb,
		"test2":   F_test2_5feee42882cd3faa875e5e551b346d74,
		"quote":   F_quote_9d5418c8b7809b2da600bfc812226bc4,
		"atom":    F_atom_8f3c75e470915a3502aaee0ca1577fcc,
		"eq":      F_eq_d0e57c13aee9aef9ae17b85b67503bf0,
		"car":     F_car_9ccdcb1ba9d132a44726ab42f27ff4e9,
		"cdr":     F_cdr_07b8b4fee0b4a65b5f5e3e5580aaa311,
		"cons":    F_cons_e5efa02a9e5b367e867a09002212f851,
		"cond":    F_cond_7b60153e8c9796298806c07f80f3c12e,
		"add":     F_add_cb8e93847499008be09106dc15141a43,
		"inc":     F_inc_c7587605642eb29d0d026da3b2833f89,
		"sub":     F_sub_c810a170d3009592f391f345d73c221a,
		"mul":     F_mul_71b5055f1965765adf59dc035730df6b,
		"exp":     F_exp_af5bcba2f722aebdc27bbc820ffab22f,
		"after":   F_after_515cbbd8f947c37b7d1dc2951e924a4d,
		"concat":  F_concat_51c24458711a1d6106cab433642d0c9c,
		"hash":    F_hash_a2c0339d4437ecca89bf201adb7d1163,
		"newkey":  F_newkey_eb320fe7889a92ac8dbdacd07152a23e,
		"pub":     F_pub_48401e110dae546272407c5eff0a4a24,
		"sign":    F_sign_5a1e6f12da842fac062579e3ff554e4b,
		"verify":  F_verify_66b0d5b8e697b8cf42702e0edd6a8d16,
		"display": F_display_ad275a4d1320cdde3e44710b8a67ddef,
		"runes":   F_runes_9497265e8a0baaf9c7e9ac783fd5c02b,
		"err":     F_err_6bbc47e9027043d4689777651239c2ff,
		"list":    F_list_69bc1c031779ce3278497999af69288e,
	}
}

var map_27 = make(map[string]func(x, first, rest Exp) Exp)

func init() {
	map_27 = map[string]func(x, first, rest Exp) Exp{
		"0": F_0_60a3caf220a5e9e8d986965c79b20f11,
	}
}

var map_28 = make(map[string]func(x, first, rest Exp) Exp)

func init() {
	map_28 = map[string]func(x, first, rest Exp) Exp{
		"0": F_0_60a3caf220a5e9e8d986965c79b20f11,
	}
}

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
		L("s", s_label),
		L("subst", subst_label),
		L("tassoc", tassoc_label),
		L("test1", test1_label),
		L("test2", test2_label),
		L("test3", test3_label),
		L("test4", test4_label),
		L("try", try_label),
		L("xlist", xlist_label),
		L("y", y_label),
	)
}
