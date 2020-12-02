package trans

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/skratchdot/open-golang/open"
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

type Storage struct {
	db trie.Database
}

func NewStorage() (*Storage, error) {
	db, err := trie.New()
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Copy() *Storage {
	return &Storage{db: s.db}
}

func (s *Storage) Reset(copy *Storage) {
	s.db = copy.db
}

func (s *Storage) IncBalance(address []byte, amount *big.Int) error {
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
	balance.Quantity.Add(balance.Quantity, amount)
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
	Quantity *big.Int `json:"q,omitempty"`
}
