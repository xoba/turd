package trie

import (
	"bytes"
	"encoding/asn1"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"sort"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/gviz"
	"github.com/xoba/turd/thash"
)

type Trie struct {
	*KeyValue
	Merkle []byte
	Next   [256]*Trie
	Count  *big.Int // number of kv pairs in this and all descendants
}

func (t *Trie) IsClean() bool {
	return len(t.Merkle) > 0 && t.Count != nil
}

func Run(c cnfg.Config) error {
	return TestMerge(c)
	return TestCOW(c)
}

// returned or wrapped to indicate key not found
var NotFound = errors.New("key not found")
var EmptyKey = errors.New("empty key")

// copy-on-write database
type Database interface {
	Set([]byte, []byte) (Database, error)
	Get([]byte) ([]byte, error)
	Delete([]byte) (Database, error)
	Hash() ([]byte, error)
	Searchable
	Visualizable
}

type Visualizable interface {
	ToGviz(file, title string) error
}

type Searchable interface {
	// error indicates not found or worse, else keyvalue is non-nil
	Search(SearchFunc) (*KeyValue, error)
}

type SearchFunc func(kv *KeyValue) bool

type KeyValue struct {
	Key, Value, Hash []byte
}

func (k KeyValue) computeHash() ([]byte, error) {
	var h hasher
	h.add(k.Key)
	h.add(k.Value)
	return h.compute()
}

func (k *KeyValue) Verify() error {
	h, err := k.computeHash()
	if err != nil {
		return err
	}
	if !bytes.Equal(h, k.Hash) {
		return fmt.Errorf("hash mismatch")
	}
	return nil
}

func New() (*Trie, error) {
	var t Trie
	if err := t.update(); err != nil {
		return nil, err
	}
	return &t, nil
}

// Copy without merkle hash, in order to mark "dirty"
func (t *Trie) Copy() *Trie {
	return &Trie{
		KeyValue: t.KeyValue,
		Next:     t.Next,
		Count:    cp(t.Count),
	}
}

func (t *Trie) Hash() ([]byte, error) {
	if err := t.update(); err != nil {
		return nil, err
	}
	return t.Merkle, nil
}

type hasher [][]byte

func (h *hasher) add(x interface{}) error {
	switch t := x.(type) {
	case nil:
		*h = append(*h, []byte{})
	case []byte:
		*h = append(*h, t)
	case int:
		*h = append(*h, big.NewInt(int64(t)).Bytes())
	default:
		return fmt.Errorf("unsupported type %T: %v", t, t)
	}
	return nil
}

func (h hasher) compute() ([]byte, error) {
	buf, err := asn1.Marshal(h)
	if err != nil {
		return nil, err
	}
	return thash.Hash(buf), nil
}

// recursively computes and sets count and merkles if unset
func (t *Trie) update() error {
	if t.IsClean() {
		return nil
	}
	var h hasher
	if kv := t.KeyValue; kv == nil {
		if err := h.add(nil); err != nil {
			return err
		}
	} else {
		if err := kv.Verify(); err != nil {
			return err
		}
		if err := h.add(kv.Hash); err != nil {
			return err
		}
	}
	count := big.NewInt(0)
	if t.KeyValue != nil {
		count = inc(count, 1)
	}
	for i, x := range t.Next {
		if x == nil {
			continue
		}
		if err := h.add(i); err != nil {
			return err
		}
		if err := x.update(); err != nil {
			return err
		}
		if err := h.add(x.Merkle); err != nil {
			return err
		}
		count = sum(count, x.Count)
	}
	hash, err := h.compute()
	if err != nil {
		return err
	}
	t.Merkle = hash
	t.Count = count
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

func (t *Trie) Search(f SearchFunc) (*KeyValue, error) {
	for _, c := range t.Next {
		if c == nil {
			continue
		}
		if c.KeyValue != nil {
			if f(c.KeyValue) {
				return c.KeyValue, nil
			}
		}
		r, err := c.Search(f)
		if err != nil && err != NotFound {
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
		return nil, EmptyKey
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

// TODO: also need a set method etc that returns *Trie
func (t *Trie) Set(key, value []byte) (Database, error) {
	if len(key) == 0 {
		return nil, EmptyKey
	}
	kv := &KeyValue{Key: key, Value: value}
	h, err := kv.computeHash()
	if err != nil {
		return nil, err
	}
	kv.Hash = h
	return t.set(key, kv)
}

func TestCOW(c cnfg.Config) error {
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
	type step struct {
		db    StringDatabase
		title string
	}
	var steps []step
	add := func(title string) {
		steps = append(steps, step{
			db:    s,
			title: fmt.Sprintf("%d. %s", len(steps), title),
		})
	}
	add("empty")
	viz := func(s step, i int) error {
		file := fmt.Sprintf("trie_%d.svg", i)
		if err := s.db.ToGviz(file, s.title); err != nil {
			return err
		}
		return open.Run(file)
	}
	show := func(db StringDatabase) error {
		h, err := db.Hash()
		if err != nil {
			return err
		}
		fmt.Printf("%x; %v\n", h[:2], db)
		return nil
	}
	set := func(key, value string) error {
		s2, err := s.Set(key, value)
		if err != nil {
			return err
		}
		s = s2
		add("set " + key)
		show(s)
		return nil
	}
	del := func(key string) error {
		if x, err := s.Delete(key); err != nil {
			return err
		} else {
			s = x
		}
		add("del " + key)
		show(s)
		return nil
	}
	log := func(e error) {
		if e != nil {
			log.Fatal(e)
		}
	}
	log(set("/a", "b"))
	log(set("/c", "d"))
	log(set("/c/1/2/3/4", "long"))
	log(set("/ab", "xyz"))
	log(set("/a/x", "c"))
	if c.Delete {
		log(del("/c/1/2/3/4"))
	}
	for i := 0; i < 4; i++ {
		log(set(fmt.Sprintf("/a/x/%d", i), fmt.Sprintf("howdy %d", i)))
	}
	if c.Delete {
		log(del("/a/x"))
		log(del("/a/x/2")) // delete a leaf node
		log(del("/a"))
	}
	fmt.Println("replaying...")
	for i, s := range steps {
		log(show(s.db))
		if err := viz(s, i); err != nil {
			return err
		}
	}
	return nil
}

func (t *Trie) ToGviz(file, title string) error {
	return ToGviz(t, file, title)
}

func (t *Trie) Delete(key []byte) (Database, error) {
	if len(key) == 0 {
		return nil, EmptyKey
	}
	return t.del(key, key)
}

// TODO: need to keep children if any; i.e., this is not a recursive delete!
func (t *Trie) del(key, original []byte) (*Trie, error) {
	t = t.Copy()
	switch len(key) {
	case 0:
		if t.KeyValue == nil {
			return nil, fmt.Errorf("1: %w", NotFound)
		}
		t.KeyValue = nil
		if one(t.Count) {
			// this node needs to be pruned
			return nil, nil
		}
	default:
		prefix := key[0]
		c := t.Next[prefix]
		if c == nil {
			return nil, fmt.Errorf("2: %w", NotFound)
		}
		c2, err := c.del(key[1:], original)
		if err != nil {
			return nil, err
		}
		t.Next[prefix] = c2
		if c2 != nil && zero(c2.Count) {
			t.Next[prefix] = nil
		}
	}
	if err := t.update(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Trie) set(key []byte, kv *KeyValue) (*Trie, error) {
	t = t.Copy()
	if len(key) == 0 {
		t.Count = big.NewInt(1)
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
	if err := t.update(); err != nil {
		return nil, err
	}
	return t, nil
}

// increments i by v
func inc(i *big.Int, v int) *big.Int {
	if i == nil {
		i = big.NewInt(0)
	}
	var x big.Int
	x.Add(i, big.NewInt(int64(v)))
	return &x
}

func cp(i *big.Int) *big.Int {
	var o big.Int
	o.Set(i)
	return &o
}

func sum(i, j *big.Int) *big.Int {
	var out big.Int
	out.Add(i, j)
	return &out
}

func zero(i *big.Int) bool {
	return i.Cmp(big.NewInt(0)) == 0
}
func one(i *big.Int) bool {
	return i.Cmp(big.NewInt(1)) == 0
}

func ToGviz(g gviz.Graph, file, title string) error {
	gv, err := gviz.Compile(g, title, nil)
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
	label := func(x string) string {
		if x == "" {
			x = "nil"
		}
		return fmt.Sprintf("%s (%d; %x)", x, t.Count, t.Merkle[:2])
	}
	if t.KeyValue == nil {
		n.shape = "ellipse"
		n.label = label(string(parent))
	} else {
		n.label = label(fmt.Sprintf("%s = %s",
			string(t.KeyValue.Key),
			string(t.KeyValue.Value),
		))
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
