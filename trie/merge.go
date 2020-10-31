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
		return b, nil
	case a != nil && b == nil: // nil case (5+6)/8
		return a, nil
	case a != nil && b != nil: // nil cases (7+8)/8
		var out *Trie
		if meet == nil {
			x, err := New()
			if err != nil {
				return nil, err
			}
			out = x
		} else {
			out = meet.Copy()
		}
		out.MarkDirty()
		switch {
		case a.KeyValue == nil && b.KeyValue == nil:
		case a.KeyValue == nil && b.KeyValue != nil:
			out.KeyValue = b.KeyValue
		case a.KeyValue != nil && b.KeyValue == nil:
			out.KeyValue = a.KeyValue
		case a.KeyValue != nil && b.KeyValue != nil:
			var mkv *KeyValue
			if meet != nil {
				mkv = meet.KeyValue
			}
			join, err := merger.Join(mkv, a.KeyValue, b.KeyValue)
			if err != nil {
				return nil, err
			}
			out.KeyValue = join
		default:
			panic("illegal")
		}
		for i, m2 := range out.Next {
			j, err := Join(m2, a.Next[i], b.Next[i], merger)
			if err != nil {
				return nil, err
			}
			out.Next[i] = j
		}
		if err := out.update(); err != nil {
			return nil, err
		}
		return out, nil
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
	b = set(m, "x/1", "x/1 value")
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
