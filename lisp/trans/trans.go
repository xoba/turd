package trans

import (
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/tnet"
)

func Run(cnfg.Config) error {
	fmt.Println("playing with transactions using lisp")
	key, err := tnet.NewKey()
	if err != nil {
		return err
	}
	var t Transaction
	in := func(i int64, s string) {
		t.Inputs = append(t.Inputs, Input{
			Quantity: big.NewInt(i),
			Script:   s,
		})
	}
	out := func(i int64, addr []byte) {
		t.Outputs = append(t.Outputs, Output{
			Quantity: big.NewInt(i),
			Address:  addr,
		})
	}
	for i := 0; i < 5; i++ {
		in(5, fmt.Sprintf("testing %d", i))
	}
	for i := 0; i < 3; i++ {
		out(10, []byte(fmt.Sprintf("xyz %d", i)))
	}

	if err := t.Sign(key); err != nil {
		return err
	}
	buf, err := t.Marshal()
	if err != nil {
		return err
	}
	fmt.Println(base64.StdEncoding.EncodeToString(buf))

	t2, err := Unmarshal(buf)
	if err != nil {
		return err
	}

	if err := t2.Verify(key.Public()); err != nil {
		return err
	}

	fmt.Printf("verified %s\n", t2)
	return nil
}

func Coinbase(key *tnet.PrivateKey, address []byte, reward *big.Int) (*Transaction, error) {
	var t Transaction
	t.Type = "coinbase"
	t.Outputs = append(t.Outputs, Output{
		Quantity: reward,
		Address:  address,
	})
	if err := t.Sign(key); err != nil {
		return nil, err
	}
	return &t, nil
}

type Block struct {
	Transactions []Transaction `asn1:"omitempty"`
}

type Transaction struct {
	Type      string   `asn1:"optional,utf8"`
	Inputs    []Input  `asn1:"omitempty"`
	Outputs   []Output `asn1:"omitempty"`
	Signature []byte
}

// quantity and script hash must match a previous transaction's output
type Input struct {
	Quantity *big.Int
	Script   string `asn1:"utf8"`
}

type Output struct {
	Quantity *big.Int
	Address  []byte
}

func (t Transaction) String() string {
	buf, _ := json.Marshal(t)
	return string(buf)
}

func (t Transaction) Verify(key *tnet.PublicKey) error {
	sig := t.Signature
	t.Signature = nil
	m, err := t.Marshal()
	if err != nil {
		return err
	}
	return key.Verify(m, sig)
}

func (t *Transaction) Sign(key *tnet.PrivateKey) error {
	t.Signature = nil
	m, err := t.Marshal()
	if err != nil {
		return err
	}
	sig, err := key.Sign(m)
	if err != nil {
		return err
	}
	t.Signature = sig
	return nil
}

func Unmarshal(buf []byte) (*Transaction, error) {
	var t Transaction
	rest, err := asn1.Unmarshal(buf, &t)
	if err != nil {
		return nil, err
	}
	if n := len(rest); n > 0 {
		return nil, fmt.Errorf("%d bytes extraneous content", n)
	}
	return &t, nil
}

func (t Transaction) Marshal() ([]byte, error) {
	return asn1.Marshal(t)
}
