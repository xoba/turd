package trie

import (
	"fmt"
)

type StringDatabase interface {
	Set(string, string) StringDatabase
	Get(string) (string, bool)
	Delete(string)
	Hash() []byte
	Search(func(kv *StringKeyValue) bool) *StringKeyValue
}

type StringKeyValue struct {
	Key, Value string
}

type stringDB struct {
	hashLen int
	x       Database
}

func NewStrings(db Database) StringDatabase {
	return stringDB{hashLen: 8, x: db}
}

func (db stringDB) Set(key string, value string) StringDatabase {
	return NewStrings(db.x.Set([]byte(key), []byte(value)))
}

func (db stringDB) Get(key string) (string, bool) {
	v, ok := db.x.Get([]byte(key))
	return string(v), ok
}

func (db stringDB) Delete(key string) {
	db.x.Delete([]byte(key))
}

func (db stringDB) Search(f func(kv *StringKeyValue) bool) *StringKeyValue {
	v := db.x.Search(func(kv *KeyValue) bool {
		return f(&StringKeyValue{
			Key:   string(kv.Key),
			Value: string(kv.Value),
		})
	})
	if v == nil {
		return nil
	}
	return &StringKeyValue{Key: string(v.Key), Value: string(v.Value)}
}

func (db stringDB) String() string {
	return fmt.Sprintf("%v", db.x)
}

func (db stringDB) Hash() []byte {
	return db.x.Hash()
}

func (db stringDB) Unwrap() Database {
	return db.x
}
