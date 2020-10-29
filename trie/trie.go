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

func Run(c cnfg.Config) error {
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
	ToGviz(string) error
}

type Searchable interface {
	// error indicates not found or worse, else keyvalue is non-nil
	Search(SearchFunc) (*KeyValue, error)
}

type SearchFunc func(kv *KeyValue) bool

type KeyValue struct {
	Key, Value []byte
}

func New() (*Trie, error) {
	var t Trie
	if err := t.update(); err != nil {
		return nil, err
	}
	return &t, nil
}

// Copy without merkle hash
func (t *Trie) Copy() *Trie {
	return &Trie{
		KeyValue: t.KeyValue,
		Next:     t.Next,
	}
}

func (t *Trie) Hash() ([]byte, error) {
	if err := t.update(); err != nil {
		return nil, err
	}
	return t.Merkle, nil
}

// recursively computes and sets count and merkles if unset
func (t *Trie) update() error {
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
	count := big.NewInt(0)
	if t.KeyValue != nil {
		count = inc(count, 1)
	}
	for i, x := range t.Next {
		if x == nil {
			continue
		}
		add(i)
		if err := x.update(); err != nil {
			return err
		}
		add(x.Merkle)
		count = sum(count, x.Count)
	}
	buf, err := asn1.Marshal(list)
	if err != nil {
		return err
	}
	t.Merkle = thash.Hash(buf)
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

func (t *Trie) Set(key, value []byte) (Database, error) {
	if len(key) == 0 {
		return nil, EmptyKey
	}
	return t.set(key, &KeyValue{Key: key, Value: value})
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
	del := func(key string) error {
		if x, err := s.Delete("/c/1/2/3/4"); err != nil {
			return err
		} else {
			s = x
		}
		steps = append(steps, s)
		return nil
	}
	log := func(e error) {
		if e != nil {
			log.Fatal(e)
		}
	}
	log(step("/a", "b"))
	log(step("/c", "d"))
	for i := 0; i < 10; i++ {
		log(step("/c/1/2/3/4", "long"))
	}
	log(step("/ab", "xyz"))
	log(step("/a/x", "c"))
	if c.Delete {
		log(del("/c/1/2/3/4"))
	}
	for i := 0; i < 5; i++ {
		log(step(fmt.Sprintf("/a/x/%d", i), fmt.Sprintf("howdy %d", i)))
	}
	viz := func(s StringDatabase, i int) error {
		file := fmt.Sprintf("trie_%d.svg", i)
		if err := s.ToGviz(file); err != nil {
			return err
		}
		return open.Run(file)
	}
	for i, db := range steps {
		log(show(db))
		if err := viz(db, i); err != nil {
			return err
		}
	}
	return nil
}

func (t *Trie) ToGviz(file string) error {
	return ToGviz(t, file)
}

func (t *Trie) Delete(key []byte) (Database, error) {
	if len(key) == 0 {
		return nil, EmptyKey
	}
	t, err := t.del(key)
	return t, err
}

func (t *Trie) del(key []byte) (*Trie, error) {
	t = t.Copy()
	if len(key) > 0 {
		prefix := key[0]
		c := t.Next[prefix]
		if c == nil {
			return nil, NotFound
		}
		c2, err := c.del(key[1:])
		if err != nil {
			return nil, err
		}
		t.Next[prefix] = c2
		if c2 != nil && zero(c2.Count) {
			// prune this branch
			t.Next[prefix] = nil
		}
		if err := t.update(); err != nil {
			return nil, err
		}
		return t, nil
	}
	return nil, nil
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

func sum(i, j *big.Int) *big.Int {
	var out big.Int
	out.Add(i, j)
	return &out
}

func zero(i *big.Int) bool {
	return i.Cmp(big.NewInt(0)) == 0
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
	label := func(x string) string {
		if x == "" {
			x = "nil"
		}
		return fmt.Sprintf("%s (%d; x%x)", x, t.Count, t.Merkle[:2])
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
