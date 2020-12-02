package trans

import (
	"encoding/json"
	"math/big"
	"math/rand"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"github.com/xoba/turd/trie"
)

func Trie() error {

	s, err := NewStorage()
	if err != nil {
		return err
	}

	r := rand.New(rand.NewSource(0))

	if err := s.UpdateBalance([]byte("abc"), big.NewInt(3)); err != nil {
		return err
	}

	keys := strings.Split("1,11,2,3,345,4,fhsdjkfhsfdk", ",")

	for x := 0; x < 1000; x++ {
		key := keys[r.Intn(len(keys))]
		i := big.NewInt(1)
		if err := s.UpdateBalance([]byte(key), i); err != nil {
			return err
		}
	}

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

func (s *Storage) UpdateBalance(address []byte, amount *big.Int) error {
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

// to be serialized in trie node corresponding to the address containing a balance
type Balance struct {
	Quantity *big.Int `json:"q,omitempty"`
}
