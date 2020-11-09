package scr

import (
	"fmt"
	"math/rand"

	"github.com/xoba/turd/cnfg"
)

type EvalFunc func(args ...*Expression) (*Expression, error)
type Eval1Func func(a *Expression) (*Expression, error)
type Eval2Func func(a, b *Expression) (*Expression, error)

type MonadFunc func(args ...Maybe) Maybe

type Maybe struct {
	*Expression
	Error error
}

func (m MonadFunc) ToEvalFunc() EvalFunc {
	return func(args ...*Expression) (*Expression, error) {
		var list []Maybe
		for _, a := range args {
			list = append(list, Maybe{Expression: a})
		}
		out := m(list...)
		return out.Expression, out.Error
	}
}

func (m Maybe) String() string {
	if m.Error != nil {
		return fmt.Sprintf("error: %v", m.Error)
	}
	return m.Expression.String()
}

func (f EvalFunc) ToMonad() MonadFunc {
	return func(args ...Maybe) Maybe {
		var list []*Expression
		for _, m := range args {
			if m.Error != nil {
				return m
			}
			if m.Expression == nil {
				return Maybe{Error: fmt.Errorf("nil expression")}
			}
			if err := m.Check(); err != nil {
				return Maybe{Error: err}
			}
			list = append(list, m.Expression)
		}
		out, err := f(list...)
		if err != nil {
			return Maybe{Error: err}
		}
		return Maybe{Expression: out}
	}
}

func (f Eval1Func) ToMonad() MonadFunc {
	return func(args ...Maybe) Maybe {
		f2 := EvalFunc(func(args ...*Expression) (*Expression, error) {
			if n := len(args); n != 1 {
				return nil, fmt.Errorf("needs 1 argument, got %d", n)
			}
			return f(args[0])
		})
		m := f2.ToMonad()
		return m(args...)
	}
}

func (f Eval2Func) ToMonad() MonadFunc {
	return func(args ...Maybe) Maybe {
		f2 := EvalFunc(func(args ...*Expression) (*Expression, error) {
			if n := len(args); n != 2 {
				return nil, fmt.Errorf("needs 2 argument, got %d", n)
			}
			return f(args[0], args[1])
		})
		m := f2.ToMonad()
		return m(args...)
	}
}

func Compose(funcs ...MonadFunc) MonadFunc {
	if len(funcs) == 0 {
		return funcs[0]
	}
	return func(args ...Maybe) Maybe {
		if len(args) != 1 {
			return Maybe{Error: fmt.Errorf("can only compose functions of one argument")}
		}
		a := args[0]
		n := len(funcs)
		for i := 0; i < n; i++ {
			f := funcs[n-i-1]
			a = f(a)
		}
		return a
	}
}

func Wrap(args ...*Expression) (out []Maybe) {
	for _, a := range args {
		out = append(out, Maybe{Expression: a})
	}
	return
}

func TestMonad(cnfg.Config) error {

	const n = 5

	randErr := func(args ...Maybe) Maybe {
		if rand.Intn(n) == 0 {
			return Maybe{Error: fmt.Errorf("fake")}
		}
		return args[0]
	}

	car := Compose(randErr, Eval1Func(Car).ToMonad())
	cdr := Compose(randErr, Eval1Func(Cdr).ToMonad())

	caddr := Compose(car, cdr, cdr)
	cadadr := Compose(car, cdr, car, cdr)
	cdddr := Compose(cdr, cdr, cdr)

	list := Maybe{
		Expression: NewList(
			NewString("a"),
			NewList(
				NewString("x"),
				NewString("y"),
				NewString("z"),
			),
			NewString("b"),
			NewString("d"),
		),
	}

	fmt.Printf("list = %s\n", list)

	fmt.Printf("car = %s\n", car(list))
	fmt.Printf("cdr = %s\n", cdr(list))

	fmt.Printf("caddr = %s\n", car(cdr(cdr(list))))
	fmt.Printf("caddr = %s\n", caddr(list))

	fmt.Printf("cadadr = %s\n", car(cdr(car(cdr(list)))))
	fmt.Printf("cadadr = %s\n", cadadr(list))

	fmt.Printf("cdddr = %s\n", cdr(cdr(cdr(list))))
	fmt.Printf("cdddr = %s\n", cdddr(list))

	return nil
}
