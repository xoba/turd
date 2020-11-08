package scr

import (
	"fmt"
	"math/rand"

	"github.com/xoba/turd/cnfg"
)

type EvalFunc func(args ...*Expression) (*Expression, error)

type MonadFunc func(args ...Maybe) Maybe

type Maybe struct {
	*Expression
	Error error
}

func (m Maybe) String() string {
	if m.Error != nil {
		return fmt.Sprintf("error: %v", m.Error)
	}
	if m.Expression == nil {
		return "nil"
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
			list = append(list, m.Expression)
		}
		out, err := f(list...)
		if err != nil {
			return Maybe{Error: err}
		}
		return Maybe{Expression: out}
	}
}

func Compose(funcs ...MonadFunc) MonadFunc {
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

func TestMonad(cnfg.Config) error {

	const n = 100

	randErr := func() error {
		if rand.Intn(n) == 0 {
			return fmt.Errorf("fake")
		}
		return nil
	}

	Car := func(args ...*Expression) (*Expression, error) {
		if err := randErr(); err != nil {
			return nil, err
		}
		v, err := Car(args[0])
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	Cdr := func(args ...*Expression) (*Expression, error) {
		if err := randErr(); err != nil {
			return nil, err
		}
		v, err := Cdr(args[0])
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	car := EvalFunc(Car).ToMonad()
	cdr := EvalFunc(Cdr).ToMonad()
	caddr := Compose(car, cdr, cdr)
	cadadr := Compose(car, cdr, car, cdr)

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

	return nil
}
