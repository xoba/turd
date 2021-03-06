package trans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"path"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/xoba/turd/lisp"
	"github.com/xoba/turd/trie"
)

func Trie() error {

	s, err := NewStorage()
	if err != nil {
		return err
	}

	r := rand.New(rand.NewSource(0))

	if err := s.IncBalance([]byte("abc"), big.NewInt(3)); err != nil {
		return err
	}

	var keys [][]byte

	for i := 0; i < 10; i++ {
		buf := make([]byte, 7)
		r.Read(buf)
		keys = append(keys, buf)
	}

	const n = 10000
	start := time.Now()
	for x := 0; x < n; x++ {
		key := keys[r.Intn(len(keys))]
		i := big.NewInt(1)
		if err := s.IncBalance(key, i); err != nil {
			return err
		}
	}
	fmt.Printf("%v/op\n", time.Since(start)/time.Duration(n))

	if err := s.db.ToGviz("trie.svg", "state"); err != nil {
		return err
	}

	return open.Run("trie.svg")
}

func Inflate(b *Balance, height *big.Int) (*Balance, error) {
	if height.Cmp(b.Height) < 0 {
		return nil, fmt.Errorf("negative inflation is illegal")
	}
	return &Balance{
		Owner:    b.Owner,
		Height:   height,
		Quantity: b.Quantity, // TODO: no inflation for now
	}, nil
}

type Storage struct {
	db trie.Database
}

// TODO: obviously having 32-deep trie's all the time is not space-efficient
func NewStorage() (*Storage, error) {
	db, err := trie.New()
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

// unique to the particular state of storage
func (s *Storage) Hash() ([]byte, error) {
	return s.db.Hash()
}

func (s *Storage) Copy() *Storage {
	return &Storage{db: s.db}
}

func (s *Storage) Reset(copy *Storage) {
	s.db = copy.db
}

// balances are at path prefix "/b/"
func balanceKey(key []byte) []byte {
	w := new(bytes.Buffer)
	w.WriteString("/")
	w.WriteRune('b')
	w.WriteString("/")
	w.Write(key)
	return w.Bytes()
}

// content are at path prefix "/c/"
func contentKey(p string) []byte {
	return []byte(path.Clean("/c/" + p))
}

type Content struct {
	Path        string
	Payload     []byte `json:",omitempty"` // generally should be empty, stored somewhere else.
	Length      *big.Int
	Hash        []byte `json:",omitempty"` // hash of the actual content, could also be a merkle root
	ContentType string `json:",omitempty"` // like a mime header
	Script      string // script to authorize content modification, if it already exists
}

func (i Content) Lisp() lisp.Exp {
	var self LispList
	add := func(name string, e Lisper) {
		self = append(self, []lisp.Exp{
			name,
			e.Lisp(),
		})
	}
	add("path", LispAtom(i.Path))
	add("length", LispInt(*i.Length))
	add("type", LispAtom(i.ContentType))
	add("payload", LispBlob(i.Payload))
	add("hash", LispBlob(i.Hash))
	return self.Lisp()
}

func (s *Storage) SetContent(c Content) error {
	var found Content
	if p := path.Clean(c.Path); p != c.Path {
		return fmt.Errorf("path not clean: %q", c.Path)
	}
	buf, err := s.db.Get(contentKey(c.Path))
	switch {
	case err == nil:
		if err := json.Unmarshal(buf, &found); err != nil {
			return err
		}
	case err == trie.NotFound:
		found = c
	default:
		return err
	}
	if buf2, err := json.Marshal(found); err != nil {
		return err
	} else {
		cp, err := s.db.Set(contentKey(c.Path), buf2)
		if err != nil {
			return err
		}
		s.db = cp
	}
	return nil
}

func (s *Storage) IncBalance(address []byte, byAmount *big.Int) error {
	address = balanceKey(address)
	var balance Balance
	buf, err := s.db.Get(address)
	switch {
	case err == nil:
		if err := json.Unmarshal(buf, &balance); err != nil {
			return err
		}
	case err == trie.NotFound:
		balance = Balance{Quantity: big.NewInt(0)}
	default:
		return err
	}
	balance.Quantity.Add(balance.Quantity, byAmount)
	buf, err = json.Marshal(balance)
	if err != nil {
		return err
	}
	db, err := s.db.Set(address, buf)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) GetBalance(address []byte) (*big.Int, error) {
	address = balanceKey(address)
	buf, err := s.db.Get(address)
	switch {
	case err == trie.NotFound:
		return big.NewInt(0), nil
	case err != nil:
		return nil, err
	default:
		var b Balance
		if err := json.Unmarshal(buf, &b); err != nil {
			return nil, err
		}
		return b.Quantity, nil
	}
}

// to be serialized in trie node corresponding to the address containing a balance
type Balance struct {
	Owner    []byte   // address of owner of this balance
	Height   *big.Int // blockchain height this balance was established
	Quantity *big.Int `json:"q,omitempty"`
}
