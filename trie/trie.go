package trie

import (
	"bytes"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/thash"
)

type Trie struct {
	kv   *KeyValue
	hash []byte
	next [256]*Trie
}

func New() *Trie {
	return &Trie{}
}

type Database interface {
	Set([]byte, []byte)
	Get([]byte) ([]byte, bool)
	Delete([]byte)
	Do(func(kv *KeyValue))
	// hash is implementation-dependent, but based only on key/values in the db
	ComputeHash() []byte
}

type StringDatabase interface {
	Set(string, string)
	Get(string) (string, bool)
	Delete(string)
	Do(func(kv *StringKeyValue))
	// hash is implementation-dependent, but based only on key/values in the db
	ComputeHash() []byte
}

type stringDB struct {
	hashLen int
	x       Database
}

func (db stringDB) Set(key string, value string) {
	db.x.Set([]byte(key), []byte(value))
}

func (db stringDB) Get(key string) (string, bool) {
	v, ok := db.x.Get([]byte(key))
	return string(v), ok
}

func (db stringDB) Delete(key string) {
	db.x.Delete([]byte(key))
}
func (db stringDB) Do(f func(kv *StringKeyValue)) {
	db.x.Do(func(kv *KeyValue) {
		f(&StringKeyValue{
			Key:   string(kv.Key),
			Value: string(kv.Value),
		})
	})
}
func (db stringDB) ComputeHash() []byte {
	h := db.x.ComputeHash()
	if n := db.hashLen; n > 0 {
		h = h[:n]
	}
	return h
}

type Iterable interface {
	Do(func(kv *KeyValue))
}

type mapdb map[string][]byte

func (m mapdb) Set(key []byte, value []byte) {
	m[string(key)] = value
}

func (m mapdb) Get(key []byte) ([]byte, bool) {
	v, ok := m[string(key)]
	return v, ok
}

func (m mapdb) Delete(key []byte) {
	delete(m, string(key))
}

func (m mapdb) Do(f func(*KeyValue)) {
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
		f(kv)
	}
}

func String(i Iterable) string {
	var list []*KeyValue
	i.Do(func(kv *KeyValue) {
		list = append(list, kv)
	})
	buf, _ := json.Marshal(list)
	return string(buf)
}

func (m mapdb) ComputeHash() []byte {
	return thash.Hash([]byte(String(m)))
}

func Run(cnfg.Config) error {
	var db1, db2 Database
	db1 = make(mapdb)
	db2 = New()
	if err := CheckDelete("map", stringDB{hashLen: 8, x: db1}); err != nil {
		return err
	}
	if err := CheckDelete("trie", stringDB{hashLen: 8, x: db2}); err != nil {
		return err
	}
	return nil
}

func CheckDelete(name string, db StringDatabase) error {
	db.Set("a", "x")
	db.Set("b", "y")
	first := db.ComputeHash()
	fmt.Printf("%s hash = %x\n", name, first)
	db.Set("c", "z")
	fmt.Printf("%s hash = %x\n", name, db.ComputeHash())
	db.Delete("c")
	second := db.ComputeHash()
	fmt.Printf("%s hash = %x\n", name, second)
	if !bytes.Equal(first, second) {
		return fmt.Errorf("%s hash mismatch: %x vs %x", name, first, second)
	}
	return nil
}

func Run3(cnfg.Config) error {
	const (
		n = 3
		m = 10000
	)
	newBuf := func() []byte {
		buf := make([]byte, 1+rand.Intn(n))
		rand.Read(buf)
		return buf
	}
	var db1, db2 Database
	db1 = make(mapdb)
	db2 = New()
	var keys [][]byte
	for i := 0; i < m; i++ {
		k, v := newBuf(), newBuf()
		keys = append(keys, k)
		db1.Set(k, v)
		db2.Set(k, v)
		if rand.Intn(10) == 0 {
			key := keys[rand.Intn(len(keys))]
			db1.Delete(key)
			db2.Delete(key)
		}
	}
	if String(db1) != String(db2) {
		return fmt.Errorf("mismatch")
	}
	return nil
}

func Run2(cnfg.Config) error {
	rand.Seed(0)
	const (
		idlen = 16
		n     = 7000
	)
	var ids []KeyValue
	{
		for i := 0; i < n; i++ {
			buf := make([]byte, idlen)
			rand.Read(buf)
			ids = append(ids, KeyValue{
				Key:   buf,
				Value: buf, // buf[:1+rand.Intn(idlen-1)],
			})
		}
	}
	if true {
		rand.Seed(time.Now().UTC().UnixNano())
		rand.Shuffle(len(ids), func(i, j int) {
			ids[i], ids[j] = ids[j], ids[i]
		})
	}
	t := New()
	all := make(map[string]string)
	m := make(map[string]bool)
	start := time.Now()
	for _, kv := range ids {
		id, prefix := kv.Key, kv.Value
		all[string(id)] = string(prefix)
		m[string(prefix)] = true
		t.Set(prefix, id)
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
	fmt.Printf("trie hash: 0x%x\n", t.ComputeHash())
	fmt.Printf("%v per iteration\n", time.Since(start)/n)
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (t *Trie) ComputeHash() []byte {
	if len(t.hash) > 0 {
		return t.hash
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
		if x == nil {
			continue
		}
		add([]byte{byte(i)})
		add(x.ComputeHash())
	}
	buf, err := asn1.Marshal(list)
	check(err)
	t.hash = thash.Hash(buf)
	return t.hash
}

type KeyValue struct {
	Key, Value []byte
}

type StringKeyValue struct {
	Key, Value string
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

func (t *Trie) Set(key, value []byte) {
	t.hash = nil
	current := t
	for _, b := range key {
		if x := current.next[b]; x != nil {
			current = x
		} else {
			newNode := New()
			current.next[b] = newNode
			current = newNode
		}
		current.hash = nil
	}
	current.kv = &KeyValue{Key: key, Value: value}
}

func (t *Trie) Delete(key []byte) {
	t.hash = nil
	current := t
	for _, b := range key {
		if x := current.next[b]; x != nil {
			current = x
		} else {
			return
		}
		current.hash = nil
	}
	current.kv = nil
}
