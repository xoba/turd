package trie

import (
	"fmt"
)

type StringDatabase interface {
	Set(string, string) (StringDatabase, error)
	Get(string) (string, error)
	Delete(string) (StringDatabase, error)
	Hash() ([]byte, error)
	Search(func(kv *StringKeyValue) bool) (*StringKeyValue, error)
	Visualizable
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

func (db stringDB) ToGviz(file string) error {
	return db.x.ToGviz(file)
}

func (db stringDB) Set(key string, value string) (StringDatabase, error) {
	x, err := db.x.Set([]byte(key), []byte(value))
	if err != nil {
		return nil, err
	}
	return NewStrings(x), nil
}

func (db stringDB) Get(key string) (string, error) {
	v, err := db.x.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (db stringDB) Delete(key string) (StringDatabase, error) {
	x, err := db.x.Delete([]byte(key))
	if err != nil {
		return nil, err
	}
	return NewStrings(x), nil
}

func (db stringDB) Search(f func(kv *StringKeyValue) bool) (*StringKeyValue, error) {
	v, err := db.x.Search(func(kv *KeyValue) bool {
		return f(&StringKeyValue{
			Key:   string(kv.Key),
			Value: string(kv.Value),
		})
	})
	if err != nil {
		return nil, err
	}
	return &StringKeyValue{Key: string(v.Key), Value: string(v.Value)}, nil
}

func (db stringDB) String() string {
	return fmt.Sprintf("%v", db.x)
}

func (db stringDB) Hash() ([]byte, error) {
	return db.x.Hash()
}
