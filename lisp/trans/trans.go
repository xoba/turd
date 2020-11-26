package trans

import (
	"bytes"
	"crypto/rand"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (t Transaction) Lisp() lisp.Exp {
	var self LispList
	add := func(list *LispList, name string, e Lisper) {
		*list = append(*list, []lisp.Exp{
			name,
			e.Lisp(),
		})
	}
	add(&self, "type", LispAtom(t.Type))
	h, err := t.Hash()
	if err != nil {
		return err
	}
	add(&self, "hash", LispBlob(h))
	{
		var list LispList
		for i, x := range t.Inputs {
			add(&list, fmt.Sprintf("%d", i), x)
		}
		add(&self, "inputs", list)
	}
	{
		var list LispList
		for i, x := range t.Outputs {
			add(&list, fmt.Sprintf("%d", i), x)
		}
		add(&self, "outputs", list)
	}
	{
		var list LispList
		for i, x := range t.Content {
			add(&list, fmt.Sprintf("%d", i), x)
		}
		add(&self, "content", list)
	}
	add(&self, "arguments", LispString(t.Arguments))
	return self.Lisp()
}

func (i Input) Lisp() lisp.Exp {
	var self LispList
	add := func(name string, e Lisper) {
		self = append(self, []lisp.Exp{
			name,
			e.Lisp(),
		})
	}
	add("quantity", LispInt(*i.Quantity))
	add("script", LispString(i.Script))
	return self.Lisp()
}

func (i Output) Lisp() lisp.Exp {
	var self LispList
	add := func(name string, e Lisper) {
		self = append(self, []lisp.Exp{
			name,
			e.Lisp(),
		})
	}
	add("quantity", LispInt(*i.Quantity))
	add("address", LispBlob(i.Address))
	return self.Lisp()
}

type LispTime time.Time

func (x LispTime) Lisp() lisp.Exp {
	return time.Time(x).Format(lisp.TimeFormat)
}

type LispBlob []byte

func (x LispBlob) Lisp() lisp.Exp {
	return marshal(x)
}

type LispInt big.Int

func (x LispInt) Lisp() lisp.Exp {
	i := big.Int(x)
	z := &i
	return z.String()
}

type LispList []lisp.Exp

func (x LispList) Lisp() lisp.Exp {
	return []lisp.Exp(x)
}

type LispExpr struct {
	lisp.Exp
}

func (x LispExpr) Lisp() lisp.Exp {
	return x.Exp
}

type LispAtom string

func (x LispAtom) Lisp() lisp.Exp {
	return string(x)
}

type LispString string

func (x LispString) Lisp() lisp.Exp {
	if x == "" {
		return []lisp.Exp{}
	}
	e, err := lisp.Parse(string(x))
	if err != nil {
		return err
	}
	return e
}

func (i Content) Lisp() lisp.Exp {
	var self LispList
	add := func(name string, e Lisper) {
		self = append(self, []lisp.Exp{
			name,
			e.Lisp(),
		})
	}
	add("key", LispBlob(i.Key))
	add("hash", LispBlob(i.Hash))
	add("value", LispBlob(i.Value))
	add("length", LispInt(*i.Length))
	return self.Lisp()
}

type Lisper interface {
	Lisp() lisp.Exp
}

func (t *Transaction) NewOutput(n int64, key *tnet.PublicKey, nonce string, after time.Time) error {
	script, err := NewScript(key, nonce, after)
	if err != nil {
		return err
	}
	t.Outputs = append(t.Outputs, Output{
		Quantity: big.NewInt(n),
		Address:  thash.Hash([]byte(script)),
	})
	return nil
}

// TODO: pass the block height and time in too!
func NewScript(key *tnet.PublicKey, nonce string, after time.Time) (string, error) {
	pubname, err := KeyName(key)
	if err != nil {
		return "", err
	}
	pub, err := key.MarshalBinary()
	if err != nil {
		return "", err
	}
	buf, err := ioutil.ReadFile("lisp/trans/script.lisp")
	if err != nil {
		return "", err
	}
	return replace(string(buf), map[string]string{
		"t0":      after.Format(lisp.TimeFormat),
		"pub":     marshal(pub),
		"nonce":   marshal([]byte(nonce)),
		"pubname": pubname,
	})
}

func (t *Transaction) NewInput(n int64, key *tnet.PublicKey, nonce string, after time.Time) error {
	script, err := NewScript(key, nonce, after)
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

func (b Block) Lisp() lisp.Exp {
	var self LispList
	add := func(list *LispList, name string, e Lisper) {
		*list = append(*list, []lisp.Exp{
			name,
			e.Lisp(),
		})
	}
	add(&self, "height", LispInt(*b.Height))
	add(&self, "time", LispTime(b.Time))
	//add(&self, "hash", LispBlob(b.Hash))
	return self.Lisp()
}

type Block struct {
	Height        *big.Int
	Time          time.Time
	Hash          Hash          // hash of entire block except ID field
	ID            Hash          // final output of all chained transaction scripts, a kind of "ID" for this block
	Transactions  []Transaction `asn1:"omitempty" json:",omitempty"`
	State         Hash          // pointer to the state trie
	ParentOutputs []Hash        // first is intra-chain, others are inter-chain
	Threshold     *big.Int      // max hash value for this block to be valid mining
	// a randomly chosen nonce for mining purposes:
	Nonce []byte
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
	Hash   []byte   `asn1:"omitempty" json:",omitempty"` // hash of key and value
	Value  []byte   `asn1:"omitempty" json:",omitempty"`
	Length *big.Int // length of the value
}

// quantity and script hash must match a previous transaction's output
// script is called with input, block, and transaction arguments,
// and has output. if output is false (nil or '()), transaction failed.
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

	now := time.Now().UTC()
	after := now.Add(-time.Millisecond)
	{
		var t Transaction
		t.Type = "turd"
		if err := t.NewOutput(1000, key1.Public(), "key1", after); err != nil {
			return err
		}
		// no need to sign mining transaction
		addt(t)
	}
	{
		var t Transaction
		t.Type = "turd"
		if err := t.NewOutput(1000, key3.Public(), "key3", after); err != nil {
			return err
		}
		// no need to sign mining transaction
		addt(t)
	}

	for i := 0; i < 100; i++ {
		var t Transaction
		t.Type = "normal"
		if err := t.NewInput(3, key1.Public(), "key1", after); err != nil {
			return err
		}
		if err := t.NewInput(4, key3.Public(), "key3", after); err != nil {
			return err
		}
		if err := t.NewOutput(7, key2.Public(), "key2", after); err != nil {
			return err
		}
		if err := t.Sign(key1, key3); err != nil {
			return err
		}
		addt(t)
	}

	balances := make(map[string]*big.Int)

	balance := func(addr []byte) *big.Int {
		key := marshal(addr)[:6]
		i, ok := balances[key]
		if !ok {
			i = big.NewInt(0)
			balances[key] = i
		}
		return i
	}
	inc := func(addr []byte, o *big.Int) {
		i := balance(addr)
		i.Add(i, o)
	}
	dec := func(addr []byte, o *big.Int) {
		inc(addr, big.NewInt(0).Neg(o))
	}

	block := Block{
		Height: big.NewInt(1000),
		Time:   time.Now().UTC(),
	}

	fmt.Printf("block = %s\n", lisp.String(block.Lisp()))

	// block hash to be chained through all inputs of all transactions
	bhash := make([]byte, 10)
	rand.Read(bhash)
	for i, t := range trans {
		fmt.Printf("%d. %s\n", i, t)

		fmt.Printf("EXPR = %s\n", lisp.String(t.Lisp()))

		if err := t.Validate(); err != nil {
			return err
		}

		for j, input := range t.Inputs {
			if b := balance(input.Address()); b.Cmp(input.Quantity) < 0 {
				return fmt.Errorf("input %s from %s", input.Quantity, b)
			}
			e, err := lisp.Parse(
				fmt.Sprintf("(%s '%s '%s '%s)",
					input.Script,
					marshal(bhash),
					lisp.String(block.Lisp()),
					lisp.String(t.Lisp()),
				),
			)
			if err != nil {
				return err
			}
			fmt.Printf("EVAL(%s)\n", lisp.String(e))
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
			dec(input.Address(), input.Quantity)
			fmt.Printf("balances: %s\n", balances)
		}

		for j, o := range t.Outputs {
			fmt.Printf("%d.%d. output %s\n", i, j, o)
			inc(o.Address, o.Quantity)
			fmt.Printf("balances: %s\n", balances)
		}
	}

	fmt.Printf("final hash = %s\n", marshal(bhash))
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
	case "turd":
		mining = true
	case "normal":
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
