package trie

import (
	"bytes"
	"fmt"
	"log"

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
	case eq(meet, a, b):
		// nothing changed
		return meet, nil
	case eq(meet, b):
		// a changed:
		return a, nil
	case eq(meet, a):
		// b changed:
		return b, nil
	default:
		// both changed:
		t, err := New()
		if err != nil {
			return nil, err
		}
		for i, m2 := range meet.Next {
			a2, b2 := a.Next[i], b.Next[i]
			j, err := Join(m2, a2, b2)
			if err != nil {
				return nil, err
			}
			t.Next[i] = j
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

	var m *Trie
	m = set(m, "a", "a value")
	m = set(m, "b", "b value")
	m = set(m, "c", "c value")
	fmt.Println(m)

	a := set(m, "x", "x value")
	fmt.Println(a)

	b := set(m, "y", "y value")
	fmt.Println(b)

	merge, err := Join(m, a, b)
	check(err)
	fmt.Println(merge)
	return nil
}
