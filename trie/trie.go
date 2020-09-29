package trie

import (
	"encoding/asn1"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/thash"
)

type Trie struct {
	kv   *KeyValue
	b    *byte
	hash []byte
	next [256]*Trie
}

func New(b *byte) *Trie {
	return &Trie{b: b}
}

func Byte(i int) *byte {
	b := byte(i)
	return &b
}

func Run(cnfg.Config) error {
	rand.Seed(0)
	const idlen = 16
	newID := func() []byte {
		buf := make([]byte, idlen)
		rand.Read(buf)
		return buf
	}
	t := New(nil)
	all := make(map[string]string)
	m := make(map[string]bool)
	const n = 10
	start := time.Now()
	for i := 0; i < n; i++ {
		id := newID()
		prefix := id[:1+rand.Intn(idlen-1)]
		all[string(id)] = string(prefix)
		m[string(prefix)] = true
		if err := t.Set(prefix, id); err != nil {
			return err
		}
		h, err := t.ComputeHash()
		if err != nil {
			return err
		}
		fmt.Printf("hash = %x\n", h)
	}
	for id, prefix := range all {
		_, ok := t.Get([]byte(prefix))
		if !ok {
			return fmt.Errorf("doesn't have %q", prefix)
		}
		for i := 1; i < len(id); i++ {
			p2 := id[:i]
			if _, ok := t.Get([]byte(p2)); ok != m[p2] {
				return fmt.Errorf("mismatch")
			}
		}
	}
	h, err := t.ComputeHash()
	if err != nil {
		return err
	}
	fmt.Printf("trie hash: 0x%x\n", h)
	fmt.Printf("%v per iteration\n", time.Since(start)/n)
	return nil
}

func (t *Trie) ComputeHash() ([]byte, error) {
	if len(t.hash) > 0 {
		return t.hash, nil
	}
	if false && t.b != nil {
		fmt.Printf("ComputeHash(%x)\n", []byte{*t.b})
	}
	var list [][]byte
	add := func(x []byte) {
		list = append(list, x)
	}
	if t.kv != nil {
		add(t.kv.Key)
		add(t.kv.Value)
	}
	for i, x := range t.next {
		add([]byte{byte(i)})
		if x == nil {
			continue
		}
		h, err := x.ComputeHash()
		if err != nil {
			return nil, err
		}
		add(h)
	}
	buf, err := asn1.Marshal(list)
	if err != nil {
		return nil, err
	}
	t.hash = thash.Hash(buf)
	return t.hash, nil
}

type KeyValue struct {
	Key, Value []byte
}

func (kv KeyValue) ComputeHash() ([]byte, error) {
	buf, err := asn1.Marshal(kv)
	if err != nil {
		return nil, err
	}
	return thash.Hash(buf), nil
}

func (t *Trie) String() string {
	var list []string
	t.Do(func(kv *KeyValue) {
		list = append(list, fmt.Sprintf("%q:%q", string(kv.Key), string(kv.Value)))
	})
	return strings.Join(list, ", ")
}

func (t *Trie) Do(f func(kv *KeyValue)) {
	for _, x := range t.next {
		if x == nil {
			continue
		}
		if x.kv != nil {
			f(x.kv)
		}
		x.Do(f)
	}
}

// FIX: returns ("",true) for non-leaf nodes
func (t *Trie) Get(key []byte) ([]byte, bool) {
	if len(key) == 0 {
		return nil, false
	}
	current := t
	for _, b := range key {
		if x := current.next[b]; x != nil {
			current = x
		} else {
			return nil, false
		}
	}
	if current.kv == nil {
		return nil, false
	}
	return current.kv.Value, true
}

func (t *Trie) Set(key, value []byte) error {
	fmt.Printf("Set(%x, %x)\n", key, value)
	t.hash = nil
	current := t
	for _, b := range key {
		if x := current.next[b]; x != nil {
			current = x
		} else {
			newNode := New(Byte(int(b)))
			current.next[b] = newNode
			current = newNode
		}
		current.hash = nil
	}
	current.kv = &KeyValue{Key: key, Value: value}
	return nil
}
