package trie

import "fmt"

type StringDatabase interface {
	Set(string, string)
	Get(string) (string, bool)
	Delete(string)
	Stats() *Stats
	Do(func(kv *StringKeyValue))
	// hash is implementation-dependent, but based only on key/values in the db
	Hash() []byte
}

type stringDB struct {
	hashLen int
	x       Database
}

func NewStrings(db Database) StringDatabase {
	return stringDB{hashLen: 8, x: db}
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
func (db stringDB) Hash() []byte {
	h := db.x.Hash()
	if n := db.hashLen; n > 0 {
		h = h[:n]
	}
	return h
}

func (db stringDB) Stats() *Stats {
	return db.x.Stats()
}

func (db stringDB) String() string {
	return fmt.Sprintf("%v", db.x)
}
