package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func tokenize(s string) ([]string, error) {
	var out []string
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	s = strings.Replace(s, "'", " ' ", -1)
	for _, x := range strings.Fields(s) {
		x = strings.TrimSpace(x)
		if len(x) == 0 {
			continue
		}
		out = append(out, x)
	}
	return out, nil
}

func NewNode(s string) (Exp, error) {
	{
		uncommented := func(s string) string {
			z := new(bytes.Buffer)
			for _, r := range s {
				if r == ';' {
					break
				}
				z.WriteRune(r)
			}
			return z.String()
		}
		w := new(bytes.Buffer)
		scanner := bufio.NewScanner(strings.NewReader(s))
		for scanner.Scan() {
			fmt.Fprintf(w, "%s\n", uncommented(scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		s = w.String()
	}
	toks, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	fmt.Printf("toks = %q\n", toks)
	nodes, err := parseTokens(toks)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func parseTokens(list []string) (Exp, error) {
	switch n := len(list); n {
	case 0:
		return nil, fmt.Errorf("can't parse empty list")
	case 1:
		return list[0], nil
	default:
		if list[0] != "(" || list[n-1] != ")" {
			return nil, fmt.Errorf("not a list: %q", list)
		}
		list = list[1 : n-1]
		var out []Exp
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
				out = append(out, c)
				current = current[:0]
			}
		}
		if indent != 0 {
			return nil, fmt.Errorf("indent = %d", indent)
		}
		return out, nil
	}
}
