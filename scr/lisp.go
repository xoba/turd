package scr

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xoba/turd/cnfg"
)

func Lisp(cnfg.Config) error {
	return testRead(`(() (defun 'foo (a b c d) (+ a b c d)))
`)
}

func parseList(list []string, prefix string) node {
	fmt.Printf("%sgot %d %q\n", prefix, len(list), list)
	switch len(list) {
	case 0:
		panic("illegal")
	case 1:
		return node{Value: list[0]}
	default:
		if list[0] != "(" || list[len(list)-1] != ")" {
			panic(fmt.Sprintf("not a list: %q", list))
		}
		list = list[1 : len(list)-1]
		var out node
		var indent int
		var current []string
		for _, x := range list {
			switch {
			case x == "(":
				indent++
			case x == ")":
				indent--
			}
			current = append(current, x)
			if indent == 0 {
				out.Children = append(out.Children, parseList(current, prefix+"  "))
				current = current[:0]
			}
		}
		if indent != 0 {
			panic(fmt.Errorf("indent = %d", indent))
		}
		return out
	}
}

func testRead(in string) error {
	e, err := Read(in)
	if err != nil {
		return err
	}
	if out := e.String(); out != in {
		return fmt.Errorf("expected %q, got %q\n", in, out)
	}
	return nil
}

type node struct {
	Value    string `json:"V,omitempty"`
	Children []node `json:"C,omitempty"`
}

func (n node) String() string {
	buf, _ := json.Marshal(n)
	return string(buf)
}

func (n node) Express() (*Expression, error) {
	panic("unimplemented")
}

func parse(s string) (*node, error) {
	list := toList(s)
	fmt.Println(parseList(list, ""))

	return nil, fmt.Errorf("read unimplemented")
}

func toList(s string) (list []string) {
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	for _, x := range strings.Fields(s) {
		x = strings.TrimSpace(x)
		if len(x) == 0 {
			continue
		}
		list = append(list, x)
	}
	return
}

func Read(s string) (*Expression, error) {
	n, err := parse(s)
	if err != nil {
		return nil, err
	}
	fmt.Printf("node: %s\n", n)
	return nil, fmt.Errorf("read unimplemented")
}

func (e Expression) Car() (*Expression, error) {
	if e.IsAtom() {
		return nil, fmt.Errorf("can't car an atom")
	}
	if e.List.Empty() {
		return nil, fmt.Errorf("can't car an empty list")
	}
	return e.List.First(), nil
}

func (e Expression) Cdr() (*Expression, error) {
	if e.IsAtom() {
		return nil, fmt.Errorf("can't car an atom")
	}
	if e.List.Empty() {
		return nil, fmt.Errorf("can't car an empty list")
	}
	return e.List.First(), nil
}

func Eval(e *Expression) (*Expression, error) {
	if e.IsList() {

	}
	return nil, fmt.Errorf("eval unimplemented")
}

// either an atom or a list
type Expression struct {
	*Atom
	*List
}

func NewString(s string) Expression {
	return Expression{
		Atom: &Atom{
			Type: "string",
			Blob: []byte(s),
		},
	}
}

func NewBlob(s []byte) Expression {
	return Expression{
		Atom: &Atom{
			Type: "blob",
			Blob: s,
		},
	}
}

func NewList(list ...Expression) Expression {
	var z List
	e := Expression{
		List: &z,
	}
	for _, x := range list {
		*e.List = append(*e.List, x)
	}
	return e
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

func (e Expression) IsAtom() bool {
	e.Check()
	return e.Atom != nil
}

func (e Expression) IsList() bool {
	e.Check()
	return e.List != nil
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

type List []Expression

func (l List) Empty() bool {
	return len(l) == 0
}

func (l List) First() *Expression {
	return &l[0]
}

func (l List) Rest() List {
	return l[1:]
}

func (l *List) Add(e Expression) {
	*l = append(*l, e)
}
