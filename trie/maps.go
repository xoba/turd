package trie

import (
	"bytes"
	"sort"

	"github.com/xoba/turd/gviz"
	"github.com/xoba/turd/thash"
)

type mapdb map[string][]byte

func (m mapdb) copy() mapdb {
	x := make(mapdb)
	for k, v := range m {
		x[k] = v
	}
	return x
}

func (m mapdb) Set(key []byte, value []byte) Database {
	x := m.copy()
	x[string(key)] = value
	return x
}

func (m mapdb) Get(key []byte) ([]byte, bool) {
	v, ok := m[string(key)]
	return v, ok
}

func (m mapdb) Delete(key []byte) {
	delete(m, string(key))
}

func (m mapdb) Search(f SearchFunc) *KeyValue {
	var list []*KeyValue
	for k, v := range m {
		list = append(list, &KeyValue{
			Key:   []byte(k),
			Value: v,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return bytes.Compare(list[i].Key, list[j].Key) == -1
	})
	for _, kv := range list {
		if f(kv) {
			return kv
		}
	}
	return nil
}

func (m mapdb) Stats() *Stats {
	var s Stats
	s.IncCount(len(m))
	for k, v := range m {
		s.IncSize(len(k))
		s.IncSize(len(v))
	}
	return &s
}

func (m mapdb) String() string {
	return String(m)
}

func (m mapdb) Hash() []byte {
	return thash.Hash([]byte(String(m)))
}

func (m mapdb) Nodes() []gviz.Node {
	return nil
}
func (m mapdb) Edges() []gviz.Edge {
	return nil
}
