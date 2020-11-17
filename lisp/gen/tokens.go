package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func tokenize(s string) ([]string, error) {
	reader := bufio.NewReader(strings.NewReader(s))
	var out []string
	current := new(bytes.Buffer)
	quote := new(bytes.Buffer)
	add := func(s string) {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			return
		}
		out = append(out, s)
	}
	reset := func(b *bytes.Buffer) {
		add(b.String())
		b.Reset()
	}
	var inQuote bool
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			reset(current)
			return out, nil
		} else if err != nil {
			return nil, err
		}
		if inQuote {
			if r == '"' {
				inQuote = false
				reset(quote)
			} else {
				next, _, err := reader.ReadRune()
				if err != nil {
					return nil, err
				}
				esc := string([]rune{r, next})
				switch esc {
				case `\"`:
					quote.WriteRune('"')
				case `\n`:
					quote.WriteRune('\n')
				default:
					quote.WriteRune(r)
					if err := reader.UnreadRune(); err != nil {
						return nil, err
					}
				}
			}

		} else {
			switch {
			case r == '"':
				reset(current)
				inQuote = true
			case r == '(':
				reset(current)
				add("(")
			case r == '\'':
				reset(current)
				add("'")
			case r == ')':
				reset(current)
				add(")")
			case unicode.IsSpace(r):
				reset(current)
			default:
				current.WriteRune(r)
			}
		}
	}
	reset(current)
	return out, nil
}

func parseTokens(list []string) (Exp, error) {
	if len(list) == 0 {
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

type List []Exp

type stack []List

func (s *stack) push(b List) {
	*s = append(*s, b)
}

func (s *stack) len() int {
	return len(*s)
}

func (s *stack) pop() List {
	n := s.len()
	if n == 0 {
		return nil
	}
	x := (*s)[n-1]
	*s = (*s)[:n-1]
	return x
}
