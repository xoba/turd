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
