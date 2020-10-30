package trie

import (
	"bytes"
	"fmt"
	"log"

	"github.com/xoba/turd/cnfg"
)

func eq(list ...[]byte) bool {
	n := len(list)
	for i := 1; i < n; i++ {
		if !bytes.Equal(list[i-1], list[i]) {
			return false
		}
	}
	return true
}

func ThreeWayMerge(meet, a, b *Trie) *Trie {
	if eq(meet.Merkle, a.Merkle, b.Merkle) {
		return meet
	}

	panic("")
}

func TestMerge(cnfg.Config) error {
	set := func(t *Trie, key, value string) *Trie {
		if t == nil {
			x, err := New()
			if err != nil {
				log.Fatal(err)
			}
			t = x
		}
		x, err := t.Set([]byte(key), []byte(value))
		if err != nil {
			log.Fatal(err)
		}
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

	merge := ThreeWayMerge(m, a, b)
	fmt.Println(merge)
	return nil
}
