package trie

import (
	"bytes"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/thash"
)

type Database interface {
	Set([]byte, []byte)
	Get([]byte) ([]byte, bool)
	Delete([]byte)
	Stats() *Stats
	Do(func(kv *KeyValue))
	// hash is implementation-dependent, but based only on what is gettable in the db
	Hash() []byte
}

func New() *Trie {
	return &Trie{dirty: true}
}

type Trie struct {
	kv     *KeyValue
	dirty  bool
	stats  *Stats
	merkle []byte
	next   [256]*Trie
}

func (t *Trie) IsDirty() bool {
	return t.dirty
}

func (t *Trie) IsClean() bool {
	return !t.IsDirty()
}

func (t *Trie) MarkDirty() {
	t.dirty = true
	t.merkle = nil
	t.stats = nil
}

func (t *Trie) MarkClean() {
	t.dirty = false
}

type Stats struct {
	Count, Size *big.Int
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

type Iterable interface {
	Do(func(kv *KeyValue))
}

func String(i Iterable) string {
	var list []*KeyValue
	i.Do(func(kv *KeyValue) {
		list = append(list, kv)
	})
	buf, _ := json.Marshal(list)
	return string(buf)
}

func TestPaths(cnfg.Config) error {
	db := NewStrings(New())
	add := func(p string) {
		db.Set(p, "nil")
	}
	add("/a")
	add("/a/x")
	add("/a/y")
	add("/a/z")
	add("/a/z/123")
	add("/b")
	fmt.Println(db)
	return nil
}

func Run(c cnfg.Config) error {
	if err := TestPaths(c); err != nil {
		return err
	}
	var db1, db2 Database
	db1 = make(mapdb)
	db2 = New()
	if err := CheckDelete("map", NewStrings(db1)); err != nil {
		return err
	}
	if err := CheckDelete("trie", NewStrings(db2)); err != nil {
		return err
	}
	if err := Run2(c); err != nil {
		return err
	}
	if err := Run3(c); err != nil {
		return err
	}
	return nil
}

func CheckDelete(name string, db StringDatabase) error {
	checkStats := func(count, size int64) error {
		if n := db.Stats().Count.Int64(); n != count {
			return fmt.Errorf("%s got count = %d, expected %d", name, n, count)
		}
		if n := db.Stats().Size.Int64(); n != size {
			return fmt.Errorf("%s got size = %d, expected %d", name, n, size)
		}
		return nil
	}
	db.Set("a", "xx")
	if err := checkStats(1, 3); err != nil {
		return err
	}
	db.Set("b", "yy")
	if err := checkStats(2, 6); err != nil {
		return err
	}
	first := db.Hash()
	fmt.Printf("%s hash = %x\n", name, first)
	db.Set("c", "zz")
	if err := checkStats(3, 9); err != nil {
		return err
	}
	fmt.Printf("%s hash = %x\n", name, db.Hash())
	db.Delete("c")
	if err := checkStats(2, 6); err != nil {
		return err
	}
	second := db.Hash()
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
	fmt.Printf("trie hash: 0x%x\n", t.Hash())
	fmt.Printf("%v per iteration\n", time.Since(start)/n)
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (t *Trie) clean() {
	if t.IsClean() {
		return
	}
	t.merkle = t.computeHash()
	t.stats = t.computeStats()
	t.MarkClean()
}

func (t *Trie) computeStats() *Stats {
	var s Stats
	if t.kv != nil {
		s.Count = inc(s.Count, 1)
		s.Size = inc(s.Size, len(t.kv.Key))
		s.Size = inc(s.Size, len(t.kv.Value))
	}
	for _, x := range t.next {
		if x == nil {
			continue
		}
		s.Inc(x.Stats())
	}
	return &s
}

func (t *Trie) computeHash() []byte {
	var list [][]byte
	add := func(x []byte) {
		list = append(list, x)
	}
	if t.kv == nil {
		add(nil)
	} else {
		add(t.kv.Key)
		add(t.kv.Value)
	}
	for i, x := range t.next {
		if x == nil {
			continue
		}
		add([]byte{byte(i)})
		add(x.Hash())
	}
	buf, err := asn1.Marshal(list)
	check(err)
	t.merkle = thash.Hash(buf)
	return t.merkle
}

func (t *Trie) Stats() *Stats {
	if t.IsDirty() {
		t.clean()
	}
	return t.stats
}

func (t *Trie) Hash() []byte {
	if t.IsDirty() {
		t.clean()
	}
	return t.merkle
}

type KeyValue struct {
	Key, Value []byte
}

type StringKeyValue struct {
	Key, Value string
}

func (kv KeyValue) xHash() ([]byte, error) {
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

// TODO: "Do" should take a range arguemnt, and handler should return bool
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
	t.MarkDirty()
	current := t
	for _, b := range key {
		if x := current.next[b]; x != nil {
			current = x
		} else {
			newNode := New()
			current.next[b] = newNode
			current = newNode
		}
		current.MarkDirty()
	}
	current.kv = &KeyValue{Key: key, Value: value}
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
		if x := current.next[b]; x != nil {
			current = x
		} else {
			return
		}
		current.MarkDirty()
	}
	current.kv = nil
	t.Prune()
}

func (t *Trie) ToGviz(file string) error {
	return fmt.Errorf("ToGviz unimplemented")
}

func (t *Trie) Prune() bool {
	if t.IsClean() {
		return false
	}
	var children int
	for i, c := range t.next {
		if c == nil {
			continue
		}
		children++
		if c.Prune() {
			t.next[i] = nil
		}
	}
	if t.kv != nil {
		return false
	}
	if children > 0 {
		return false
	}
	return true
}
