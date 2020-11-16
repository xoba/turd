package lisp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xoba/turd/cnfg"
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
	add("functions", `	append cap close complex copy delete imag len
	make new panic print println real recover
`)
	add("constants", `	true false iota
`)
	add("zero", "nil")
	add("types", `	bool byte complex64 complex128 error float32 float64
	int int8 int16 int32 int64 rune string
	uint uint8 uint16 uint32 uint64 uintptr
`)
	for _, x := range list {
		e = translateAtoms(x, "go_sanitized_"+x, e)
	}
	return e
}

func CompileDef(cnfg.Config) error {
	if err := os.MkdirAll(pkg, os.ModePerm); err != nil {
		return err
	}
	file := filepath.Join(pkg, "gen.go")
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "// THIS FILE IS AUTOGENERATED, DO NOT EDIT!\n\n")
	buf, err := ioutil.ReadFile("lisp/gensource/source.go")
	if err != nil {
		return err
	}
	if _, err := f.Write(buf); err != nil {
		return err
	}
	fmt.Fprint(f, "\n\n\n")
	var defs []string
	{
		const dir = "defs"
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, fi := range files {
			if name := fi.Name(); filepath.Ext(name) == ".lisp" {
				defs = append(defs, filepath.Join(dir, name))
			}
		}
	}
	//fmt.Printf("%d defs = %q\n", len(defs), defs)
	var names []string
	for _, def := range defs {
		buf, err := ioutil.ReadFile(def)
		if err != nil {
			return err
		}
		e, err := Read(string(buf))
		if err != nil {
			return err
		}
		//fmt.Printf("%s -> %s\n", def, e)
		e = SanitizeGo(e)
		{
			name, env, err := ToEnv(e)
			if err != nil {
				return err
			}
			fmt.Fprintf(f, "var env_%s = %s\n", name, string(env))
		}

		name, code, err := Tofunc(e)
		if err != nil {
			return err
		}
		fmt.Fprint(f, string(code))

		names = append(names, name)

	}

	fmt.Fprintf(f, "\n\nfunc init() { env = list(")
	for _, n := range names {
		fmt.Fprintf(f, "env_%s,", n)
	}
	fmt.Fprintf(f, ")}\n\n")

	if err := f.Close(); err != nil {
		return err
	}
	if err := Gofmt(file); err != nil {
		return err
	}
	return nil
}

// TODO: instead, express this as a string (not as code) that can be parsed like lisp.Read()
func ToEnv(defun exp.Expression) (string, []byte, error) {
	return ToEnv0(defun)
}

func ToEnv0(defun exp.Expression) (string, []byte, error) {
	if Car(defun).String() != "defun" {
		return "", nil, fmt.Errorf("not a defun")
	}
	name := Car(Cdr(defun))
	args := Car(Cdr(Cdr(defun)))
	body := Car(Cdr(Cdr(Cdr(defun))))
	q := func(s string) exp.Expression {
		return exp.NewString(s)
	}
	nl := func(args ...exp.Expression) exp.Expression {
		return exp.NewList(args...)
	}
	e := nl(
		name,
		nl(
			q("label"),
			name,
			nl(
				q("lambda"),
				args,
				body,
			),
		),
	)
	exp, err := ToExpression(e)
	return name.String(), exp, err
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

func Tofunc(defun exp.Expression) (string, []byte, error) {
	if Car(defun).String() != "defun" {
		return "", nil, fmt.Errorf("not a defun")
	}
	name := Car(Cdr(defun)).String()
	args := Car(Cdr(Cdr(defun)))
	body := Car(Cdr(Cdr(Cdr(defun))))
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "func %[1]s(args ... Exp) Exp {\n", name)
	//fmt.Fprintf(w, "debug(%q,args...)\n", name)
	fmt.Fprintf(w, "checklen(%d,args)\n", len(args.List()))
	for i, a := range args.List() {
		if !IsAtom(a) {
			return name, nil, fmt.Errorf("not atom: %s", a)
		}
		fmt.Fprintf(w, "%s := args[%d];\n", a.Atom(), i)
	}
	code, err := Compile(body, true)
	if err != nil {
		return name, nil, err
	}
	fmt.Fprintf(w, "return %s\n}\n\n", string(code))
	return name, w.Bytes(), nil
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

func CompileLazy(e exp.Expression) ([]byte, error) {
	w := new(bytes.Buffer)
	list := e.List()
	if len(list) != 2 {
		return nil, fmt.Errorf("malformed cond: %s %s", e)
	}
	pb, err := Compile(list[0], false)
	if err != nil {
		return nil, err
	}
	eb, err := Compile(list[1], false)
	if err != nil {
		return nil, err
	}
	f := func(s string) string {
		return fmt.Sprintf(`Func(func(...Exp) Exp {
return %s
})`, s)
	}
	fmt.Fprintf(w, "list(\n%s,\n%s,\n)", f(string(pb)), f(string(eb)))
	return w.Bytes(), nil
}

func Compile(e exp.Expression, indent bool) ([]byte, error) {
	//	fmt.Printf("%d compile(%q)\n", len(e.List()), e)
	w := new(bytes.Buffer)
	emit := func(msg string, list []string) {
		if indent {
			fmt.Fprintf(w, "%s(\n%s,\n)", msg, strings.Join(list, ",\n"))
		} else {
			fmt.Fprintf(w, "%s(%s)", msg, strings.Join(list, ","))
		}
	}
	switch {
	case e.Atom() != nil:
		var x string
		x = e.Atom().String()
		fmt.Fprint(w, x)
	default:
		n := len(e.List())
		switch {
		case n == 0:
			fmt.Fprintf(w, "Nil")
		case n == 1:
			return nil, fmt.Errorf("illegal list of length 1")
		case n == 2 && in(e.List()[0].String(), "quote"):
			x := e.List()[1]
			switch {
			case x.Atom() != nil:
				fmt.Fprintf(w, "%q", x)
			default:

				compiled, err := Compile(e.List()[1], false)
				if err != nil {
					return nil, err
				}
				fmt.Fprintf(w, string(compiled))
			}
		case n > 1 && in(e.List()[0].String(), "cond"):
			var list []string
			for i, a := range e.List() {
				var f func(exp.Expression) ([]byte, error)
				if i == 0 {
					f = func(e exp.Expression) ([]byte, error) {
						return Compile(e, true)
					}
				} else {
					f = CompileLazy
				}
				sub, err := f(a)
				if err != nil {
					return nil, err
				}
				list = append(list, string(sub))
			}
			emit("apply", list)
		default:
			var list []string
			for _, a := range e.List() {
				sub, err := Compile(a, indent)
				if err != nil {
					return nil, err
				}
				list = append(list, string(sub))
			}
			emit("apply", list)
		}
	}
	return w.Bytes(), nil
}

func Gofmt(file string) error {
	w := new(bytes.Buffer)
	cmd := exec.Command("gofmt", "-w", file)
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %q", err, w)
	}
	return nil
}
