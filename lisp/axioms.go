package lisp

import "github.com/xoba/turd/lisp/exp"

// #1
func Quote(e exp.Expression) exp.Expression {
	return e
}

// #2
func Atom(e exp.Expression) exp.Expression {
	if e.Atom() != nil {
		return True()
	}
	if len(e.List()) == 0 {
		return True()
	}
	return False()
}

// #3
func Eq(x, y exp.Expression) exp.Expression {
	ret := func(v bool) exp.Expression {
		if v {
			return True()
		}
		return False()
	}
	switch {
	case IsAtom(x) && IsAtom(y):
		return ret(AtomsEqual(x.Atom(), y.Atom()))
	case IsList(x) && IsList(y):
		return ret(IsEmpty(x) && IsEmpty(y))
	default:
		return False()
	}
}

// #4
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

// #5
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

// #6
func Cons(x, y exp.Expression) exp.Expression {
	if !IsList(y) {
		return exp.Errorf("second arg not a list: %s", y)
	}
	var args []exp.Expression
	add := func(e exp.Expression) {
		args = append(args, e)
	}
	add(x)
	for _, e := range y.List() {
		add(e)
	}
	return exp.NewList(args...)
}

// #7
func Cond(args ...exp.Expression) exp.Expression {
	for _, a := range args {
		if Boolean(car(a)) {
			return cadr(a)
		}
	}
	return exp.Errorf("cond fallthrough: %s", args)
}
