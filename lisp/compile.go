package lisp

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xoba/turd/lisp/exp"
)

const (
	pkg = "lisp/gen"
)

// TODO: also be able to invert this mapping... or simply reject reserved identifiers
func SanitizeGo(e exp.Expression) exp.Expression {
	// from the go spec
	var list []string
	add := func(category, words string) {
		list = append(list, strings.Fields(words)...)
	}

	add("keywords", `break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
`)
	add("functions", `      append cap close complex copy delete imag len
       make new panic print println real recover
`)
	add("constants", `      true false iota
`)
	add("zero", "nil")
	add("types", `  bool byte complex64 complex128 error float32 float64
       int int8 int16 int32 int64 rune string
       uint uint8 uint16 uint32 uint64 uintptr
`)
	for _, x := range list {
		e = translateAtoms(x, "go_sanitized_"+x, e)
	}
	return e
}

func ToExpression(e exp.Expression) ([]byte, error) {
	w := new(bytes.Buffer)
	switch {
	case e.Atom() != nil:
		fmt.Fprintf(w, "quote(%q)", e)
	default:
		list := e.List()
		var parts []string
		for _, x := range list {
			buf, err := ToExpression(x)
			if err != nil {
				return nil, err
			}
			parts = append(parts, string(buf))
		}
		fmt.Fprintf(w, "list(%s)", strings.Join(parts, ","))
	}
	return w.Bytes(), nil
}

func in(a string, list ...string) bool {
	for _, x := range list {
		if a == x {
			return true
		}
	}
	return false
}

func translateAtoms(from, to string, e exp.Expression) exp.Expression {
	a := e.Atom()
	if a == nil {
		var out []exp.Expression
		for _, c := range e.List() {
			out = append(out, translateAtoms(from, to, c))
		}
		return exp.NewList(out...)
	}
	if a.String() == from {
		return exp.NewString(to)
	}
	return e
}
