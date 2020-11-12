package lisp

import (
	"github.com/xoba/turd/lisp/exp"
)

type MonadFunc func(exp.Expression) exp.Expression

func Compose(funcs ...MonadFunc) MonadFunc {
	if len(funcs) == 0 {
		return funcs[0]
	}
	return func(a exp.Expression) exp.Expression {
		n := len(funcs)
		for i := 0; i < n; i++ {
			f := funcs[n-i-1]
			a = f(a)
		}
		return a
	}
}
