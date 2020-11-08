// expressions
package exp

// atom must be nil if length of list is nonzero
type Expression interface {
	Atom() (*Atom, error)
	List() ([]Expression, error)
}

type Atom struct {
	Type string
	Blob []byte
}

type expr struct {
	atom *Atom
	list []Expression
	lazy func() (Expression, error)
}

func (e *expr) eval() error {
	if e.lazy == nil {
		return nil
	}
	o, err := e.lazy()
	e.lazy = nil
	if err != nil {
		return err
	}
	atom, err := o.Atom()
	if err != nil {
		return err
	}
	e.atom = atom
	list, err := o.List()
	if err != nil {
		return err
	}
	e.list = list
	return nil
}

func (e expr) Atom() (*Atom, error) {
	if err := e.eval(); err != nil {
		return nil, err
	}
	return e.atom, nil
}

func (e expr) List() ([]Expression, error) {
	if err := e.eval(); err != nil {
		return nil, err
	}
	return e.list, nil
}

func NewAtom(a *Atom) Expression {
	return expr{atom: a}
}

func NewLazy(f func() (Expression, error)) Expression {
	return expr{lazy: f}
}

func NewList(list []Expression) Expression {
	return expr{list: list}
}
