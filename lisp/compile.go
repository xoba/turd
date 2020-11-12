package lisp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/lisp/exp"
)

func CompileDef(cnfg.Config) error {

	const (
		pkg = "lisp/gen"
	)

	if err := os.MkdirAll(pkg, os.ModePerm); err != nil {
		return err
	}

	file := filepath.Join(pkg, "gen.go")
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, `package %s

import (
	"github.com/xoba/turd/lisp/exp"
)

`, path.Base(pkg))

	for _, x := range strings.Split("list,newlist,atom,cond,cons,eq,quote,Nil,True", ",") {
		fmt.Fprintf(f, `

func %s(...exp.Expression) exp.Expression{
panic("")
}

`, x)
	}

	fmt.Fprintf(f, `

type Func func(...exp.Expression) exp.Expression 

func apply(Func, ...exp.Expression) exp.Expression {
panic("")
}

var t = True()

`)

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
		fmt.Println(defs)
	}
	for _, def := range defs { // strings.Split("caddr,caar,null", ",") {
		buf, err := ioutil.ReadFile(def)
		if err != nil {
			return err
		}
		e, err := Read(string(buf))
		if err != nil {
			return err
		}
		code, err := Tofunc(e)
		if err != nil {
			return err
		}
		fmt.Fprintf(f, string(code))
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := Gofmt(file); err != nil {
		return err
	}
	return nil
}

func Tofunc(defun exp.Expression) ([]byte, error) {
	if Car(defun).String() != "defun" {
		return nil, fmt.Errorf("not a defun")
	}
	name := Car(Cdr(defun))
	args := Car(Cdr(Cdr(defun)))
	body := Car(Cdr(Cdr(Cdr(defun))))
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "func %s(args ... exp.Expression) exp.Expression {\n", name)
	for i, a := range args.List() {
		if !IsAtom(a) {
			return nil, fmt.Errorf("not atom: %s", a)
		}
		fmt.Fprintf(w, "%s := args[%d];\n", a.Atom(), i)
	}
	code, err := Compile(body, false)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(w, "return %s\n}\n\n", string(code))
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

func Compile(e exp.Expression, noapply bool) ([]byte, error) {
	w := new(bytes.Buffer)
	switch {
	case e.Atom() != nil:
		var atom string
		switch a := e.Atom().String(); a {
		case "quote", "atom", "eq", "car", "cdr", "cons", "cond":
			atom = a
		default:
			atom = a
		}
		fmt.Fprint(w, atom)
	default:
		n := len(e.List())
		switch {
		case n == 0:
			fmt.Fprintf(w, "Nil()")
		case n == 1:
			sub, err := Compile(e.List()[0], false)
			if err != nil {
				return nil, err
			}
			fmt.Fprintf(w, "(%s)", string(sub))
		case n > 1 && in(e.List()[0].String(), "cond", "quote"):
			var list []string
			for _, a := range e.List() {
				sub, err := Compile(a, true)
				if err != nil {
					return nil, err
				}
				list = append(list, string(sub))
			}
			fmt.Fprintf(w, "apply(/* cond */ %s)", strings.Join(list, ","))
		default:
			var list []string
			for _, a := range e.List() {
				sub, err := Compile(a, false)
				if err != nil {
					return nil, err
				}
				list = append(list, string(sub))
			}
			if noapply {
				fmt.Fprintf(w, "newlist")
			} else {
				fmt.Fprintf(w, "apply")
			}
			fmt.Fprintf(w, "(%s)", strings.Join(list, ","))
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
