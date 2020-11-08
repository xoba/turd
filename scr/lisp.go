package scr

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/xoba/turd/cnfg"
)

func Lisp(cnfg.Config) error {
	m := make(map[string]bool)
	test := func(i int, in, expect string) {
		wrap := func(e error) error {
			if e == nil {
				return nil
			}
			return fmt.Errorf("#%d. %w", i, e)
		}
		if in == "" {
			return
		}
		if m[in] {
			panic("duplication: " + in)
		}
		m[in] = true
		fmt.Printf("%2d. %-50s -> %s\n", i, in, expect)
		check := func(e error) {
			if e != nil {
				log.Fatal(e)
			}
		}
		e, err := Read(in)
		check(wrap(err))
		a := NewList()
		x, err := Eval(e, a)
		check(wrap(err))
		if got := x.String(); got != expect {
			check(wrap(fmt.Errorf("expected %q, got %q", expect, got)))
		}
	}

	test2 := func(x, y string) {
		test(0, x, y)
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
	test(7, "(cond ((eq 'a 'b) 'first) ((atom 'a) 'second))", "second")

	test(0, "", "")
	test(0, "", "")
	test(0, "", "")
	test(0, "", "")
	test(0, "", "")
	test(0, "", "")

	return nil
}

func Assoc(x, y *Expression) (*Expression, error) {
	caar, err := Caar(y)
	if err != nil {
		return nil, err
	}
	eq, err := Eq(caar, x)
	if err != nil {
		return nil, err
	}
	if eq.Boolean() {
		return Cadar(y)
	}
	cdr, err := Cdr(y)
	if err != nil {
		return nil, err
	}
	return Assoc(x, cdr)
}

func Eq(x, y *Expression) (*Expression, error) {
	r := func(v bool) (*Expression, error) {
		if v {
			return NewString("t"), nil
		}
		return NewList(), nil
	}
	switch {
	case x.IsAtom() && y.IsAtom():
		return r(AtomsEqual(x.Atom, y.Atom))
	case x.IsList() && y.IsList():
		return r(x.IsEmpty() && y.IsEmpty())
	default:
		return r(false)
	}
}

func Eval(e, a *Expression) (*Expression, error) {
	f := MonadFunc(MEval).ToEvalFunc()
	return f(e, a)
}

func MEval(args ...Maybe) Maybe {
	e, a := args[0], args[1]
	car := Eval1Func(Car).ToMonad()
	cdr := Eval1Func(Cdr).ToMonad()
	cadr := Compose(car, cdr)
	caddr := Compose(car, cdr, cdr)
	eq := Eval2Func(Eq).ToMonad()
	eval := Eval2Func(Eval).ToMonad()
	atom := Eval1Func(AtomF).ToMonad()
	assoc := Eval2Func(Assoc).ToMonad()
	cons := Eval2Func(Cons).ToMonad()

	var evcon MonadFunc
	if false {
		evcon = func(args ...Maybe) Maybe {
			return EvconM(args[0], args[1])
		}
	} else {
		evcon = Eval2Func(Evcon).ToMonad()
	}

	if e.Error != nil {
		return e
	}

	if e.IsAtom() {
		return assoc(e, a)
	}

	care := car(e)

	switch {
	case care.IsAtom():
		switch care.String() {
		case "quote":
			return cadr(e)
		case "atom":
			return atom(eval(cadr(e), a))
		case "eq":
			return eq(
				eval(cadr(e), a),
				eval(caddr(e), a),
			)
		case "car":
			return car(eval(cadr(e), a))
		case "cdr":
			return cdr(eval(cadr(e), a))
		case "cons":
			return cons(
				eval(cadr(e), a),
				eval(caddr(e), a),
			)
		case "cond":
			return evcon(cdr(e), a)
		default:
			return eval(cons(assoc(car(e), a), cdr(e)), a)
		}
	default:
		return Maybe{Error: fmt.Errorf("eval unimplemented for %q", args)}
	}
}

func Cons(x, y *Expression) (*Expression, error) {
	if !x.IsAtom() {
		return nil, fmt.Errorf("first arg not an atom: %s", x)
	}
	if !y.IsList() {
		return nil, fmt.Errorf("second arg not a list: %s", y)
	}
	var args []*Expression
	add := func(e *Expression) {
		args = append(args, e)
	}
	add(x)
	for _, e := range *y.List {
		add(e)
	}
	return NewList(args...), nil
}

func QuoteAtom(s string) *Expression {
	return Quote(NewString(s))
}

func Quote(e *Expression) *Expression {
	return NewList(NewString("quote"), e)
}

func EvconM(c, a Maybe) Maybe {
	fmt.Printf("evcon(%q, %q)\n", c, a)
	car := Eval1Func(Car).ToMonad()
	cdr := Eval1Func(Cdr).ToMonad()
	caar := Compose(car, car)
	cadar := Compose(car, cdr, car)
	eval := Eval2Func(Eval).ToMonad()
	cond := EvalFunc(Cond).ToMonad()
	list := func(args ...Maybe) Maybe {
		var list []*Expression
		for _, a := range args {
			if a.Error != nil {
				return a
			}
		}
		return Maybe{Expression: NewList(list...)}
	}
	quote := func(m *Expression) Maybe {
		return Maybe{Expression: Quote(m)}
	}
	return cond(
		list(
			eval(caar(c), a),
			NewLazyM(func() Maybe {
				return eval(cadar(c), a)
			}),
		),
		NewLazyM(func() Maybe {
			return list(quote(NewString("t")), EvconM(cdr(c), a))
		}),
	)
}

func Evcon(c, a *Expression) (*Expression, error) {
	if c == nil || a == nil {
		return nil, fmt.Errorf("nil arguments")
	}
	if !(c.IsList() && a.IsList()) {
		return nil, fmt.Errorf("needs lists")
	}
	for _, arg := range *c.List {
		car, err := Car(arg)
		if err != nil {
			return nil, err
		}
		r, err := Eval(car, a)
		if err != nil {
			return nil, err
		}
		if r.String() == "t" {
			cdr, err := Cdr(arg)
			if err != nil {
				return nil, err
			}
			cadr, err := Car(cdr)
			if err != nil {
				return nil, err
			}
			return Eval(cadr, a)
		}
	}
	return nil, fmt.Errorf("no condition satisfied")
}

func Cond(args ...*Expression) (*Expression, error) {
	fmt.Printf("cond(%q)\n", args)
	for _, a := range args {
		p, err := Car(a)
		if err != nil {
			return nil, err
		}
		if err := p.EvalLazy(); err != nil {
			return nil, err
		}
		e, err := Eval(p, NewList())
		if err != nil {
			return nil, err
		}
		if e.String() == "t" {
			cdr, err := Cdr(a)
			if err != nil {
				return nil, err
			}
			if err := cdr.EvalLazy(); err != nil {
				return nil, err
			}
			return cdr, nil
		}
	}
	return nil, fmt.Errorf("no condition satisfied")
}

func Read(s string) (*Expression, error) {
	n, err := parse(s)
	if err != nil {
		return nil, err
	}
	return n.Expression()
}

func Cadr(e *Expression) (*Expression, error) {
	a, err := Cdr(e)
	if err != nil {
		return nil, err
	}
	return Car(a)
}

func Caddr(e *Expression) (*Expression, error) {
	a, err := Cdr(e)
	if err != nil {
		return nil, err
	}
	b, err := Cdr(a)
	if err != nil {
		return nil, err
	}
	return Car(b)
}

func Caar(e *Expression) (*Expression, error) {
	a, err := Car(e)
	if err != nil {
		return nil, err
	}
	return Car(a)
}

func Cadar(e *Expression) (*Expression, error) {
	a, err := Car(e)
	if err != nil {
		return nil, err
	}
	b, err := Cdr(a)
	if err != nil {
		return nil, err
	}
	return Car(b)
}

func Car(e *Expression) (*Expression, error) {
	if e == nil {
		return Nil(), nil
	}
	if err := e.EvalLazy(); err != nil {
		return nil, err
	}
	if !e.IsList() {
		return nil, fmt.Errorf("can only car a list: %q", e)
	}
	if e.List.Empty() {
		return Nil(), nil
	}
	return e.List.First(), nil
}

func Cdr(e *Expression) (*Expression, error) {
	if e == nil {
		return Nil(), nil
	}
	if err := e.EvalLazy(); err != nil {
		return nil, err
	}
	if !e.IsList() {
		return nil, fmt.Errorf("can only cdr a list: %q", e)
	}
	if e.List.Empty() {
		return Nil(), nil
	}
	return NewList(e.List.Rest()...), nil
}

// either an atom or a list
type Expression struct {
	*Atom
	*List
	Lazy func() error // for lazy evaluation
}

func (e Expression) EvalLazy() error {
	if e.Lazy == nil {
		return nil
	}
	return e.Lazy()
}

func NewLazyM(f func() Maybe) Maybe {
	return Maybe{Expression: NewLazy(func() (*Expression, error) {
		out := f()
		return out.Expression, out.Error
	})}
}

func NewLazy(f func() (*Expression, error)) *Expression {
	var out Expression
	out.Lazy = func() error {
		e, err := f()
		if err != nil {
			return err
		}
		out = *e
		return nil
	}
	return &out
}

func NewString(s string) *Expression {
	return &Expression{
		Atom: &Atom{
			Type: "string",
			Blob: []byte(s),
		},
	}
}

func NewBlob(s []byte) *Expression {
	return &Expression{
		Atom: &Atom{
			Type: "blob",
			Blob: s,
		},
	}
}

func Nil() *Expression {
	return NewList()
}

func True() *Expression {
	return NewString("t")
}

func False() *Expression {
	return Nil()
}

func NewList(list ...*Expression) *Expression {
	var z List
	e := Expression{
		List: &z,
	}
	for _, x := range list {
		*e.List = append(*e.List, x)
	}
	return &e
}

func (e Expression) String() string {
	switch {
	case e.Lazy != nil:
		return "[LAZY]"
	case e.IsAtom():
		if e.Atom == nil {
			return "()"
		}
		return e.Atom.String()
	case e.IsList():
		var list []string
		for _, x := range *e.List {
			list = append(list, x.String())
		}
		return fmt.Sprintf("(%s)", strings.Join(list, " "))
	default:
		panic("illegal expression")
	}
}

func (e Expression) Boolean() bool {
	if e.IsList() {
		return false
	}
	return e.Atom.String() == "t"
}

func AtomF(e *Expression) (*Expression, error) {
	if e.IsAtom() {
		return NewString("t"), nil
	}
	return NewList(), nil
}

func (e Expression) IsAtom() bool {
	if e.Atom != nil {
		return true
	}
	if e.List == nil || len(*e.List) == 0 {
		return true
	}
	return false
}

func (e Expression) IsList() bool {
	return e.List != nil
}

func (e Expression) IsEmpty() bool {
	return e.List.Empty()
}

func (e Expression) Check() error {
	switch {
	case e.Atom == nil && e.List == nil && e.Lazy == nil:
		return fmt.Errorf("empty expression")
	case e.Atom != nil && e.List != nil:
		return fmt.Errorf("paradoxical expression")
	}
	return nil
}

type Atom struct {
	Type string
	Blob []byte
}

// TODO: resolve overall issues of null Atoms
func AtomsEqual(a, b *Atom) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil:
		return false
	case b == nil:
		return false
	}

	if a.Type != b.Type {
		return false
	}
	if !bytes.Equal(a.Blob, b.Blob) {
		return false
	}
	return true
}

func (a Atom) String() string {
	switch a.Type {
	case "string":
		return string(a.Blob)
	case "blob":
		return fmt.Sprintf("0x%x", a.Blob)
	default:
		panic("illegal type")
	}
}

type List []*Expression

func (l List) Empty() bool {
	return len(l) == 0
}

func (l List) First() *Expression {
	if len(l) == 0 {
		return nil
	}
	return l[0]
}

func (l List) Rest() List {
	if len(l) == 0 {
		return nil
	}
	return l[1:]
}

func (l *List) Add(e *Expression) {
	*l = append(*l, e)
}
