package scr

import (
	"encoding/json"
	"fmt"
	"strings"
)

type node struct {
	Value    string  `json:"V,omitempty"`
	Children []*node `json:"C,omitempty"`
}

func (n *node) Expression() (*Expression, error) {
	if len(n.Value) > 0 {
		e := NewString(n.Value)
		return &e, nil
	}
	var list []*Expression
	for _, c := range n.Children {
		e, err := c.Expression()
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return NewList(list...), nil
}

func parse(s string) (*node, error) {
	toks, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	nodes, err := parseTokens(toks)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func parseTokens(list []string) (*node, error) {
	switch len(list) {
	case 0:
		return nil, fmt.Errorf("can't parse empty list")
	case 1:
		return &node{Value: list[0]}, nil
	default:
		if list[0] != "(" || list[len(list)-1] != ")" {
			return nil, fmt.Errorf("not a list: %q", list)
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
				c, err := parseTokens(current)
				if err != nil {
					return nil, err
				}
				out.Children = append(out.Children, c)
				current = current[:0]
			}
		}
		if indent != 0 {
			panic(fmt.Errorf("indent = %d", indent))
		}
		return &out, nil
	}
}

func (n node) String() string {
	buf, _ := json.Marshal(n)
	return string(buf)
}

func tokenize(s string) ([]string, error) {
	var out []string
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	for _, x := range strings.Fields(s) {
		x = strings.TrimSpace(x)
		if len(x) == 0 {
			continue
		}
		out = append(out, x)
	}
	return out, nil
}
