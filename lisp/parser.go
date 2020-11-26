package lisp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/big"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Parse(s string) (Exp, error) {
	if !utf8.ValidString(s) {
		return nil, fmt.Errorf("invalid utf8")
	}
	u, err := uncomment(s)
	if err != nil {
		return nil, err
	}
	toks, err := tokenize(u)
	if err != nil {
		return nil, err
	}
	e, err := parseTokens(toks)
	if err != nil {
		return nil, err
	}
	x, ok := e.([]Exp)
	if !ok {
		return nil, fmt.Errorf("expression not a list")
	}
	if n := len(x); n != 1 {
		return nil, fmt.Errorf("need just one element, got %d", n)
	}
	return x[0], nil
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
	var s stack
	var current []Exp
	for _, x := range list {
		switch x {
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
	return expandQuotes(current)
}

func expandQuotes(e Exp) (Exp, error) {
	switch t1 := e.(type) {
	case string:
		if t1 == "'" {
			return nil, fmt.Errorf("unexpanded quote")
		}
		return t1, nil
	case []Exp:
		var out []Exp
		n := len(t1)
		for i := 0; i < n; i++ {
			c := t1[i]
			switch t2 := c.(type) {
			case string:
				if t2 == "'" {
					if i == n-1 {
						return nil, fmt.Errorf("bad quote")
					}
					x, err := expandQuotes(t1[i+1])
					if err != nil {
						return nil, err
					}
					c = []Exp{
						"quote",
						x,
					}
					i++
				}
			default:
				z, err := expandQuotes(c)
				if err != nil {
					return nil, err
				}
				c = z
			}
			out = append(out, c)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("illegal type %T", t1)
	}
}

type stack [][]Exp

func (s *stack) push(b []Exp) {
	*s = append(*s, b)
}

func (s *stack) len() int {
	return len(*s)
}

func (s *stack) pop() []Exp {
	n := s.len()
	if n == 0 {
		return nil
	}
	x := (*s)[n-1]
	*s = (*s)[:n-1]
	return x
}

// from the go spec
func GoIdentifiers() (list []string) {
	add := func(category, words string) {
		list = append(list, strings.Fields(words)...)
	}

	add("keywords", `break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
`)
	add("functions",
		`      append cap close complex copy delete imag len
       make new panic print println real recover
`)
	add("constants",
		`      true false iota
`)
	add("zero", "nil")
	add("types", `  bool byte complex64 complex128 error float32 float64
       int int8 int16 int32 int64 rune string
       uint uint8 uint16 uint32 uint64 uintptr
`)
	return
}

func sanitize(s string) string {
	return s + "_go_sanitized"
}

func UnsanitizeGo(e Exp) Exp {
	for _, x := range GoIdentifiers() {
		e = translateAtoms(e, sanitize(x), x)
	}
	return e
}

func SanitizeGo(e Exp) Exp {
	for _, x := range GoIdentifiers() {
		e = translateAtoms(e, x, sanitize(x))
	}
	return e
}

func translateAtoms(e Exp, from, to string) Exp {
	switch t := e.(type) {
	case string:
		if t == from {
			return to
		}
		return t
	case *big.Int, []byte:
		return t
	case []Exp:
		var out []Exp
		for _, c := range t {
			out = append(out, translateAtoms(c, from, to))
		}
		return out
	case error:
		return t
	default:
		return fmt.Errorf("can't translate %T %v", t, t)
	}
	return e
}
