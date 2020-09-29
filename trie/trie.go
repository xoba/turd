package trie

import (
	"encoding/json"
	"fmt"

	"github.com/xoba/turd/cnfg"
)

type Trie struct {
	key   string
	value string
	next  [256]*Trie
}

func New() *Trie {
	return &Trie{}
}

func Run(cnfg.Config) error {
	t := New()
	const (
		key   = "key"
		value = "value"
	)
	for i := 0; i < 10; i++ {
		t.Set(fmt.Sprintf("%s-%d", key, i), fmt.Sprintf("%s-%d", value, i))
	}
	fmt.Println(t.Get(key))
	for i := 0; i < 12; i++ {
		k := fmt.Sprintf("%s-%d", key, i)
		v, ok := t.Get(k)
		if ok {
			fmt.Printf("got %q -> %q\n", k, v)
		} else {
			fmt.Printf("no such key: %q\n", k)
		}
	}
	fmt.Println(t)
	return nil
}

type KeyValue struct {
	Key, Value string
}

func (t *Trie) String() string {
	var list []KeyValue
	t.Do(func(key, value string) {
		list = append(list, KeyValue{
			Key:   key,
			Value: value,
		})
	})
	buf, _ := json.MarshalIndent(list, "", "  ")
	return string(buf)
}

func (t *Trie) Do(f func(string, string)) {
	for _, x := range t.next {
		if x == nil {
			continue
		}
		if x.key != "" {
			f(x.key, x.value)
		}
		x.Do(f)
	}
}

// FIX: returns ("",true) for non-leaf nodes
func (t *Trie) Get(key string) (string, bool) {
	current := t
	for _, b := range []byte(key) {
		if x := current.next[b]; x != nil {
			current = x
		} else {
			return "", false
		}
	}
	return current.value, true
}

func (t *Trie) Set(key, value string) {
	current := t
	for _, b := range []byte(key) {
		if x := current.next[b]; x != nil {
			current = x
		} else {
			newNode := New()
			current.next[b] = newNode
			current = newNode
		}
	}
	current.key = key
	current.value = value
}
