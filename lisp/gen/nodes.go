package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

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
	n := len(list)

	if n == 0 {
		return nil, fmt.Errorf("can't parse empty list of tokens")
	}

	s := new(stack)
	var current List
	for _, x := range list {
		switch x {
		case "'":
			return nil, fmt.Errorf("can't handle quote")
		case "(":
			s.push(current)
			current = make([]Exp, 0)
		case ")":
			x := s.pop()
			x = append(x, current)
			current = x
		default:
			current = append(current, x)
		}
	}
	return current, nil

}
