package scr

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func parseList(list []string) node {
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
				out.Children = append(out.Children, parseList(current))
				current = current[:0]
			}
		}
		if indent != 0 {
			panic(fmt.Errorf("indent = %d", indent))
		}
		return out
	}
}

type node struct {
	Value    string `json:"V,omitempty"`
	Children []node `json:"C,omitempty"`
}

func (n node) String() string {
	buf, _ := json.Marshal(n)
	return string(buf)
}

func (n node) Expression() (*Expression, error) {
	panic("unimplemented")
}

func parse(s string) (*node, error) {
	list := toList(s)
	fmt.Println(parseList(list))

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
