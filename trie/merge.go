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
	p := func(t *Trie) string {
		if t == nil {
			return "nil"
		}
		return fmt.Sprintf("%x (%v)", t.Merkle[:2], t)
	}

	if false {
		fmt.Printf("meet = %s, a = %s, b = %s\n", p(meet), p(a), p(b))
	}

	both := func(t *Trie) (*Trie, error) {
		t.Merkle = nil
		switch {
		case a.KeyValue == nil && b.KeyValue == nil:
		case a.KeyValue == nil && b.KeyValue != nil:
			t.KeyValue = b.KeyValue
		case a.KeyValue != nil && b.KeyValue == nil:
			t.KeyValue = a.KeyValue
		case a.KeyValue != nil && b.KeyValue != nil:
			if !bytes.Equal(a.KeyValue.Hash, b.KeyValue.Hash) {
				return nil, fmt.Errorf("key value conflict: %v vs %v", a.KeyValue, b.KeyValue)
			}
		default:
			panic("illegal")
		}
		for i, m2 := range t.Next {
			j, err := Join(m2, a.Next[i], b.Next[i])
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

	switch {
	case eq(a, b):
		// no conflict whatsoever:
		return a, nil
	case b == nil || eq(meet, b):
		// only a is new or changed:
		return a, nil
	case a == nil || eq(meet, a):
		// only b is new or changed:
		return b, nil
	case a == nil && b != nil: // nil case (3+4)/8
		return b, nil
	case a != nil && b == nil: // nil case (5+6)/8
		return a, nil
	case a != nil && b != nil: // nil cases (7+8)/8
		fmt.Printf("%x vs %x\n", a.Merkle[:2], b.Merkle[:2])
		var x *Trie
		if meet == nil {
			t, err := New()
			if err != nil {
				return nil, err
			}
			x = t
		} else {
			x = meet.Copy()
		}
		return both(x)
	default:
		panic("default")
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
	if true {
		a = set(a, "x/1", "x/1 value")
		a = set(a, "x/2", "x2 value")
		a = set(a, "x/3", "x3 value")
		a = set(a, "y/1", "y1 value")
	}
	check(viz(a, "a"))

	b = set(m, "y", "y value")
	b = set(m, "x/1", "x/1 value")
	check(viz(b, "b"))

	j, err := Join(m, a, b)
	check(err)
	check(viz(j, "join"))
	return nil
}
