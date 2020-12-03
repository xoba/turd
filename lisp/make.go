package lisp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/thash"
)

const (
	optimize = true
	pkg      = "lisp"
)

// autogenerate an eval func, TODO also "try"
func EvalTemplate(cnfg.Config) error {
	if err := GenEval("defs/compiled/eval.lisp", map[string]string{
		"nextlambda_prefix": "",
		"nextlambda_suffix": "",
		"defun":             "eval",
		"args":              "(e a)",
		"eval":              "eval",
		"assoc":             "assoc",
		"evlis":             "evlis",
		"evcon":             "evcon",
		"comment":           autogen,
	}); err != nil {
		return err
	}
	if err := GenEval("defs/compiled/teval.lisp", map[string]string{
		"nextlambda_prefix": "((lambda (t2)",
		"nextlambda_suffix": "t2))",
		"defun":             "teval",
		"args":              "(t e a)",
		"eval":              "teval t2",
		"assoc":             "tassoc t2",
		"evlis":             "tevlis t2",
		"evcon":             "tevcon t2",
		"comment":           autogen,
	}); err != nil {
		return err
	}
	return nil
}

func GenEval(file string, args map[string]string) error {
	addkv := func(k, v string, m map[string]string) map[string]string {
		out := map[string]string{
			k: v,
		}
		for k, v := range m {
			out[k] = v
		}
		return out
	}
	type compiled struct {
		name  string
		class string
		args  int
	}
	m := make(map[string]compiled)
	counts := make(map[string]int)
	set := func(c compiled) {
		if c.name == "" {
			return
		}
		counts[c.class]++
		if _, ok := m[c.name]; ok {
			log.Fatalf("duplicate %q", c.name)
		}
		m[c.name] = c
	}
	load := func(name string, args int) {
		set(compiled{name: name, class: "loaded", args: args})
	}
	add := func(name string, args int) {
		set(compiled{name: name, class: "manual", args: args})
	}
	axiom := func(name string, args int) {
		set(compiled{name: name, class: "axiom", args: args})
	}
	axiom("atom", 1)
	axiom("eq", 2)
	axiom("car", 1)
	axiom("cdr", 1)
	axiom("cons", 2)

	add("display", 1)
	add("exp", 3)
	add("mul", 2)
	add("add", 2)
	add("sub", 2)
	add("hash", 1)
	add("hashed", 1)
	add("concat", 2)
	add("newkey", 0)
	add("pub", 1)
	add("sign", 2)
	add("verify", 3)
	add("after", 2)
	add("err", 1)
	add("runes", 1)

	const dirname = "defs/compiled"
	dir, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	c := NewContext()
	for _, fi := range dir {
		file := filepath.Join(dirname, fi.Name())
		if filepath.Ext(file) != ".lisp" {
			continue
		}
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		e, err := Parse(string(buf))
		if err != nil {
			return err
		}
		name, args, _, err := DefunCode(c, e)
		if err != nil {
			return err
		}
		load(name, len(args))
	}
	var sorted []string
	for k := range m {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	w := new(bytes.Buffer)
	for _, k := range sorted {
		c := m[k]
		fmt.Fprintf(w, ";; %q with %d args (%s)\n", c.name, c.args, c.class)
		emit := func(s string) {
			t := template.Must(template.New("expr").Parse(s + "\n\n"))
			if err := t.Execute(w, addkv("name", c.name, args)); err != nil {
				log.Fatal(err)
			}
		}
		switch c.args {
		case 0:
			emit(`((eq op '{{.name}}) ({{.name}}))`)
		case 1:
			emit(`((eq op '{{.name}}) ({{.name}} ({{.eval}} first a)))`)
		case 2:
			emit(`((eq op '{{.name}}) ({{.name}} ({{.eval}} first a) ({{.eval}} second a)))`)
		case 3:
			emit(`((eq op '{{.name}}) ({{.name}} ({{.eval}} first a) ({{.eval}} second a) ({{.eval}} third a)))`)
		default:
			return fmt.Errorf("illegal args: %d", c.args)
		}
	}

	t, err := template.ParseFiles("defs/template/eval.lisp")
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := t.Execute(f, addkv("compiled", "\n"+w.String(), args)); err != nil {
		return err
	}
	fmt.Printf("%s compiled: %v\n", filepath.Base(file), counts)
	return nil
}

const autogen = "THIS FILE IS AUTOGENERATED, DO NOT EDIT!"

func CompileDefuns(cnfg.Config) error {
	if err := os.MkdirAll(pkg, os.ModePerm); err != nil {
		return err
	}
	const filename = "gen.go"
	file := filepath.Join(pkg, filename)
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, `// %[1]s

package %[2]s

import "fmt"

func init() {
return
fmt.Println("%[3]s: %[4]s");
}
`, autogen, path.Base(pkg), filename, autogen)
	fmt.Fprintln(f, `var( L = list
A = apply
)
`)
	fmt.Fprintf(f, `func parse_env(s string) Exp {
e,err:= Parse(s)
if err != nil {
panic(err)
}
return e
}

`)

	// definitions referenced by eval need
	// to be compiled, since eval is compiled
	type definition struct {
		file     string
		name     string
		expr     Exp
		compiled bool
	}

	var defs []*definition

	load := func(dir string, compiled bool) error {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, fi := range files {
			if name := fi.Name(); filepath.Ext(name) == ".lisp" {
				def := &definition{
					file:     filepath.Join(dir, name),
					compiled: compiled,
				}
				buf, err := ioutil.ReadFile(def.file)
				if err != nil {
					return err
				}
				e, err := Parse(string(buf))
				if err != nil {
					return err
				}
				def.name = String(cadr(e))
				def.expr = e
				defs = append(defs, def)
			}
		}
		return nil
	}

	load("defs/compiled", true)
	load("defs/interpreted", false)

	sort.Slice(defs, func(i, j int) bool {
		return defs[i].name < defs[j].name
	})

	c := NewContext()

	for _, def := range defs {
		e := def.expr
		e = SanitizeGo(e)
		name := String(cadr(e))
		label, err := LabelExpr(e)
		if err != nil {
			return err
		}
		var msg string
		if def.compiled {
			msg = "compiled"
		} else {
			msg = "interpreted"
		}
		fmt.Fprintf(f, "\n\n//\n// %s (%s)\n//\n\n\n", def.name, msg)
		fmt.Fprintf(f, "var %[1]s_label = parse_env(%[2]q)\n", name, String(label))
		name, _, code, err := DefunCode(c, e)
		if err != nil {
			return err
		}
		if def.compiled {
			fmt.Fprint(f, string(code))
		}
		def.name = name
	}

	if true {
		fmt.Fprintln(f, c.emit())
	}

	fmt.Fprintf(f, "\n\nfunc init() { env = L(\n")
	for _, def := range defs {
		fmt.Fprintf(f, "L(%q,%s_label),\n", def.name, def.name)
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

func DefunCode(c context, defun Exp) (string, []string /* args */, []byte, error) {
	if String(car(defun)) != "defun" {
		return "", nil, nil, fmt.Errorf("not a defun")
	}
	name := String(cadr(defun))
	var args []Exp
	if x, ok := caddr(defun).([]Exp); ok {
		args = x
	} else {
		args = []Exp{caddr(defun)}
	}
	var arglist []string
	for _, a := range args {
		arglist = append(arglist, String(a))
	}
	body := cadddr(defun)
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "func %[1]s(_args ... Exp) Exp {\n", name)
	var vars []string
	for i, a := range args {
		if !isString(a) {
			return name, nil, nil, fmt.Errorf("not a string: %s", a)
		}
		v := String(a)
		fmt.Fprintf(w, "%s := _args[%d];\n", v, i)
		vars = append(vars, v)
	}
	code, err := Compile(c, body, true, dedup(vars))
	if err != nil {
		return name, nil, nil, err
	}
	fmt.Fprintf(w, "return %s\n}\n\n", string(code))
	return name, arglist, w.Bytes(), nil
}

func isString(e Exp) bool {
	switch e.(type) {
	case string:
		return true
	default:
		return false
	}
}

type context struct {
	cases map[string]string
	funcs map[string]string
}

func NewContext() context {
	return context{
		cases: make(map[string]string),
		funcs: make(map[string]string),
	}
}

func (c context) emit() string {
	w := new(bytes.Buffer)

	f := func(name string, m map[string]string) {
		var list []string
		for k := range m {
			list = append(list, k)
		}
		sort.Strings(list)
		fmt.Fprintf(w, "\n// %s:\n\n", name)
		for _, k := range list {
			fmt.Fprintf(w, ` 
%s
`, m[k])
		}
	}
	f("cases", c.cases)
	//f("funcs", c.funcs)

	return w.String()
}

// non-negligible chance of collision, but maybe worth it for brevity:
func smallHash(s string) string {
	return hex.EncodeToString(thash.Hash([]byte(s)))[:4]
}

func funcName(name, code string, vars []string) string {
	clean := func(x string) string {
		w := new(bytes.Buffer)
		for _, r := range x {
			switch {
			case unicode.IsDigit(r), unicode.IsLetter(r):
				w.WriteRune(r)
			default:
				w.WriteRune('π')
			}
		}
		return w.String()
	}
	args := strings.Join(vars, ",")
	if name == "" {
		return fmt.Sprintf("F_%s_%s", smallHash(code), smallHash(args))
	}
	return fmt.Sprintf("F_%s_%s_%s", clean(name), smallHash(code), smallHash(args))
}

func CompileLazy(c context, e Exp, vars []string) ([]byte, error) {
	w := new(bytes.Buffer)
	list, ok := e.([]Exp)
	if !ok {
		return nil, fmt.Errorf("lazy not a list")
	}
	if len(list) != 2 {
		return nil, fmt.Errorf("malformed cond with %d parts: %s", len(list), e)
	}
	f := func(s string) string {
		name := funcName("", s, vars)
		c.funcs[name] = fmt.Sprintf(`func %s(...Exp) Exp {
return %s
}
`, name, s)
		return fmt.Sprintf(`Func(func(...Exp) Exp {
return %s
})`, s)
	}
	fc := func(e Exp) (string, error) {
		switch t := e.(type) {
		case string:
			return t, nil
		case []Exp:
			switch {
			case len(t) == 0:
				return "Nil", nil
			case len(t) == 2 && String(t[0]) == "quote":
				return compileQuote(c, t[1], vars)
			}
		}
		pb, err := Compile(c, e, false, vars)
		if err != nil {
			return "", err
		}
		return f(string(pb)), nil
	}
	pf, err := fc(list[0])
	if err != nil {
		return nil, err
	}
	ef, err := fc(list[1])
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(w, `L(
%s,
%s,
)`, pf, ef)
	return w.Bytes(), nil
}

func compileQuote(c context, x Exp, vars []string) (string, error) {
	switch t := x.(type) {
	case string:
		return fmt.Sprintf("%q", t), nil
	default:
		compiled, err := Compile(c, t, false, vars)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(string(compiled)), nil
	}
}

func Compile(c context, e Exp, indent bool, vars []string) ([]byte, error) {
	w := new(bytes.Buffer)
	emit := func(list []string) {
		if indent {
			fmt.Fprintf(w, "A(\n%s,\n)", strings.Join(list, ",\n"))
		} else {
			fmt.Fprintf(w, "A(%s)", strings.Join(list, ","))
		}
	}
	lambda := func(e Exp, name string) error {
		var args []string
		for _, e := range cadar(e).([]Exp) {
			args = append(args, String(e))
		}
		var arglist []string
		for _, x := range cdr(e).([]Exp) {
			arg, err := Compile(c, x, false, vars)
			if err != nil {
				return err
			}
			arglist = append(arglist, string(arg))
		}
		fmt.Fprintf(w, `func() Exp {
var %[1]s func(... Exp) Exp
%[1]s = func(_args ... Exp) Exp {
`,
			name,
		)
		for i, a := range args {
			fmt.Fprintf(w, "%s := _args[%d]\n", a, i)
			vars = append(vars, a)
		}
		vars = dedup(vars)
		body, err := Compile(c, caddar(e), false, vars)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, `return %[2]s
}
return %[1]s(%[3]s)
}()
`,
			name,
			string(body),
			strings.Join(arglist, ","),
		)
		return nil
	}

	switch e := e.(type) {
	case string:
		fmt.Fprint(w, e)
	case []Exp:
		n := len(e)
		switch {

		case n == 0:
			fmt.Fprintf(w, "Nil")

		case ExpToBool(eq(car(e), "quote")):
			q, err := compileQuote(c, e[1], vars)
			if err != nil {
				return nil, err
			}
			fmt.Fprint(w, q)

		case ExpToBool(eq(car(e), "cond")):

			if op, ok := isMaplikeCond(e); optimize && ok {

				var funcs, names []string

				for i, a := range e {
					if i == 0 {
						continue
					}
					if i == len(e)-1 {
						continue
					}
					eq := car(a)
					key := car(cdr(car(cdr(cdr(eq)))))
					sub, err := Compile(c, cadr(a), false, vars)
					if err != nil {
						return nil, err
					}
					name := funcName(String(key), string(sub), vars)
					w := new(bytes.Buffer)
					fmt.Fprintf(w, `func %[1]s(%[2]s Exp) Exp {
return %[3]s
}`,
						name,
						strings.Join(vars, ","),
						string(sub),
					)
					f := w.String()
					names = append(names, String(key))
					funcs = append(funcs, name)
					c.cases[name] = f
				}

				fmt.Fprintf(w, "func() Exp {\n")
				kv := new(bytes.Buffer)
				for i, n := range names {
					fmt.Fprintf(kv, "%q: %s,\n", n, funcs[i])
				}
				mapname := fmt.Sprintf("map_%s", smallHash(kv.String()))
				{
					g := new(bytes.Buffer)
					// TODO: move map "m" to global and just reference it here, not re-create it
					fmt.Fprintf(g, "var %s =make( map[string]func(%s Exp) Exp)\n", mapname, strings.Join(vars, ","))
					fmt.Fprintf(g, "func init() {\n")
					fmt.Fprintf(g, " %s = map[string]func(%s Exp) Exp {\n", mapname, strings.Join(vars, ","))
					fmt.Fprint(g, kv)
					fmt.Fprintf(g, "}}\n")
					c.cases[mapname] = g.String()
				}
				t, err := Compile(c, cadr(e[len(e)-1]), false, vars)
				if err != nil {
					return nil, err
				}
				fmt.Fprintf(w, `if f,ok := %[4]s[String(%[1]s)]; ok {
return f(%[2]s)
}
return %[3]s
`,
					op,
					strings.Join(vars, ","),
					string(t),
					mapname,
				)
				fmt.Fprintf(w, "}()\n")
			} else {

				var list []string
				for i, a := range e {
					var f func(context, Exp, []string) ([]byte, error)
					if i == 0 {
						f = func(c context, e Exp, vars []string) ([]byte, error) {
							return Compile(c, e, true, vars)
						}
					} else {
						f = CompileLazy
					}
					sub, err := f(c, a, vars)
					if err != nil {
						return nil, err
					}
					list = append(list, string(sub))
				}
				emit(list)
			}

		case ExpToBool(eq(caar(e), "lambda")) || ExpToBool(eq(caar(e), "λ")):
			if err := lambda(e, "λ"); err != nil {
				return nil, err
			}

		case ExpToBool(eq(caar(e), "label")):
			expr := cons(car(cdr(cdr(car(e)))), cdr(e))
			if err := lambda(expr, String(cadar(e))); err != nil {
				return nil, err
			}

		default:
			var list []string
			for _, a := range e {
				sub, err := Compile(c, a, indent, vars)
				if err != nil {
					return nil, err
				}
				list = append(list, string(sub))
			}
			emit(list)
		}

	default:
		return nil, fmt.Errorf("can't compile %T", e)
	}
	return w.Bytes(), nil
}

func dedup(vars []string) (out []string) {
	m := make(map[string]bool)
	for _, v := range vars {
		m[v] = true
	}
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return
}

func isMaplikeCond(e []Exp) (string, bool) {
	out := true
	if len(e) < 2 {
		out = false
	}
	ops := make(map[string]bool)
	for i, a := range e {
		if i == 0 {
			continue
		}
		test := car(a)
		if i == len(e)-1 {
			if String(test) != "'t" {
				out = false
			}
		} else {
			first := car(test)
			second := cadr(test)
			third := caddr(test)
			// need pattern (eq atom quote); i.e., atom==quote
			if String(first) != "eq" {
				out = false
			}
			op, ok := second.(string)
			if !ok {
				out = false
			}
			ops[op] = true
			if String(car(third)) != "quote" {
				out = false
			}
		}
	}
	if len(ops) != 1 {
		out = false
	}
	var op string
	for k := range ops {
		op = k
	}
	return op, out
}

func LabelExpr(defun Exp) (Exp, error) {
	if String(car(defun)) != "defun" {
		return nil, fmt.Errorf("not a defun: %s", String(defun))
	}
	name := cadr(defun)
	args := caddr(defun)
	body := cadddr(defun)
	q := func(s string) Exp {
		return s
	}
	nl := func(args ...Exp) Exp {
		return args
	}
	e := nl(
		q("label"),
		name,
		nl(
			q("lambda"),
			args,
			body,
		),
	)
	return e, nil
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
