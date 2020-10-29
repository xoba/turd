package trie

import (
	"encoding/asn1"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"sort"
	"strings"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/gviz"
	"github.com/xoba/turd/thash"
)

func Run(c cnfg.Config) error {
	return TestCOW()
}

// returned or wrapped to indicate key not found
var NotFound = errors.New("key not found")

type Database interface {
	Set([]byte, []byte) (Database, error)
	Get([]byte) ([]byte, error)
	Delete([]byte) (Database, error)
	Hash() ([]byte, error)
	Searchable
	Visualizable
}

type Visualizable interface {
	ToGviz(string) error
}

type Searchable interface {
	Search(SearchFunc) (*KeyValue, error)
}

type SearchFunc func(kv *KeyValue) bool

type KeyValue struct {
	Key, Value []byte
}

func New() (*Trie, error) {
	var t Trie
	if err := t.computeHash(); err != nil {
		return nil, err
	}
	return &t, nil
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

func (t *Trie) Hash() ([]byte, error) {
	if err := t.computeHash(); err != nil {
		return nil, err
	}
	return t.Merkle, nil
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

func TestPaths(cnfg.Config) error {
	t, err := New()
	if err != nil {
		return err
	}
	db := NewStrings(t)
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

// recursively computes and sets merkles if unset
func (t *Trie) computeHash() error {
	if len(t.Merkle) > 0 {
		return nil
	}
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
		if err := x.computeHash(); err != nil {
			return err
		}
		add(x.Merkle)
	}
	buf, err := asn1.Marshal(list)
	if err != nil {
		return err
	}
	t.Merkle = thash.Hash(buf)
	return nil
}

func String(t Searchable) string {
	var list []string
	t.Search(func(kv *KeyValue) bool {
		list = append(list, fmt.Sprintf("%q:%q", string(kv.Key), string(kv.Value)))
		return false
	})
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	return strings.Join(list, ", ")
}

func (t *Trie) String() string {
	return String(t)
}

// TODO: "Do" should take a range argument, and handler should return bool
func (t *Trie) Search(f SearchFunc) (*KeyValue, error) {
	for _, x := range t.Next {
		if x == nil {
			continue
		}
		if x.KeyValue != nil {
			if f(x.KeyValue) {
				return x.KeyValue, nil
			}
		}
		r, err := x.Search(f)
		if err != nil {
			return nil, err
		}
		if r != nil {
			return r, nil
		}
	}
	return nil, NotFound
}

func (t *Trie) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("nil key")
	}
	current := t
	for _, b := range key {
		if x := current.Next[b]; x != nil {
			current = x
		} else {
			return nil, NotFound
		}
	}
	if current.KeyValue == nil {
		return nil, NotFound
	}
	return current.KeyValue.Value, nil
}

func (t *Trie) Set(key, value []byte) (Database, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("nil key")
	}
	d, err := t.set(key, &KeyValue{Key: key, Value: value})
	if err != nil {
		return nil, err
	}
	return d, nil
}

func TestCOW() error {
	var db Database
	if true {
		t, err := New()
		if err != nil {
			return err
		}
		db = t
	} else {
		db = make(mapdb)
	}
	s := NewStrings(db)
	var steps []StringDatabase
	show := func(db StringDatabase) error {
		h, err := db.Hash()
		if err != nil {
			return err
		}
		fmt.Printf("%x; %v\n", h, db)
		return nil
	}
	step := func(key, value string) error {
		s2, err := s.Set(key, value)
		if err != nil {
			return err
		}
		steps = append(steps, s2)
		s = s2
		show(s)
		return nil
	}
	log := func(e error) {
		if e != nil {
			log.Fatal(e)
		}
	}
	log(step("/a", "b"))
	log(step("/ab", "xyz"))
	log(step("/a/x", "c"))
	for i := 0; i < 5; i++ {
		log(step(fmt.Sprintf("/a/x/%d", i), fmt.Sprintf("howdy %d", i)))
	}
	for _, db := range steps {
		log(show(db))
	}
	return s.ToGviz("trie.svg")
}

func (t *Trie) ToGviz(file string) error {
	return ToGviz(t, file)
}

func (t *Trie) set(key []byte, kv *KeyValue) (*Trie, error) {
	t = t.Copy()
	if len(key) == 0 {
		t.KeyValue = kv
	} else {
		b0 := key[0]
		child := t.Next[b0]
		if child == nil {
			c, err := New()
			if err != nil {
				return nil, err
			}
			child = c
		}
		c, err := child.set(key[1:], kv)
		if err != nil {
			return nil, err
		}
		t.Next[b0] = c
	}
	if err := t.computeHash(); err != nil {
		return nil, err
	}
	return t, nil
}

func inc(i *big.Int, v int) *big.Int {
	if i == nil {
		i = big.NewInt(0)
	}
	var x big.Int
	x.Add(i, big.NewInt(int64(v)))
	return &x
}

func ToGviz(g gviz.Graph, file string) error {
	gv, err := gviz.Compile(g, nil)
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
	return t.nodes(nil)
}

func (t *Trie) nodes(parent []byte) (out []gviz.Node) {
	n := node{
		id: fmt.Sprintf("%x", t.Merkle),
	}
	if t.KeyValue == nil {
		n.shape = "ellipse"
		n.label = string(parent)
	} else {
		n.label = fmt.Sprintf("%s = %s", string(t.KeyValue.Key), string(t.KeyValue.Value))
	}
	if n.label == "" {
		n.label = "nil"
	}
	out = append(out, n)
	for i, n := range t.Next {
		if n == nil {
			continue
		}
		child := byte(i)
		p2 := parent
		p2 = append(p2, child)
		out = append(out, n.nodes(p2)...)
	}
	return
}

type node struct {
	id, label, shape string
}

func (n node) ID() string {
	return n.id
}
func (n node) Label() string {
	return n.label
}

func (t node) Group() string {
	return ""
}
func (t node) Shape() string {
	return t.shape
}

func (t *Trie) Edges() (out []gviz.Edge) {
	id := func(t *Trie) string {
		return fmt.Sprintf("%x", t.Merkle)
	}
	for _, n := range t.Next {
		if n == nil {
			continue
		}
		out = append(out, edge{
			from: id(t),
			to:   id(n),
		})
		out = append(out, n.Edges()...)
	}
	return
}

type edge struct {
	from, to string
}

func (e edge) From() string {
	return e.from
}
func (e edge) To() string {
	return e.to
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

// TODO: this has to be COW too
func (t *Trie) Delete(key []byte) (Database, error) {
	panic("unimplemented")
	/*
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
		t.Prune()*/
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
