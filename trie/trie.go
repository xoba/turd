package trie

import (
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/gviz"
	"github.com/xoba/turd/thash"
)

type Database interface {
	Set([]byte, []byte) Database
	Get([]byte) ([]byte, bool)
	Delete([]byte)
	Hash() []byte
	Searchable
}

type SearchFunc func(kv *KeyValue) bool

type KeyValue struct {
	Key, Value []byte
}

func New() *Trie {
	t := &Trie{}
	t.Merkle = computeHash(t)
	return t
}

// TODO: make it copy-on-write
type Trie struct {
	*KeyValue
	Merkle []byte
	Next   [256]*Trie
}

// Copy without merkle hash
func (t *Trie) Copy() *Trie {
	return &Trie{
		KeyValue: t.KeyValue,
		Next:     t.Next,
	}
}

type Stats struct {
	Count, Size *big.Int
}

type Searchable interface {
	Search(SearchFunc) *KeyValue
}

func (t *Trie) Hash() []byte {
	return t.Merkle
}

func (t *Trie) IsDirty() bool {
	return len(t.Merkle) == 0
}

func (t *Trie) IsClean() bool {
	return !t.IsDirty()
}

func (t *Trie) MarkDirty() {
	t.Merkle = nil
}

func (s *Stats) Inc(o *Stats) {
	if s.Count == nil {
		s.Count = big.NewInt(0)
	}
	if s.Size == nil {
		s.Size = big.NewInt(0)
	}
	s.Count.Add(s.Count, o.Count)
	s.Size.Add(s.Size, o.Size)
}

func (s *Stats) IncCount(i int) {
	s.Count = inc(s.Count, i)
}

func (s *Stats) IncSize(i int) {
	s.Size = inc(s.Size, i)
}

func String(i Searchable) string {
	var list []*KeyValue
	i.Search(func(kv *KeyValue) bool {
		list = append(list, kv)
		return false
	})
	buf, _ := json.Marshal(list)
	return string(buf)
}

func TestPaths(cnfg.Config) error {
	db := NewStrings(New())
	add := func(p string) {
		db.Set(p, "value for "+p)
	}
	add("/a")
	add("/a/x")
	add("/a/y")
	add("/a/z")
	add("/a/z/123")
	add("/b")
	fmt.Println(db)

	for _, key := range strings.Split("/a,/a/q,/a/z", ",") {
		r, ok := db.Get(key)
		fmt.Printf("get(%q) = %q, %v\n", key, r, ok)
	}

	return nil
}

func Run(c cnfg.Config) error {
	return TestCOW()

	if err := TestPaths(c); err != nil {
		return err
	}
	if err := Run2(c); err != nil {
		return err
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
				Value: buf,
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
	fmt.Printf("trie hash: 0x%x\n", t.Merkle)
	fmt.Printf("%v per iteration\n", time.Since(start)/n)
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// not recursive, doesn't change Trie
func computeHash(t *Trie) []byte {
	var list [][]byte
	add := func(x interface{}) {
		switch t := x.(type) {
		case nil:
			list = append(list, nil)
		case []byte:
			if len(t) == 0 {
				panic("empty byte buffer")
			}
			list = append(list, t)
		case int:
			list = append(list, big.NewInt(int64(t)).Bytes())
		default:
			panic(fmt.Errorf("unsupported type %T", t))
		}
	}
	if t.KeyValue == nil {
		add(nil)
	} else {
		add(t.KeyValue.Key)
		add(t.KeyValue.Value)
	}
	for i, x := range t.Next {
		if x == nil {
			continue
		}
		add(i)
		add(x.Merkle)
	}
	buf, err := asn1.Marshal(list)
	check(err)
	return thash.Hash(buf)
}

func (t *Trie) String() string {
	var list []string
	t.Search(func(kv *KeyValue) bool {
		list = append(list, fmt.Sprintf("%q:%q", string(kv.Key), string(kv.Value)))
		return false
	})
	return strings.Join(list, ", ")
}

// TODO: "Do" should take a range argument, and handler should return bool
func (t *Trie) Search(f SearchFunc) *KeyValue {
	for _, x := range t.Next {
		if x == nil {
			continue
		}
		if x.KeyValue != nil {
			if f(x.KeyValue) {
				return x.KeyValue
			}
		}
		if r := x.Search(f); r != nil {
			return r
		}
	}
	return nil
}

func (t *Trie) Get(key []byte) ([]byte, bool) {
	if len(key) == 0 {
		return nil, false
	}
	current := t
	for _, b := range key {
		if x := current.Next[b]; x != nil {
			current = x
		} else {
			return nil, false
		}
	}
	if current.KeyValue == nil {
		return nil, false
	}
	return current.KeyValue.Value, true
}

func (t *Trie) Set(key, value []byte) Database {
	return t.set(key, &KeyValue{Key: key, Value: value})
}

func TestCOW() error {
	t := New()
	s := NewStrings(t)
	var steps []StringDatabase
	show := func(db StringDatabase) {
		fmt.Printf("%x; %v\n", db.Hash(), db)
	}
	step := func(key, value string) {
		s2 := s.Set(key, value)
		steps = append(steps, s2)
		s = s2
		show(s)
	}
	step("a", "b")
	step("a/x", "c")
	for i := 0; i < 5; i++ {
		step(fmt.Sprintf("a/x/%d", i), fmt.Sprintf("howdy %d", i))
	}
	for _, db := range steps {
		show(db)
	}

	return t.ToGviz("trie.svg")
}

func (t *Trie) set(key []byte, kv *KeyValue) *Trie {
	t = t.Copy()
	if len(key) == 0 {
		t.KeyValue = kv
	} else {
		b0 := key[0]
		child := t.Next[b0]
		if child == nil {
			child = New()
		}
		t.Next[b0] = child.set(key[1:], kv)
	}
	t.Merkle = computeHash(t)
	return t
}

func (t *Trie) OLD_Set(key, value []byte) {
	t.MarkDirty()
	current := t
	for _, b := range key {
		if x := current.Next[b]; x != nil {
			current = x
		} else {
			newNode := New()
			current.Next[b] = newNode
			current = newNode
		}
		current.MarkDirty()
	}
	current.KeyValue = &KeyValue{Key: key, Value: value}
}

func inc(i *big.Int, v int) *big.Int {
	if i == nil {
		i = big.NewInt(0)
	}
	var x big.Int
	x.Add(i, big.NewInt(int64(v)))
	return &x
}

func (t *Trie) Delete(key []byte) {
	t.MarkDirty()
	current := t
	for _, b := range key {
		if x := current.Next[b]; x != nil {
			current = x
		} else {
			return
		}
		current.MarkDirty()
	}
	current.KeyValue = nil
	t.Prune()
}

func (t *Trie) ToGviz(file string) error {
	gv, err := gviz.Compile(t, nil)
	if err != nil {
		return err
	}
	const in = "g.gv"
	if err := ioutil.WriteFile(in, gv, os.ModePerm); err != nil {
		return err
	}
	return gviz.Dot(in, file)
}

func (t *Trie) Nodes() (out []gviz.Node) {
	out = append(out, t)
	return
}

func (t *Trie) ID() string {
	return fmt.Sprintf("%x", t.Merkle)
}
func (t *Trie) Group() string {
	return ""
}
func (t *Trie) Shape() string {
	return ""
}
func (t *Trie) Label() string {
	return t.ID()[:4]
}

func (t *Trie) Edges() []gviz.Edge {
	return nil
}

// Prune returns true if this trie node can be discarded
func (t *Trie) Prune() bool {
	if t.IsClean() {
		return false
	}
	var children int
	for i, c := range t.Next {
		if c == nil {
			continue
		}
		children++
		if c.Prune() {
			t.Next[i] = nil
		}
	}
	if t.KeyValue != nil {
		return false
	}
	if children > 0 {
		return false
	}
	return true
}
