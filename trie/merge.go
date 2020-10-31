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

type Joiner interface {
	Join(meet, a, b *KeyValue) (*KeyValue, error)
}

// Join is a three-way merge
func Join(meet, a, b *Trie, merger Joiner) (*Trie, error) {
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
		// only b exists:
		return b, nil
	case a != nil && b == nil: // nil case (5+6)/8
		// only a exists:
		return a, nil
	case a != nil && b != nil: // nil cases (7+8)/8
		// advance meet to the join of a and b:
		if meet == nil {
			x, err := New()
			if err != nil {
				return nil, err
			}
			meet = x
		} else {
			meet = meet.Copy()
		}
		meet.MarkDirty()
		switch {
		case a.KeyValue == nil && b.KeyValue == nil:
		case a.KeyValue == nil && b.KeyValue != nil:
			meet.KeyValue = b.KeyValue
		case a.KeyValue != nil && b.KeyValue == nil:
			meet.KeyValue = a.KeyValue
		case a.KeyValue != nil && b.KeyValue != nil:
			var kv *KeyValue
			if meet != nil {
				kv = meet.KeyValue
			}
			join, err := merger.Join(kv, a.KeyValue, b.KeyValue)
			if err != nil {
				return nil, err
			}
			meet.KeyValue = join
		default:
			panic("illegal")
		}
		for i, m := range meet.Next {
			j, err := Join(m, a.Next[i], b.Next[i], merger)
			if err != nil {
				return nil, err
			}
			meet.Next[i] = j
		}
		if err := meet.update(); err != nil {
			return nil, err
		}
		return meet, nil
	default:
		panic("illegal")
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
	b = set(b, "x/1", "x/1 value")
	check(viz(b, "b"))

	j, err := Join(m, a, b, mergeFunc(func(m, a, b *KeyValue) (*KeyValue, error) {
		if !bytes.Equal(a.Hash, b.Hash) {
			return nil, fmt.Errorf("conflict")
		}
		return a, nil
	}))
	check(err)
	check(viz(j, "join"))
	return nil
}

type mergeFunc func(meet, a, b *KeyValue) (*KeyValue, error)

func (m mergeFunc) Join(meet, a, b *KeyValue) (*KeyValue, error) {
	return m(meet, a, b)
}
