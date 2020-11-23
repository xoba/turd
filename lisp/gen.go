// THIS FILE IS AUTOGENERATED, DO NOT EDIT!

package lisp

var L = list

var and_label = L("label", "and", L("lambda", L("x", "y"), L("cond", L("x", L("cond", L("y", L("quote", "t")), L(L("quote", "t"), L()))), L(L("quote", "t"), L("quote", L())))))

func and(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	return apply(
		cond,
		[]Exp{
			x,
			Func(func(...Exp) Exp {
				return apply(cond, []Exp{
					y,
					"t",
				}, []Exp{
					"t",
					Nil,
				})
			}),
		},
		[]Exp{
			"t",
			Nil,
		},
	)
}

var append_go_sanitized_label = L("label", "append_go_sanitized", L("lambda", L("x", "y"), L("cond", L(L("null", "x"), "y"), L(L("quote", "t"), L("cons", L("car", "x"), L("append_go_sanitized", L("cdr", "x"), "y"))))))

func append_go_sanitized(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	return apply(
		cond,
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(null, x)
			}),
			y,
		},
		[]Exp{
			"t",
			Func(func(...Exp) Exp {
				return apply(cons, apply(car, x), apply(append_go_sanitized, apply(cdr, x), y))
			}),
		},
	)
}

var assoc_label = L("label", "assoc", L("lambda", L("x", "y"), L("cond", L(L("eq", L("caar", "y"), "x"), L("cadar", "y")), L(L("quote", "t"), L("assoc", "x", L("cdr", "y"))))))

func assoc(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	return apply(
		cond,
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(eq, apply(caar, y), x)
			}),
			Func(func(...Exp) Exp {
				return apply(cadar, y)
			}),
		},
		[]Exp{
			"t",
			Func(func(...Exp) Exp {
				return apply(assoc, x, apply(cdr, y))
			}),
		},
	)
}

var caaaar_label = L("label", "caaaar", L("lambda", L("x"), L("car", L("car", L("car", L("car", "x"))))))

func caaaar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			car,
			apply(
				car,
				apply(
					car,
					x,
				),
			),
		),
	)
}

var caaadr_label = L("label", "caaadr", L("lambda", L("x"), L("car", L("car", L("car", L("cdr", "x"))))))

func caaadr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			car,
			apply(
				car,
				apply(
					cdr,
					x,
				),
			),
		),
	)
}

var caaar_label = L("label", "caaar", L("lambda", L("x"), L("car", L("car", L("car", "x")))))

func caaar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			car,
			apply(
				car,
				x,
			),
		),
	)
}

var caadar_label = L("label", "caadar", L("lambda", L("x"), L("car", L("car", L("cdr", L("car", "x"))))))

func caadar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			car,
			apply(
				cdr,
				apply(
					car,
					x,
				),
			),
		),
	)
}

var caaddr_label = L("label", "caaddr", L("lambda", L("x"), L("car", L("car", L("cdr", L("cdr", "x"))))))

func caaddr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			car,
			apply(
				cdr,
				apply(
					cdr,
					x,
				),
			),
		),
	)
}

var caadr_label = L("label", "caadr", L("lambda", L("x"), L("car", L("car", L("cdr", "x")))))

func caadr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			car,
			apply(
				cdr,
				x,
			),
		),
	)
}

var caar_label = L("label", "caar", L("lambda", L("x"), L("car", L("car", "x"))))

func caar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			car,
			x,
		),
	)
}

var cadaar_label = L("label", "cadaar", L("lambda", L("x"), L("car", L("cdr", L("car", L("car", "x"))))))

func cadaar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			cdr,
			apply(
				car,
				apply(
					car,
					x,
				),
			),
		),
	)
}

var cadadr_label = L("label", "cadadr", L("lambda", L("x"), L("car", L("cdr", L("car", L("cdr", "x"))))))

func cadadr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			cdr,
			apply(
				car,
				apply(
					cdr,
					x,
				),
			),
		),
	)
}

var cadar_label = L("label", "cadar", L("lambda", L("x"), L("car", L("cdr", L("car", "x")))))

func cadar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			cdr,
			apply(
				car,
				x,
			),
		),
	)
}

var caddar_label = L("label", "caddar", L("lambda", L("x"), L("car", L("cdr", L("cdr", L("car", "x"))))))

func caddar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			cdr,
			apply(
				cdr,
				apply(
					car,
					x,
				),
			),
		),
	)
}

var cadddr_label = L("label", "cadddr", L("lambda", L("x"), L("car", L("cdr", L("cdr", L("cdr", "x"))))))

func cadddr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			cdr,
			apply(
				cdr,
				apply(
					cdr,
					x,
				),
			),
		),
	)
}

var caddr_label = L("label", "caddr", L("lambda", L("x"), L("car", L("cdr", L("cdr", "x")))))

func caddr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			cdr,
			apply(
				cdr,
				x,
			),
		),
	)
}

var cadr_label = L("label", "cadr", L("lambda", L("x"), L("car", L("cdr", "x"))))

func cadr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		car,
		apply(
			cdr,
			x,
		),
	)
}

var cdaaar_label = L("label", "cdaaar", L("lambda", L("x"), L("cdr", L("car", L("car", L("car", "x"))))))

func cdaaar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			car,
			apply(
				car,
				apply(
					car,
					x,
				),
			),
		),
	)
}

var cdaadr_label = L("label", "cdaadr", L("lambda", L("x"), L("cdr", L("car", L("car", L("cdr", "x"))))))

func cdaadr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			car,
			apply(
				car,
				apply(
					cdr,
					x,
				),
			),
		),
	)
}

var cdaar_label = L("label", "cdaar", L("lambda", L("x"), L("cdr", L("car", L("car", "x")))))

func cdaar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			car,
			apply(
				car,
				x,
			),
		),
	)
}

var cdadar_label = L("label", "cdadar", L("lambda", L("x"), L("cdr", L("car", L("cdr", L("car", "x"))))))

func cdadar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			car,
			apply(
				cdr,
				apply(
					car,
					x,
				),
			),
		),
	)
}

var cdaddr_label = L("label", "cdaddr", L("lambda", L("x"), L("cdr", L("car", L("cdr", L("cdr", "x"))))))

func cdaddr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			car,
			apply(
				cdr,
				apply(
					cdr,
					x,
				),
			),
		),
	)
}

var cdadr_label = L("label", "cdadr", L("lambda", L("x"), L("cdr", L("car", L("cdr", "x")))))

func cdadr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			car,
			apply(
				cdr,
				x,
			),
		),
	)
}

var cdar_label = L("label", "cdar", L("lambda", L("x"), L("cdr", L("car", "x"))))

func cdar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			car,
			x,
		),
	)
}

var cddaar_label = L("label", "cddaar", L("lambda", L("x"), L("cdr", L("cdr", L("car", L("car", "x"))))))

func cddaar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			cdr,
			apply(
				car,
				apply(
					car,
					x,
				),
			),
		),
	)
}

var cddadr_label = L("label", "cddadr", L("lambda", L("x"), L("cdr", L("cdr", L("car", L("cdr", "x"))))))

func cddadr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			cdr,
			apply(
				car,
				apply(
					cdr,
					x,
				),
			),
		),
	)
}

var cddar_label = L("label", "cddar", L("lambda", L("x"), L("cdr", L("cdr", L("car", "x")))))

func cddar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			cdr,
			apply(
				car,
				x,
			),
		),
	)
}

var cdddar_label = L("label", "cdddar", L("lambda", L("x"), L("cdr", L("cdr", L("cdr", L("car", "x"))))))

func cdddar(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			cdr,
			apply(
				cdr,
				apply(
					car,
					x,
				),
			),
		),
	)
}

var cddddr_label = L("label", "cddddr", L("lambda", L("x"), L("cdr", L("cdr", L("cdr", L("cdr", "x"))))))

func cddddr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			cdr,
			apply(
				cdr,
				apply(
					cdr,
					x,
				),
			),
		),
	)
}

var cdddr_label = L("label", "cdddr", L("lambda", L("x"), L("cdr", L("cdr", L("cdr", "x")))))

func cdddr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			cdr,
			apply(
				cdr,
				x,
			),
		),
	)
}

var cddr_label = L("label", "cddr", L("lambda", L("x"), L("cdr", L("cdr", "x"))))

func cddr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cdr,
		apply(
			cdr,
			x,
		),
	)
}

var eval_label = L("label", "eval", L("lambda", L("e", "a"), L("cond", L(L("atom", "e"), L("assoc", "e", "a")), L(L("atom", L("car", "e")), L("cond", L(L("eq", L("car", "e"), L("quote", "quote")), L("cadr", "e")), L(L("eq", L("car", "e"), L("quote", "atom")), L("atom", L("eval", L("cadr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "eq")), L("eq", L("eval", L("cadr", "e"), "a"), L("eval", L("caddr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "plus")), L("plus", L("eval", L("cadr", "e"), "a"), L("eval", L("caddr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "minus")), L("minus", L("eval", L("cadr", "e"), "a"), L("eval", L("caddr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "mult")), L("mult", L("eval", L("cadr", "e"), "a"), L("eval", L("caddr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "after")), L("after", L("eval", L("cadr", "e"), "a"), L("eval", L("caddr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "display")), L("display", L("eval", L("cadr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "concat")), L("concat", L("eval", L("cadr", "e"), "a"), L("eval", L("caddr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "hash")), L("hash", L("eval", L("cadr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "newkey")), L("newkey")), L(L("eq", L("car", "e"), L("quote", "pub")), L("pub", L("eval", L("cadr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "sign")), L("sign", L("eval", L("cadr", "e"), "a"), L("eval", L("caddr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "verify")), L("verify", L("eval", L("cadr", "e"), "a"), L("eval", L("caddr", "e"), "a"), L("eval", L("cadddr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "car")), L("car", L("eval", L("cadr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "cdr")), L("cdr", L("eval", L("cadr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "cons")), L("cons", L("eval", L("cadr", "e"), "a"), L("eval", L("caddr", "e"), "a"))), L(L("eq", L("car", "e"), L("quote", "cond")), L("evcon", L("cdr", "e"), "a")), L(L("eq", L("car", "e"), L("quote", "list")), L("evlis", L("cdr", "e"), "a")), L(L("quote", "t"), L("eval", L("cons", L("assoc", L("car", "e"), "a"), L("cdr", "e")), "a")))), L(L("eq", L("caar", "e"), L("quote", "label")), L("eval", L("cons", L("caddar", "e"), L("cdr", "e")), L("cons", L("list", L("cadar", "e"), L("car", "e")), "a"))), L(L("eq", L("caar", "e"), L("quote", "lambda")), L("cond", L(L("atom", L("cadar", "e")), L("eval", L("caddar", "e"), L("cons", L("list", L("cadar", "e"), L("evlis", L("cdr", "e"), "a")), "a"))), L(L("quote", "t"), L("eval", L("caddar", "e"), L("append_go_sanitized", L("pair", L("cadar", "e"), L("evlis", L("cdr", "e"), "a")), "a"))))))))

func eval(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	e := args[0]
	a := args[1]
	return apply(
		cond,
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(atom, e)
			}),
			Func(func(...Exp) Exp {
				return apply(assoc, e, a)
			}),
		},
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(atom, apply(car, e))
			}),
			Func(func(...Exp) Exp {
				return apply(cond, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "quote")
					}),
					Func(func(...Exp) Exp {
						return apply(cadr, e)
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "atom")
					}),
					Func(func(...Exp) Exp {
						return apply(atom, apply(eval, apply(cadr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "eq")
					}),
					Func(func(...Exp) Exp {
						return apply(eq, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "plus")
					}),
					Func(func(...Exp) Exp {
						return apply(plus, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "minus")
					}),
					Func(func(...Exp) Exp {
						return apply(minus, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "mult")
					}),
					Func(func(...Exp) Exp {
						return apply(mult, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "after")
					}),
					Func(func(...Exp) Exp {
						return apply(after, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "display")
					}),
					Func(func(...Exp) Exp {
						return apply(display, apply(eval, apply(cadr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "concat")
					}),
					Func(func(...Exp) Exp {
						return apply(concat, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "hash")
					}),
					Func(func(...Exp) Exp {
						return apply(hash, apply(eval, apply(cadr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "newkey")
					}),
					Func(func(...Exp) Exp {
						return apply(newkey)
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "pub")
					}),
					Func(func(...Exp) Exp {
						return apply(pub, apply(eval, apply(cadr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "sign")
					}),
					Func(func(...Exp) Exp {
						return apply(sign, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "verify")
					}),
					Func(func(...Exp) Exp {
						return apply(verify, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a), apply(eval, apply(cadddr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "car")
					}),
					Func(func(...Exp) Exp {
						return apply(car, apply(eval, apply(cadr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "cdr")
					}),
					Func(func(...Exp) Exp {
						return apply(cdr, apply(eval, apply(cadr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "cons")
					}),
					Func(func(...Exp) Exp {
						return apply(cons, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "cond")
					}),
					Func(func(...Exp) Exp {
						return apply(evcon, apply(cdr, e), a)
					}),
				}, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "list")
					}),
					Func(func(...Exp) Exp {
						return apply(evlis, apply(cdr, e), a)
					}),
				}, []Exp{
					"t",
					Func(func(...Exp) Exp {
						return apply(eval, apply(cons, apply(assoc, apply(car, e), a), apply(cdr, e)), a)
					}),
				})
			}),
		},
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(eq, apply(caar, e), "label")
			}),
			Func(func(...Exp) Exp {
				return apply(eval, apply(cons, apply(caddar, e), apply(cdr, e)), apply(cons, apply(list, apply(cadar, e), apply(car, e)), a))
			}),
		},
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(eq, apply(caar, e), "lambda")
			}),
			Func(func(...Exp) Exp {
				return apply(cond, []Exp{
					Func(func(...Exp) Exp {
						return apply(atom, apply(cadar, e))
					}),
					Func(func(...Exp) Exp {
						return apply(eval, apply(caddar, e), apply(cons, apply(list, apply(cadar, e), apply(evlis, apply(cdr, e), a)), a))
					}),
				}, []Exp{
					"t",
					Func(func(...Exp) Exp {
						return apply(eval, apply(caddar, e), apply(append_go_sanitized, apply(pair, apply(cadar, e), apply(evlis, apply(cdr, e), a)), a))
					}),
				})
			}),
		},
	)
}

var evcon_label = L("label", "evcon", L("lambda", L("c", "a"), L("cond", L(L("eval", L("caar", "c"), "a"), L("eval", L("cadar", "c"), "a")), L(L("quote", "t"), L("evcon", L("cdr", "c"), "a")))))

func evcon(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	c := args[0]
	a := args[1]
	return apply(
		cond,
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(eval, apply(caar, c), a)
			}),
			Func(func(...Exp) Exp {
				return apply(eval, apply(cadar, c), a)
			}),
		},
		[]Exp{
			"t",
			Func(func(...Exp) Exp {
				return apply(evcon, apply(cdr, c), a)
			}),
		},
	)
}

var evlis_label = L("label", "evlis", L("lambda", L("m", "a"), L("cond", L(L("null", "m"), L("quote", L())), L(L("quote", "t"), L("cons", L("eval", L("car", "m"), "a"), L("evlis", L("cdr", "m"), "a"))))))

func evlis(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	m := args[0]
	a := args[1]
	return apply(
		cond,
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(null, m)
			}),
			Nil,
		},
		[]Exp{
			"t",
			Func(func(...Exp) Exp {
				return apply(cons, apply(eval, apply(car, m), a), apply(evlis, apply(cdr, m), a))
			}),
		},
	)
}

var factorial_label = L("label", "factorial", L("lambda", L("n"), L("cond", L(L("eq", L("quote", "0"), "n"), L("quote", "1")), L(L("quote", "t"), L("mult", "n", L("factorial", L("minus", "n", L("quote", "1"))))))))

func factorial(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	n := args[0]
	return apply(
		cond,
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(eq, "0", n)
			}),
			"1",
		},
		[]Exp{
			"t",
			Func(func(...Exp) Exp {
				return apply(mult, n, apply(factorial, apply(minus, n, "1")))
			}),
		},
	)
}

var length_label = L("label", "length", L("lambda", L("x"), L("cond", L(L("atom", "x"), L("quote", "0")), L(L("quote", "t"), L("plus", L("quote", "1"), L("length", L("cdr", "x")))))))

func length(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cond,
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(atom, x)
			}),
			"0",
		},
		[]Exp{
			"t",
			Func(func(...Exp) Exp {
				return apply(plus, "1", apply(length, apply(cdr, x)))
			}),
		},
	)
}

var not_label = L("label", "not", L("lambda", L("x"), L("cond", L("x", L("quote", L())), L(L("quote", "t"), L("quote", "t")))))

func not(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cond,
		[]Exp{
			x,
			Nil,
		},
		[]Exp{
			"t",
			"t",
		},
	)
}

var null_label = L("label", "null", L("lambda", L("x"), L("eq", "x", L("quote", L()))))

func null(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		eq,
		x,
		Nil,
	)
}

var pair_label = L("label", "pair", L("lambda", L("x", "y"), L("cond", L(L("and", L("null", "x"), L("null", "y")), L("quote", L())), L(L("and", L("not", L("atom", "x")), L("not", L("atom", "y"))), L("cons", L("list", L("car", "x"), L("car", "y")), L("pair", L("cdr", "x"), L("cdr", "y")))))))

func pair(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	return apply(
		cond,
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(and, apply(null, x), apply(null, y))
			}),
			Nil,
		},
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(and, apply(not, apply(atom, x)), apply(not, apply(atom, y)))
			}),
			Func(func(...Exp) Exp {
				return apply(cons, apply(list, apply(car, x), apply(car, y)), apply(pair, apply(cdr, x), apply(cdr, y)))
			}),
		},
	)
}

var subst_label = L("label", "subst", L("lambda", L("x", "y", "z"), L("cond", L(L("atom", "z"), L("cond", L(L("eq", "z", "y"), "x"), L(L("quote", "t"), "z"))), L(L("quote", "t"), L("cons", L("subst", "x", "y", L("car", "z")), L("subst", "x", "y", L("cdr", "z")))))))

func subst(args ...Exp) Exp {
	if err := checklen(3, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	z := args[2]
	return apply(
		cond,
		[]Exp{
			Func(func(...Exp) Exp {
				return apply(atom, z)
			}),
			Func(func(...Exp) Exp {
				return apply(cond, []Exp{
					Func(func(...Exp) Exp {
						return apply(eq, z, y)
					}),
					x,
				}, []Exp{
					"t",
					z,
				})
			}),
		},
		[]Exp{
			"t",
			Func(func(...Exp) Exp {
				return apply(cons, apply(subst, x, y, apply(car, z)), apply(subst, x, y, apply(cdr, z)))
			}),
		},
	)
}

var xlist_label = L("label", "xlist", L("lambda", "x", "x"))

func xlist(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return x
}

func init() {
	env = L(
		L("and", and_label),
		L("append_go_sanitized", append_go_sanitized_label),
		L("assoc", assoc_label),
		L("caaaar", caaaar_label),
		L("caaadr", caaadr_label),
		L("caaar", caaar_label),
		L("caadar", caadar_label),
		L("caaddr", caaddr_label),
		L("caadr", caadr_label),
		L("caar", caar_label),
		L("cadaar", cadaar_label),
		L("cadadr", cadadr_label),
		L("cadar", cadar_label),
		L("caddar", caddar_label),
		L("cadddr", cadddr_label),
		L("caddr", caddr_label),
		L("cadr", cadr_label),
		L("cdaaar", cdaaar_label),
		L("cdaadr", cdaadr_label),
		L("cdaar", cdaar_label),
		L("cdadar", cdadar_label),
		L("cdaddr", cdaddr_label),
		L("cdadr", cdadr_label),
		L("cdar", cdar_label),
		L("cddaar", cddaar_label),
		L("cddadr", cddadr_label),
		L("cddar", cddar_label),
		L("cdddar", cdddar_label),
		L("cddddr", cddddr_label),
		L("cdddr", cdddr_label),
		L("cddr", cddr_label),
		L("eval", eval_label),
		L("evcon", evcon_label),
		L("evlis", evlis_label),
		L("factorial", factorial_label),
		L("length", length_label),
		L("not", not_label),
		L("null", null_label),
		L("pair", pair_label),
		L("subst", subst_label),
		L("xlist", xlist_label),
	)
}
