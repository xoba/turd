package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func tokenize0(s string) ([]string, error) {
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

func uncomment(s string) (string, error) {
	w := new(bytes.Buffer)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		z := new(bytes.Buffer)
		for _, r := range line {
			if r == ';' {
				break
			}
			z.WriteRune(r)
		}
		fmt.Fprintf(w, "%s\n", z.String())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return w.String(), nil
}

func NewNode(s string) (Exp, error) {
	u, err := uncomment(s)
	if err != nil {
		return nil, err
	}
	s = u
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
