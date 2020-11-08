package scr

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/xoba/turd/cnfg"
)

func Lisp(cnfg.Config) error {
	test := func(in, expect string) {
		check := func(e error) {
			if e != nil {
				log.Fatal(e)
			}
		}
		e, err := Read(in)
		check(err)
		a := NewList()
		x, err := Eval(e, a)
		check(err)
		fmt.Printf("(eval %s %s) = %s\n", e, a, x)
		if got := x.String(); got != expect {
			check(fmt.Errorf("expected %q, got %q", expect, got))
		}
	}
	test("(quote z)", "z")
	test("(atom 'a)", "t")
	test("(atom '())", "()")
	test("(atom '(1 2 3))", "()")
	test("(eq 'a 'a)", "t")
	test("(eq 'a 'b)", "()")
	test("(eq 'b 'a)", "()")
	test("(eq 'b '())", "()")
	test("(eq '() '())", "t")
	return nil
}

func Assoc(x, y *Expression) (*Expression, error) {
	fmt.Printf("assoc(%s, %s)\n", x, y)
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
	if e.IsAtom() {
		return Assoc(e, a)
	}
	car, err := Car(e)
	if err != nil {
		return nil, err
	}
	if car.IsAtom() {
		switch car.String() {
		case "quote":
			return Cadr(e)
		case "atom":
			cadr, err := Cadr(e)
			if err != nil {
				return nil, err
			}
			x, err := Eval(cadr, a)
			if err != nil {
				return nil, err
			}
			if x.IsAtom() {
				return True(), nil
			}
			return False(), nil
		case "eq":
			cadr, err := Cadr(e)
			if err != nil {
				return nil, err
			}
			caddr, err := Caddr(e)
			if err != nil {
				return nil, err
			}
			first, err := Eval(cadr, a)
			if err != nil {
				return nil, err
			}
			second, err := Eval(caddr, a)
			if err != nil {
				return nil, err
			}
			return Eq(first, second)
		case "car":
		case "cdr":
		case "cons":
		case "cond":
		default:
		}
	}
	return nil, fmt.Errorf("eval unimplemented")
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
	if e.IsAtom() {
		return nil, fmt.Errorf("can't car an atom")
	}
	if e.List.Empty() {
		return Nil(), nil
	}
	return e.List.First(), nil
}

func Cdr(e *Expression) (*Expression, error) {
	if e.IsAtom() {
		return nil, fmt.Errorf("can't cdr an atom")
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
}

func NewQuote(e *Expression) *Expression {
	out := NewList()
	out.Add(NewString("quote"))
	out.Add(e)
	return out
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
	case e.IsAtom():
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

func (e Expression) IsAtom() bool {
	e.Check()
	return e.Atom != nil
}

func (e Expression) IsList() bool {
	return !e.IsAtom()
}

func (e Expression) IsEmpty() bool {
	e.Check()
	return e.List.Empty()
}

func (e Expression) Check() {
	switch {
	case e.Atom == nil && e.List == nil:
		panic("empty expression")
	case e.Atom != nil && e.List != nil:
		panic("paradoxical expression")
	}
}

type Atom struct {
	Type string
	Blob []byte
}

func AtomsEqual(a, b *Atom) bool {
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
