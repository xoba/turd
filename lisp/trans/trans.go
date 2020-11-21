package trans

import (
	"crypto/rand"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/lisp"
	"github.com/xoba/turd/thash"
	"github.com/xoba/turd/tnet"
)

// script output: either a blob or empty list.
// empty list means the validation failed, blob should
// be a mining hash, calculated in appropriate manner.

// TODO: needs potentially multiple signatures
// TODO: maybe signature block should be generic arguments to scripts,
// to be accessed via "assoc". makes sense to use some sort of hash of key
// for signature names.
type Transaction struct {
	Type      string    `asn1:"optional,utf8" json:",omitempty"`
	Inputs    []Input   `asn1:"omitempty" json:",omitempty"`
	Outputs   []Output  `asn1:"omitempty" json:",omitempty"`
	Content   []Content `asn1:"omitempty" json:",omitempty"`
	Arguments string    `asn1:"omitempty" json:",omitempty"`
}

type Block struct {
	Height       *big.Int
	Time         time.Time
	Transactions []Transaction `asn1:"omitempty" json:",omitempty"`
	State        Hash          // pointer to the state trie
	Parents      []Hash        // first is intra-chain, others are inter-chain
	Threshold    *big.Int      // max hash value for this block to be valid mining
	Nonce        []byte
	Hash         Hash // hash of this block
}

type Hash []byte

// content, compatible with a trie's KeyValue
type Content struct {
	Key    []byte   `asn1:"omitempty" json:",omitempty"`
	Value  []byte   `asn1:"omitempty" json:",omitempty"`
	Hash   []byte   `asn1:"omitempty" json:",omitempty"` // hash of key and value
	Length *big.Int // length of the value
}

// quantity and script hash must match a previous transaction's output
// script is called with hash and named arguments
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

func formatScript(s string) (string, error) {
	exp, err := lisp.Parse(s)
	if err != nil {
		return "", err
	}
	return lisp.String(exp), nil
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
	pubname, err := KeyName(key.Public())
	if err != nil {
		return err
	}
	script, err := formatScript(fmt.Sprintf(
		`
(lambda
  (nonce thash args) ; block hash, transaction hash, other arguments
  ((lambda (sig)
     (cond
      ((verify '%s thash sig) ; if signature verified
       (hash nonce))          ; hash the nonce
      ('t ())))               ; else return "false"
  (assoc '%s args)))
`,
		marshal(pub), pubname,
	))
	if err != nil {
		return err
	}

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
		hash, err := t2.Hash()
		if err != nil {
			return err
		}
		fmt.Println(t2)
		nonce := make([]byte, 5)
		rand.Read(nonce)
		for _, i := range t2.Inputs {
			addr := marshal(thash.Hash([]byte(t2.Inputs[0].Script)))
			o, ok := outputs[addr]
			if !ok {
				return fmt.Errorf("no such address: %s", addr)
			}
			fmt.Printf("processing %s -> %s\n", o, i)
			e, err := lisp.Parse(
				fmt.Sprintf("(%s '%s '%s '%s)",
					i.Script,
					marshal(nonce),
					marshal(hash),
					t2.Arguments,
				),
			)
			if err != nil {
				return err
			}
			fmt.Println(lisp.String(e))
			res := lisp.Eval(e)
			fmt.Printf("res = %s\n", lisp.String(res))
		}
	}

	return nil
}

func KeyName(key *tnet.PublicKey) (string, error) {
	buf, err := key.MarshalBinary()
	if err != nil {
		return "", err
	}
	return marshal(thash.Hash(buf)), nil
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

func (t Transaction) Hash() (Hash, error) {
	t.Arguments = ""
	m, err := t.Marshal()
	if err != nil {
		return nil, err
	}
	return thash.Hash(m), nil
}

func (t *Transaction) Sign(keys ...*tnet.PrivateKey) error {
	h, err := t.Hash()
	if err != nil {
		return err
	}
	var signature []lisp.Exp
	pair := func(a string, b []byte) {
		signature = append(signature, []lisp.Exp{a, marshal(b)})
	}
	for _, key := range keys {
		sig, err := key.Sign(h)
		if err != nil {
			return err
		}
		name, err := KeyName(key.Public())
		if err != nil {
			return err
		}
		pair(name, sig)
	}
	t.Arguments = lisp.String(signature)
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
