// THIS FILE IS AUTOGENERATED, DO NOT EDIT!

package lisp

var env_and = list(quote("and"), list(quote("label"), quote("and"), list(quote("lambda"), list(quote("x"), quote("y")), list(quote("cond"), list(quote("x"), list(quote("cond"), list(quote("y"), list(quote("quote"), quote("t"))), list(list(quote("quote"), quote("t")), list()))), list(list(quote("quote"), quote("t")), list(quote("quote"), list()))))))

func and(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return x
			}),
			Func(func(...Exp) Exp {
				return apply(cond, list(
					Func(func(...Exp) Exp {
						return y
					}),
					Func(func(...Exp) Exp {
						return "t"
					}),
				), list(
					Func(func(...Exp) Exp {
						return "t"
					}),
					Func(func(...Exp) Exp {
						return Nil
					}),
				))
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return "t"
			}),
			Func(func(...Exp) Exp {
				return Nil
			}),
		),
	)
}

var env_go_sanitized_append = list(quote("go_sanitized_append"), list(quote("label"), quote("go_sanitized_append"), list(quote("lambda"), list(quote("x"), quote("y")), list(quote("cond"), list(list(quote("null"), quote("x")), quote("y")), list(list(quote("quote"), quote("t")), list(quote("cons"), list(quote("car"), quote("x")), list(quote("go_sanitized_append"), list(quote("cdr"), quote("x")), quote("y"))))))))

func go_sanitized_append(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return apply(null, x)
			}),
			Func(func(...Exp) Exp {
				return y
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return "t"
			}),
			Func(func(...Exp) Exp {
				return apply(cons, apply(car, x), apply(go_sanitized_append, apply(cdr, x), y))
			}),
		),
	)
}

var env_assoc = list(quote("assoc"), list(quote("label"), quote("assoc"), list(quote("lambda"), list(quote("x"), quote("y")), list(quote("cond"), list(list(quote("eq"), list(quote("caar"), quote("y")), quote("x")), list(quote("cadar"), quote("y"))), list(list(quote("quote"), quote("t")), list(quote("assoc"), quote("x"), list(quote("cdr"), quote("y"))))))))

func assoc(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return apply(eq, apply(caar, y), x)
			}),
			Func(func(...Exp) Exp {
				return apply(cadar, y)
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return "t"
			}),
			Func(func(...Exp) Exp {
				return apply(assoc, x, apply(cdr, y))
			}),
		),
	)
}

var env_caaaar = list(quote("caaaar"), list(quote("label"), quote("caaaar"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("car"), list(quote("car"), list(quote("car"), quote("x"))))))))

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

var env_caaadr = list(quote("caaadr"), list(quote("label"), quote("caaadr"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("car"), list(quote("car"), list(quote("cdr"), quote("x"))))))))

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

var env_caaar = list(quote("caaar"), list(quote("label"), quote("caaar"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("car"), list(quote("car"), quote("x")))))))

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

var env_caadar = list(quote("caadar"), list(quote("label"), quote("caadar"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("car"), list(quote("cdr"), list(quote("car"), quote("x"))))))))

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

var env_caaddr = list(quote("caaddr"), list(quote("label"), quote("caaddr"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("car"), list(quote("cdr"), list(quote("cdr"), quote("x"))))))))

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

var env_caadr = list(quote("caadr"), list(quote("label"), quote("caadr"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("car"), list(quote("cdr"), quote("x")))))))

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

var env_caar = list(quote("caar"), list(quote("label"), quote("caar"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("car"), quote("x"))))))

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

var env_cadaar = list(quote("cadaar"), list(quote("label"), quote("cadaar"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("cdr"), list(quote("car"), list(quote("car"), quote("x"))))))))

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

var env_cadadr = list(quote("cadadr"), list(quote("label"), quote("cadadr"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("cdr"), list(quote("car"), list(quote("cdr"), quote("x"))))))))

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

var env_cadar = list(quote("cadar"), list(quote("label"), quote("cadar"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("cdr"), list(quote("car"), quote("x")))))))

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

var env_caddar = list(quote("caddar"), list(quote("label"), quote("caddar"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("cdr"), list(quote("cdr"), list(quote("car"), quote("x"))))))))

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

var env_cadddr = list(quote("cadddr"), list(quote("label"), quote("cadddr"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("cdr"), list(quote("cdr"), list(quote("cdr"), quote("x"))))))))

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

var env_caddr = list(quote("caddr"), list(quote("label"), quote("caddr"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("cdr"), list(quote("cdr"), quote("x")))))))

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

var env_cadr = list(quote("cadr"), list(quote("label"), quote("cadr"), list(quote("lambda"), list(quote("x")), list(quote("car"), list(quote("cdr"), quote("x"))))))

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

var env_cdaaar = list(quote("cdaaar"), list(quote("label"), quote("cdaaar"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("car"), list(quote("car"), list(quote("car"), quote("x"))))))))

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

var env_cdaadr = list(quote("cdaadr"), list(quote("label"), quote("cdaadr"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("car"), list(quote("car"), list(quote("cdr"), quote("x"))))))))

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

var env_cdaar = list(quote("cdaar"), list(quote("label"), quote("cdaar"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("car"), list(quote("car"), quote("x")))))))

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

var env_cdadar = list(quote("cdadar"), list(quote("label"), quote("cdadar"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("car"), list(quote("cdr"), list(quote("car"), quote("x"))))))))

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

var env_cdaddr = list(quote("cdaddr"), list(quote("label"), quote("cdaddr"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("car"), list(quote("cdr"), list(quote("cdr"), quote("x"))))))))

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

var env_cdadr = list(quote("cdadr"), list(quote("label"), quote("cdadr"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("car"), list(quote("cdr"), quote("x")))))))

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

var env_cdar = list(quote("cdar"), list(quote("label"), quote("cdar"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("car"), quote("x"))))))

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

var env_cddaar = list(quote("cddaar"), list(quote("label"), quote("cddaar"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("cdr"), list(quote("car"), list(quote("car"), quote("x"))))))))

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

var env_cddadr = list(quote("cddadr"), list(quote("label"), quote("cddadr"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("cdr"), list(quote("car"), list(quote("cdr"), quote("x"))))))))

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

var env_cddar = list(quote("cddar"), list(quote("label"), quote("cddar"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("cdr"), list(quote("car"), quote("x")))))))

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

var env_cdddar = list(quote("cdddar"), list(quote("label"), quote("cdddar"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("cdr"), list(quote("cdr"), list(quote("car"), quote("x"))))))))

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

var env_cddddr = list(quote("cddddr"), list(quote("label"), quote("cddddr"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("cdr"), list(quote("cdr"), list(quote("cdr"), quote("x"))))))))

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

var env_cdddr = list(quote("cdddr"), list(quote("label"), quote("cdddr"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("cdr"), list(quote("cdr"), quote("x")))))))

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

var env_cddr = list(quote("cddr"), list(quote("label"), quote("cddr"), list(quote("lambda"), list(quote("x")), list(quote("cdr"), list(quote("cdr"), quote("x"))))))

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

var env_eval = list(quote("eval"), list(quote("label"), quote("eval"), list(quote("lambda"), list(quote("e"), quote("a")), list(quote("cond"), list(list(quote("atom"), quote("e")), list(quote("assoc"), quote("e"), quote("a"))), list(list(quote("atom"), list(quote("car"), quote("e"))), list(quote("cond"), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("quote"))), list(quote("cadr"), quote("e"))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("atom"))), list(quote("atom"), list(quote("eval"), list(quote("cadr"), quote("e")), quote("a")))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("eq"))), list(quote("eq"), list(quote("eval"), list(quote("cadr"), quote("e")), quote("a")), list(quote("eval"), list(quote("caddr"), quote("e")), quote("a")))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("plus"))), list(quote("plus"), list(quote("eval"), list(quote("cadr"), quote("e")), quote("a")), list(quote("eval"), list(quote("caddr"), quote("e")), quote("a")))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("minus"))), list(quote("minus"), list(quote("eval"), list(quote("cadr"), quote("e")), quote("a")), list(quote("eval"), list(quote("caddr"), quote("e")), quote("a")))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("mult"))), list(quote("mult"), list(quote("eval"), list(quote("cadr"), quote("e")), quote("a")), list(quote("eval"), list(quote("caddr"), quote("e")), quote("a")))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("display"))), list(quote("display"), list(quote("eval"), list(quote("cadr"), quote("e")), quote("a")))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("car"))), list(quote("car"), list(quote("eval"), list(quote("cadr"), quote("e")), quote("a")))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("cdr"))), list(quote("cdr"), list(quote("eval"), list(quote("cadr"), quote("e")), quote("a")))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("cons"))), list(quote("cons"), list(quote("eval"), list(quote("cadr"), quote("e")), quote("a")), list(quote("eval"), list(quote("caddr"), quote("e")), quote("a")))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("cond"))), list(quote("evcon"), list(quote("cdr"), quote("e")), quote("a"))), list(list(quote("eq"), list(quote("car"), quote("e")), list(quote("quote"), quote("list"))), list(quote("evlis"), list(quote("cdr"), quote("e")), quote("a"))), list(list(quote("quote"), quote("t")), list(quote("eval"), list(quote("cons"), list(quote("assoc"), list(quote("car"), quote("e")), quote("a")), list(quote("cdr"), quote("e"))), quote("a"))))), list(list(quote("eq"), list(quote("caar"), quote("e")), list(quote("quote"), quote("label"))), list(quote("eval"), list(quote("cons"), list(quote("caddar"), quote("e")), list(quote("cdr"), quote("e"))), list(quote("cons"), list(quote("list"), list(quote("cadar"), quote("e")), list(quote("car"), quote("e"))), quote("a")))), list(list(quote("eq"), list(quote("caar"), quote("e")), list(quote("quote"), quote("lambda"))), list(quote("eval"), list(quote("caddar"), quote("e")), list(quote("go_sanitized_append"), list(quote("pair"), list(quote("cadar"), quote("e")), list(quote("evlis"), list(quote("cdr"), quote("e")), quote("a"))), quote("a"))))))))

func eval(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	e := args[0]
	a := args[1]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return apply(atom, e)
			}),
			Func(func(...Exp) Exp {
				return apply(assoc, e, a)
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return apply(atom, apply(car, e))
			}),
			Func(func(...Exp) Exp {
				return apply(cond, list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "quote")
					}),
					Func(func(...Exp) Exp {
						return apply(cadr, e)
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "atom")
					}),
					Func(func(...Exp) Exp {
						return apply(atom, apply(eval, apply(cadr, e), a))
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "eq")
					}),
					Func(func(...Exp) Exp {
						return apply(eq, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "plus")
					}),
					Func(func(...Exp) Exp {
						return apply(plus, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "minus")
					}),
					Func(func(...Exp) Exp {
						return apply(minus, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "mult")
					}),
					Func(func(...Exp) Exp {
						return apply(mult, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "display")
					}),
					Func(func(...Exp) Exp {
						return apply(display, apply(eval, apply(cadr, e), a))
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "car")
					}),
					Func(func(...Exp) Exp {
						return apply(car, apply(eval, apply(cadr, e), a))
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "cdr")
					}),
					Func(func(...Exp) Exp {
						return apply(cdr, apply(eval, apply(cadr, e), a))
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "cons")
					}),
					Func(func(...Exp) Exp {
						return apply(cons, apply(eval, apply(cadr, e), a), apply(eval, apply(caddr, e), a))
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "cond")
					}),
					Func(func(...Exp) Exp {
						return apply(evcon, apply(cdr, e), a)
					}),
				), list(
					Func(func(...Exp) Exp {
						return apply(eq, apply(car, e), "list")
					}),
					Func(func(...Exp) Exp {
						return apply(evlis, apply(cdr, e), a)
					}),
				), list(
					Func(func(...Exp) Exp {
						return "t"
					}),
					Func(func(...Exp) Exp {
						return apply(eval, apply(cons, apply(assoc, apply(car, e), a), apply(cdr, e)), a)
					}),
				))
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return apply(eq, apply(caar, e), "label")
			}),
			Func(func(...Exp) Exp {
				return apply(eval, apply(cons, apply(caddar, e), apply(cdr, e)), apply(cons, apply(list, apply(cadar, e), apply(car, e)), a))
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return apply(eq, apply(caar, e), "lambda")
			}),
			Func(func(...Exp) Exp {
				return apply(eval, apply(caddar, e), apply(go_sanitized_append, apply(pair, apply(cadar, e), apply(evlis, apply(cdr, e), a)), a))
			}),
		),
	)
}

var env_evcon = list(quote("evcon"), list(quote("label"), quote("evcon"), list(quote("lambda"), list(quote("c"), quote("a")), list(quote("cond"), list(list(quote("eval"), list(quote("caar"), quote("c")), quote("a")), list(quote("eval"), list(quote("cadar"), quote("c")), quote("a"))), list(list(quote("quote"), quote("t")), list(quote("evcon"), list(quote("cdr"), quote("c")), quote("a")))))))

func evcon(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	c := args[0]
	a := args[1]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return apply(eval, apply(caar, c), a)
			}),
			Func(func(...Exp) Exp {
				return apply(eval, apply(cadar, c), a)
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return "t"
			}),
			Func(func(...Exp) Exp {
				return apply(evcon, apply(cdr, c), a)
			}),
		),
	)
}

var env_evlis = list(quote("evlis"), list(quote("label"), quote("evlis"), list(quote("lambda"), list(quote("m"), quote("a")), list(quote("cond"), list(list(quote("null"), quote("m")), list(quote("quote"), list())), list(list(quote("quote"), quote("t")), list(quote("cons"), list(quote("eval"), list(quote("car"), quote("m")), quote("a")), list(quote("evlis"), list(quote("cdr"), quote("m")), quote("a"))))))))

func evlis(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	m := args[0]
	a := args[1]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return apply(null, m)
			}),
			Func(func(...Exp) Exp {
				return Nil
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return "t"
			}),
			Func(func(...Exp) Exp {
				return apply(cons, apply(eval, apply(car, m), a), apply(evlis, apply(cdr, m), a))
			}),
		),
	)
}

var env_factorial = list(quote("factorial"), list(quote("label"), quote("factorial"), list(quote("lambda"), list(quote("n")), list(quote("cond"), list(list(quote("eq"), list(quote("quote"), quote("0")), quote("n")), list(quote("quote"), quote("1"))), list(list(quote("quote"), quote("t")), list(quote("mult"), quote("n"), list(quote("factorial"), list(quote("minus"), quote("n"), list(quote("quote"), quote("1"))))))))))

func factorial(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	n := args[0]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return apply(eq, "0", n)
			}),
			Func(func(...Exp) Exp {
				return "1"
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return "t"
			}),
			Func(func(...Exp) Exp {
				return apply(mult, n, apply(factorial, apply(minus, n, "1")))
			}),
		),
	)
}

var env_not = list(quote("not"), list(quote("label"), quote("not"), list(quote("lambda"), list(quote("x")), list(quote("cond"), list(quote("x"), list(quote("quote"), list())), list(list(quote("quote"), quote("t")), list(quote("quote"), quote("t")))))))

func not(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return x
			}),
			Func(func(...Exp) Exp {
				return Nil
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return "t"
			}),
			Func(func(...Exp) Exp {
				return "t"
			}),
		),
	)
}

var env_null = list(quote("null"), list(quote("label"), quote("null"), list(quote("lambda"), list(quote("x")), list(quote("eq"), quote("x"), list(quote("quote"), list())))))

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

var env_pair = list(quote("pair"), list(quote("label"), quote("pair"), list(quote("lambda"), list(quote("x"), quote("y")), list(quote("cond"), list(list(quote("and"), list(quote("null"), quote("x")), list(quote("null"), quote("y"))), list(quote("quote"), list())), list(list(quote("and"), list(quote("not"), list(quote("atom"), quote("x"))), list(quote("not"), list(quote("atom"), quote("y")))), list(quote("cons"), list(quote("list"), list(quote("car"), quote("x")), list(quote("car"), quote("y"))), list(quote("pair"), list(quote("cdr"), quote("x")), list(quote("cdr"), quote("y")))))))))

func pair(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return apply(and, apply(null, x), apply(null, y))
			}),
			Func(func(...Exp) Exp {
				return Nil
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return apply(and, apply(not, apply(atom, x)), apply(not, apply(atom, y)))
			}),
			Func(func(...Exp) Exp {
				return apply(cons, apply(list, apply(car, x), apply(car, y)), apply(pair, apply(cdr, x), apply(cdr, y)))
			}),
		),
	)
}

var env_subst = list(quote("subst"), list(quote("label"), quote("subst"), list(quote("lambda"), list(quote("x"), quote("y"), quote("z")), list(quote("cond"), list(list(quote("atom"), quote("z")), list(quote("cond"), list(list(quote("eq"), quote("z"), quote("y")), quote("x")), list(list(quote("quote"), quote("t")), quote("z")))), list(list(quote("quote"), quote("t")), list(quote("cons"), list(quote("subst"), quote("x"), quote("y"), list(quote("car"), quote("z"))), list(quote("subst"), quote("x"), quote("y"), list(quote("cdr"), quote("z")))))))))

func subst(args ...Exp) Exp {
	if err := checklen(3, args); err != nil {
		return err
	}
	x := args[0]
	y := args[1]
	z := args[2]
	return apply(
		cond,
		list(
			Func(func(...Exp) Exp {
				return apply(atom, z)
			}),
			Func(func(...Exp) Exp {
				return apply(cond, list(
					Func(func(...Exp) Exp {
						return apply(eq, z, y)
					}),
					Func(func(...Exp) Exp {
						return x
					}),
				), list(
					Func(func(...Exp) Exp {
						return "t"
					}),
					Func(func(...Exp) Exp {
						return z
					}),
				))
			}),
		),
		list(
			Func(func(...Exp) Exp {
				return "t"
			}),
			Func(func(...Exp) Exp {
				return apply(cons, apply(subst, x, y, apply(car, z)), apply(subst, x, y, apply(cdr, z)))
			}),
		),
	)
}

var env_testing = list(quote("testing"), list(quote("label"), quote("testing"), list(quote("lambda"), list(quote("x")), list(quote("display"), list(quote("car"), quote("x"))))))

func testing(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	return apply(
		display,
		apply(
			car,
			x,
		),
	)
}

func init() {
	env = list(env_and, env_go_sanitized_append, env_assoc, env_caaaar, env_caaadr, env_caaar, env_caadar, env_caaddr, env_caadr, env_caar, env_cadaar, env_cadadr, env_cadar, env_caddar, env_cadddr, env_caddr, env_cadr, env_cdaaar, env_cdaadr, env_cdaar, env_cdadar, env_cdaddr, env_cdadr, env_cdar, env_cddaar, env_cddadr, env_cddar, env_cdddar, env_cddddr, env_cdddr, env_cddr, env_eval, env_evcon, env_evlis, env_factorial, env_not, env_null, env_pair, env_subst, env_testing)
}
