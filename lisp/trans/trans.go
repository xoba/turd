package trans

import (
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/lisp"
	"github.com/xoba/turd/thash"
	"github.com/xoba/turd/tnet"
)

// TODO: needs potentially multiple signatures
type Transaction struct {
	Type      string    `asn1:"optional,utf8" json:",omitempty"`
	Inputs    []Input   `asn1:"omitempty" json:",omitempty"`
	Outputs   []Output  `asn1:"omitempty" json:",omitempty"`
	Content   []Content `asn1:"omitempty" json:",omitempty"`
	Signature Signature `asn1:"omitempty" json:",omitempty"`
}

type Signature []byte

type Content struct {
	Key   []byte `asn1:"omitempty" json:",omitempty"`
	Value []byte `asn1:"omitempty" json:",omitempty"`
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

func marshal(buf []byte) string {
	return base64.RawStdEncoding.EncodeToString(buf)
}

func unmarshal(s string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(s)
}

func Run(cnfg.Config) error {
	key, err := tnet.NewKey()
	if err != nil {
		return err
	}
	pub, err := key.Public().MarshalBinary()
	if err != nil {
		return err
	}
	script := fmt.Sprintf(
		// outputs hash of nonce and verification:
		`
(lambda (hash signatures nonce) 
    (list (hash nonce) (verify '%s hash (car signatures))))`,
		marshal(pub),
	)

	outputs := make(map[string]Output)

	in := func(t *Transaction, i int64, s string) {
		t.Inputs = append(t.Inputs, Input{
			Quantity: big.NewInt(i),
			Script:   s,
		})
	}
	out := func(t *Transaction, i int64, addr []byte) {
		o := Output{
			Quantity: big.NewInt(i),
			Address:  addr,
		}
		t.Outputs = append(t.Outputs, o)
		outputs[marshal(o.Address)] = o
	}
	address := thash.Hash([]byte(script))

	var t1 Transaction
	{
		t1.Type = "coinbase"
		if false {
			in(&t1, 10, "")
		}
		out(&t1, 10, address)
		if err := t1.Sign(key); err != nil {
			return err
		}
		if err := t1.Validate(); err != nil {
			return err
		}
		if err := t1.Verify(key.Public()); err != nil {
			return err
		}
		fmt.Println(t1)
	}

	var t2 Transaction
	{
		in(&t2, 1, script)
		if err := t2.Sign(key); err != nil {
			return err
		}
		if err := t2.Validate(); err != nil {
			return err
		}
		if err := t2.Verify(key.Public()); err != nil {
			return err
		}
		fmt.Printf("script hash = %s\n", marshal(thash.Hash([]byte(t2.Inputs[0].Script))))
		fmt.Println(t2)
		for _, i := range t2.Inputs {
			addr := marshal(thash.Hash([]byte(t2.Inputs[0].Script)))
			o, ok := outputs[addr]
			if !ok {
				return fmt.Errorf("no such address: %s", addr)
			}
			fmt.Printf("processing %s -> %s\n", o, i)
			hash, err := t2.Hash()
			if err != nil {
				return err
			}
			e, err := lisp.Parse(
				fmt.Sprintf("(%s '%s '(%s) '%s)",
					i.Script,
					marshal(hash),
					marshal(t2.Signature),
					marshal([]byte(fmt.Sprintf("test %d", rand.Intn(3)))),
				),
			)
			if err != nil {
				return err
			}
			fmt.Println(lisp.String(e))
			res := lisp.Eval(e)
			fmt.Printf("res = %v\n", res)
		}
	}

	return nil
}

func test1() error {
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
	fmt.Println(marshal(buf))
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

func (t Transaction) String() string {
	buf, _ := json.Marshal(t)
	return string(buf)
}
func (t Input) String() string {
	buf, _ := json.Marshal(t)
	return string(buf)
}
func (t Output) String() string {
	buf, _ := json.Marshal(t)
	return string(buf)
}

func (t *Transaction) Validate() error {
	var coinbase bool
	switch x := t.Type; x {
	case "coinbase":
		coinbase = true
	case "":
	default:
		return fmt.Errorf("bad type: %q", x)
	}
	pos := func(i *big.Int) bool {
		return i.Cmp(big.NewInt(0)) > 0
	}
	neg := func(i *big.Int) bool {
		return i.Cmp(big.NewInt(0)) < 0
	}
	var i big.Int
	for _, x := range t.Inputs {
		if n := x.Quantity; !pos(n) {
			return fmt.Errorf("negative amount: %s", n)
		}
		i.Add(&i, x.Quantity)
	}
	for _, x := range t.Outputs {
		if n := x.Quantity; !pos(n) {
			return fmt.Errorf("negative amount: %s", n)
		}
		i.Sub(&i, x.Quantity)
	}
	if !coinbase && neg(&i) {
		return fmt.Errorf("negative transaction fee: %s", &i)
	}
	return nil
}

func (t *Transaction) Fee() *big.Int {
	var i big.Int
	for _, x := range t.Inputs {
		i.Add(&i, x.Quantity)
	}
	for _, x := range t.Outputs {
		i.Sub(&i, x.Quantity)
	}
	return &i
}

func (t Transaction) Hash() ([]byte, error) {
	t.Signature = nil
	m, err := t.Marshal()
	if err != nil {
		return nil, err
	}
	return thash.Hash(m), nil
}

func (t *Transaction) Sign(key *tnet.PrivateKey) error {
	h, err := t.Hash()
	if err != nil {
		return err
	}
	sig, err := key.Sign(h)
	if err != nil {
		return err
	}
	t.Signature = sig
	return nil
}

func (t Transaction) Verify(key *tnet.PublicKey) error {
	h, err := t.Hash()
	if err != nil {
		return err
	}
	return key.Verify(h, t.Signature)
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
