package trans

import (
	"bytes"
	"crypto/rand"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"text/template"
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

func (t *Transaction) NewOutput(n int64, key *tnet.PublicKey, nonce []byte) error {
	script, err := NewScript(key, nonce)
	if err != nil {
		return err
	}
	t.Outputs = append(t.Outputs, Output{
		Quantity: big.NewInt(n),
		Address:  thash.Hash([]byte(script)),
	})
	return nil
}

func NewScript(key *tnet.PublicKey, nonce []byte) (string, error) {
	pubname, err := KeyName(key)
	if err != nil {
		return "", err
	}
	pub, err := key.MarshalBinary()
	if err != nil {
		return "", err
	}
	return replace(`
(lambda
  (input thash args) ; input, transaction hash, other arguments
  ((lambda (sig)
     (cond
      ((verify '{{.pub}} thash sig)        ; if signature verified:
       (hash (concat '{{.nonce}} input)))  ; hash the input with nonce
      ('t ())))                            ; else: return "false"
  (assoc '{{.pubname}} args)))
`, map[string]string{
		"pub":     marshal(pub),
		"nonce":   marshal(nonce),
		"pubname": pubname,
	})
}

func (t *Transaction) NewInput(n int64, key *tnet.PublicKey, nonce []byte) error {
	script, err := NewScript(key, nonce)
	if err != nil {
		return err
	}
	t.Inputs = append(t.Inputs, Input{
		Quantity: big.NewInt(n),
		Script:   script,
	})
	return nil
}

func (i Input) Address() []byte {
	return thash.Hash([]byte(i.Script))
}

type Block struct {
	Height        *big.Int
	Time          time.Time
	Transactions  []Transaction `asn1:"omitempty" json:",omitempty"`
	State         Hash          // pointer to the state trie
	ParentOutputs []Hash        // first is intra-chain, others are inter-chain
	Threshold     *big.Int      // max hash value for this block to be valid mining
	// a randomly chosen nonce for mining purposes:
	Nonce []byte
	// hash of this block, including all above fields,
	// but not including Output field below:
	Hash   Hash
	Output Hash // output of transactions, a kind of "ID" for this block
}

/*

for processing transactions:

pass in block hash to first transaction script, then for each subsequent one:

if previous transaction output is false ("()"), that transaction is invalid.
otherwise, pass in that output as the input to next transaction.

last transaction's output is the block's output. if it satisfies the mining
requirement, then block is valid. otherwise, repeat this procedure.

*/

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

func replace(s string, m map[string]string) (string, error) {
	t := template.New("script")
	if _, err := t.Parse(s); err != nil {
		return "", err
	}
	w := new(bytes.Buffer)
	if err := t.Execute(w, m); err != nil {
		return "", err
	}
	return formatScript(w.String())
}

func Run(cnfg.Config) error {

	key1, err := tnet.NewKey()
	if err != nil {
		return err
	}
	key2, err := tnet.NewKey()
	if err != nil {
		return err
	}
	key3, err := tnet.NewKey()
	if err != nil {
		return err
	}

	var trans []Transaction
	addt := func(t Transaction) {
		trans = append(trans, t)
	}

	{
		var t Transaction
		t.Type = "mining"
		if err := t.NewOutput(10, key1.Public(), []byte{0}); err != nil {
			return err
		}
		// no need to sign mining transaction
		addt(t)
	}
	{
		var t Transaction
		t.Type = "mining"
		if err := t.NewOutput(10, key3.Public(), []byte{0}); err != nil {
			return err
		}
		// no need to sign mining transaction
		addt(t)
	}

	{
		var t Transaction
		if err := t.NewInput(3, key1.Public(), []byte{0}); err != nil {
			return err
		}
		if err := t.NewInput(4, key3.Public(), []byte{0}); err != nil {
			return err
		}
		if err := t.NewOutput(2, key2.Public(), []byte{1}); err != nil {
			return err
		}
		if err := t.Sign(key1, key3); err != nil {
			return err
		}
		addt(t)
	}

	outputs := make(map[string][]Output)

	// block hash to be chained through all inputs of all transactions
	bhash := make([]byte, 10)
	rand.Read(bhash)

	for i, t := range trans {
		fmt.Printf("%d. %s\n", i, t)
		if err := t.Validate(); err != nil {
			return err
		}

		hash, err := t.Hash()
		if err != nil {
			return err
		}

		for j, input := range t.Inputs {
			addr := marshal(thash.Hash([]byte(input.Script)))
			_, ok := outputs[addr]
			if !ok {
				return fmt.Errorf("no such address: %s", addr)
			}
			e, err := lisp.Parse(
				fmt.Sprintf("(%s '%s '%s '%s)",
					input.Script,
					marshal(bhash),
					marshal(hash),
					t.Arguments,
				),
			)
			if err != nil {
				return err
			}
			res := lisp.Eval(e)
			switch t := res.(type) {
			case string:
				buf, err := base64.RawStdEncoding.DecodeString(t)
				if err != nil {
					return err
				}
				fmt.Printf("%d.%d. %s\n", i, j, marshal(buf))
				bhash = buf
			default:
				return fmt.Errorf("bad result: %s\n", lisp.String(res))
			}
		}

		for _, o := range t.Outputs {
			key := marshal(o.Address)
			outputs[key] = append(outputs[key], o)
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
	var mining bool
	switch x := t.Type; x {
	case "mining":
		mining = true
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
	if !mining && neg(&i) {
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
