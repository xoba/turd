// expressions
package exp

import (
	"fmt"
	"math/big"
	"strings"
)

// if atom is nil, expression is an atom, otherwise a list
type Expression interface {
	Atom() Atom
	List() []Expression
	Error() error
	fmt.Stringer
}

type Atom interface {
	fmt.Stringer
}

type expr struct {
	atom Atom
	list list
	err  error
	lazy func() Expression
}

func (e expr) Atom() Atom {
	e.eval()
	return e.atom
}

func (e expr) List() []Expression {
	e.eval()
	return e.list
}

func (e expr) Error() error {
	e.eval()
	return e.err
}

func (e expr) String() string {
	e.eval()
	if a := e.atom; a != nil {
		return a.String()
	}
	if err := e.err; err != nil {
		return fmt.Sprintf("ERROR: %q", err.Error())
	}
	return e.list.String()
}

func (e *expr) eval() {
	if e.lazy == nil {
		return
	}
	defer func() {
		e.lazy = nil
	}()
	o := e.lazy()
	e.atom = o.Atom()
	e.list = o.List()
	e.err = o.Error()
}

type atom struct {
	v interface{}
}

func (a atom) String() string {
	switch t := a.v.(type) {
	case fmt.Stringer:
		return t.String()
	case string:
		return t
	case []byte:
		return fmt.Sprintf("0x%x", t)
	default:
		panic(fmt.Errorf("illegal atom type %T", t))
	}
}

type list []Expression

func (l list) String() string {
	var list []string
	for _, e := range l {
		list = append(list, e.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(list, " "))
}

func NewError(e error) Expression {
	return expr{err: e}
}

func Errorf(format string, a ...interface{}) Expression {
	return NewError(fmt.Errorf(format, a...))
}

func NewAtom(a fmt.Stringer) Expression {
	return expr{atom: a}
}

func NewString(x string) Expression {
	return NewAtom(atom{
		v: x,
	})
}

func NewBlob(x []byte) Expression {
	return NewAtom(atom{
		v: x,
	})
}

func NewInt(x int64) Expression {
	return NewAtom(atom{
		v: big.NewInt(x),
	})
}

func NewLazy(f func() Expression) Expression {
	return expr{lazy: f}
}

func NewList(list ...Expression) Expression {
	return expr{list: list}
}
