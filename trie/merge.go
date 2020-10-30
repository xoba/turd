package trie

import (
	"bytes"
	"fmt"
	"log"

	"github.com/skratchdot/open-golang/open"
	"github.com/xoba/turd/cnfg"
)

func eq(list ...*Trie) bool {
	m := func(t *Trie) []byte {
		if t == nil {
			return nil
		}
		return t.Merkle
	}
	n := len(list)
	for i := 1; i < n; i++ {
		if !bytes.Equal(m(list[i-1]), m(list[i])) {
			return false
		}
	}
	return true
}

// Join is a three-way merge
func Join(meet, a, b *Trie) (*Trie, error) {
	switch {
	case eq(a, b):
		// no conflict whatsoever:
		return a, nil
	case eq(meet, b):
		// only a changed:
		return a, nil
	case eq(meet, a):
		// only b changed:
		return b, nil
	default:
		// TODO: when meet doesn't exist?
		switch {
		case a == nil && b != nil:
			// a does not exist
			fmt.Println("*B")
			return b, nil
		case a != nil && b == nil:
			// b does not exist
			fmt.Println("*A")
			return a, nil
		}
		// both a & b changed:
		t := meet.Copy()
		for i, m2 := range t.Next {
			a2, b2 := a.Next[i], b.Next[i]
			j, err := Join(m2, a2, b2) // TODO: what if m2 == nil?
			if err != nil {
				return nil, err
			}
			t.Next[i] = j
		}
		if err := t.update(); err != nil {
			return nil, err
		}
		return t, nil
	}
}

func TestMerge(cnfg.Config) error {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	set := func(t *Trie, key, value string) *Trie {
		if t == nil {
			x, err := New()
			check(err)
			t = x
		}
		x, err := t.Set([]byte(key), []byte(value))
		check(err)
		return x.(*Trie)
	}

	viz := func(t *Trie, name string) error {
		fmt.Println(t)
		file := fmt.Sprintf("trie_%s.svg", name)
		if err := t.ToGviz(file, name); err != nil {
			return err
		}
		return open.Run(file)
	}

	var a, b, m *Trie

	m = set(m, "a", "a value")
	m = set(m, "b", "b value")
	m = set(m, "c", "c value")
	check(viz(m, "meet"))

	a = set(m, "x", "x value")
	a = set(a, "x/1", "x1 value")
	a = set(a, "x/2", "x2 value")
	a = set(a, "x/3", "x3 value")
	a = set(a, "y/1", "y1 value")
	check(viz(a, "a"))

	b = set(m, "y", "y value")
	//b = set(m, "x", "conflicting value")
	check(viz(b, "b"))

	j, err := Join(m, a, b)
	check(err)
	check(viz(j, "join"))
	return nil
}
