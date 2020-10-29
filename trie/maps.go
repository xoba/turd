package trie

import (
	"bytes"
	"fmt"
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

func (m mapdb) Set(key []byte, value []byte) (Database, error) {
	x := m.copy()
	x[string(key)] = value
	return x, nil
}

func (m mapdb) Get(key []byte) ([]byte, error) {
	v, ok := m[string(key)]
	if ok {
		return v, nil
	}
	return nil, NotFound
}

func (m mapdb) Delete(key []byte) (Database, error) {
	sk := string(key)
	if _, ok := m[sk]; ok {
		x := m.copy()
		delete(x, sk)
		return x, nil
	}
	return m, NotFound
}

func (m mapdb) Search(f SearchFunc) (*KeyValue, error) {
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
			return kv, nil
		}
	}
	return nil, NotFound
}

func (m mapdb) String() string {
	return String(m)
}

func (m mapdb) Hash() ([]byte, error) {
	return thash.Hash([]byte(String(m))), nil
}

func (m mapdb) Nodes() (out []gviz.Node) {
	for k, v := range m {
		out = append(out, node{
			id:    k,
			label: fmt.Sprintf("%s = %s", k, string(v)),
		})
	}
	return
}

func (m mapdb) Edges() (out []gviz.Edge) {
	var list []string
	for k := range m {
		list = append(list, k)
	}
	sort.Strings(list)
	var last string
	for _, x := range list {
		if last != "" {
			out = append(out, edge{
				from: last,
				to:   x,
			})
		}
		last = x
	}
	return
}

func (m mapdb) ToGviz(file string) error {
	return ToGviz(m, file)
}
