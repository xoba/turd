package trans

import (
	"bytes"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"sort"
	"text/template"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/lisp"
	"github.com/xoba/turd/thash"
	"github.com/xoba/turd/tnet"
)

// script output: either a blob or empty list.
// empty list means the validation failed, blob should
// be a mining hash, calculated in appropriate manner.

// TODO: needs potentially multiple signatures
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

func (t *Transaction) NewContent(key *tnet.PrivateKey, path string, content []byte) error {
	buf, err := key.Public().MarshalBinary()
	if err != nil {
		return err
	}
	c := Content{
		Path:    path,
		Owner:   buf,
		Payload: content,
		Length:  big.NewInt(int64(len(content))),
		Hash:    thash.Hash(content),
	}
	if err := c.Sign(key); err != nil {
		return err
	}
	t.Content = append(t.Content, c)
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
		Max:      big.NewInt(12),
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

func (b Block) ComputeHash() ([]byte, error) {
	b.Hash = nil
	b.FinalState = nil
	b.ID = nil
	buf, err := asn1.Marshal(b)
	if err != nil {
		return nil, err
	}
	return thash.Hash(buf), nil
}

// returns true if hash could be verified
func (b Block) VerifyHash() bool {
	h, err := b.ComputeHash()
	if err != nil {
		return false
	}
	return bytes.Equal(h, b.Hash)
}

/*

for processing transactions:

pass in block hash to first transaction script, then for each subsequent one:

if previous transaction output is false ("()"), that transaction is invalid.
otherwise, pass in that output as the input to next transaction.

last transaction's output is the block's output. if it satisfies the mining
requirement, then block is valid. otherwise, repeat this procedure.

*/

// quantity and script hash must match a previous transaction's output
// script is called with input, block, and transaction arguments,
// and has output. if output is false (nil or '()), transaction failed.
type Input struct {
	Quantity *big.Int
	Script   string   `asn1:"utf8"`
	Max      *big.Int // claimed max eval depth (like a "time")
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

type Block struct {
	Parent       []byte // parent in blockchain picture
	Height       *big.Int
	Time         time.Time
	Transactions []Transaction `json:",omitempty"`
	Threshold    *big.Int      // max hash value for this block to be valid mining
	Nonce        []byte        // a randomly chosen nonce for mining purposes
	Hash         []byte        // hash of entire block except ID field, this is the input to transactions
	Output       []byte        // output of all transactions
	FinalState   []byte        // pointer to the state trie after the transactions are processed
	ID           []byte        // hash of hash, output, and final state
}

func (b Block) String() string {
	buf, _ := json.MarshalIndent(b, "", "  ")
	return string(buf)
}

func Run(cnfg.Config) error {

	//return Trie()

	block := Block{
		Height:    big.NewInt(1000),
		Time:      time.Now().UTC(),
		Threshold: Difficulty(MaxHash(32), big.NewInt(30)),
		Nonce:     make([]byte, 20),
	}

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

	addt := func(t Transaction) {
		block.Transactions = append(block.Transactions, t)
	}

	after := block.Time.Add(-time.Millisecond)
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

	for i := 0; i < 10; i++ {
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
		if err := t.NewContent(key1, "/mykey", []byte("abc")); err != nil {
			return err
		}
		if err := t.Sign(key1, key3); err != nil {
			return err
		}
		addt(t)
	}

	fmt.Printf("block = %s\n", lisp.String(block.Lisp()))

	var allRounds []int

	times := make(map[string][]time.Duration)

	timing := func(name string, f func() error) error {
		t0 := time.Now()
		err := f()
		dt := time.Since(t0)
		times[name] = append(times[name], dt)
		return err
	}

	median := func(name string) time.Duration {
		list := times[name]
		sort.Slice(list, func(i, j int) bool {
			return list[i] < list[j]
		})
		return list[len(list)/2]
	}
	count := func(name string) int {
		return len(times[name])
	}

	start := time.Now()
	for time.Since(start) < time.Minute {
		if err := timing("round", func() error {
			var rounds int
			for {
				rounds++

				rand.Read(block.Nonce)

				state, err := NewStorage()
				if err != nil {
					return err
				}

				key := func(addr []byte) []byte {
					return []byte(fmt.Sprintf("%x", addr)[:4])
				}

				balance := func(addr []byte) *big.Int {
					i, err := state.GetBalance(key(addr))
					if err != nil {
						log.Fatal(err)
					}
					return i
				}

				inc := func(addr []byte, o *big.Int) {
					if err := state.IncBalance(key(addr), o); err != nil {
						log.Fatal(err)
					}
				}
				dec := func(addr []byte, o *big.Int) {
					inc(addr, big.NewInt(0).Neg(o))
				}

				type compiled struct {
					transaction lisp.Exp
					inputs      []lisp.Exp
					lengths     []*big.Int
				}

				var compiledTrans []compiled
				for _, t := range block.Transactions {
					if err := t.Validate(); err != nil {
						return err
					}
					c := compiled{
						transaction: t.Lisp(),
					}
					for _, input := range t.Inputs {
						x, err := lisp.Parse(input.Script)
						if err != nil {
							return err
						}
						c.inputs = append(c.inputs, x)
						if input.Max == nil || input.Max.Cmp(big.NewInt(0)) <= 0 {
							return fmt.Errorf("script max %v", input.Max)
						}
						c.lengths = append(c.lengths, input.Max)
					}
					compiledTrans = append(compiledTrans, c)
				}

				blockLisp := block.Lisp()
				quote := func(e lisp.Exp) lisp.Exp {
					return []lisp.Exp{"quote", e}
				}

				{
					bhash, err := block.ComputeHash()
					if err != nil {
						return err
					}
					block.Hash = bhash
				}

				concat := func(list ...[]byte) (out []byte) {
					for _, x := range list {
						out = append(out, x...)
					}
					return
				}

				// block hash to be chained through all inputs of all transactions
				if err := timing("proc", func() error {
					// TODO: output of each transaction is a "receipt"
					// for data it updates in trie.
					// also check that output of transaction is hashed!

					bhash := block.Hash

					for i, t := range block.Transactions {
						var output []byte
						if err := timing("trans", func() error {
							if err := t.Validate(); err != nil {
								return err
							}
							if err := timing("input", func() error {
								for j, input := range t.Inputs {
									if b := balance(input.Address()); b.Cmp(input.Quantity) < 0 {
										return fmt.Errorf("input %s from %s", input.Quantity, b)
									}
									e := []lisp.Exp{
										compiledTrans[i].inputs[j],
										quote(bhash),
										quote(blockLisp),
										quote(compiledTrans[i].transaction),
									}
									var res lisp.Exp
									if err := timing("eval", func() error {
										res = lisp.Try(e, compiledTrans[i].lengths[j])
										return nil
									}); err != nil {
										return err
									}
									var hashed bool
									switch t := res.(type) {
									case string:
										buf, err := base64.RawStdEncoding.DecodeString(t)
										if err != nil {
											return err
										}
										output = buf
									case []byte:
										output = t
									case *lisp.Blob:
										hashed = t.Hashed
										output = t.Content
									default:
										return fmt.Errorf("%d. bad result: %s\n", i, lisp.String(res))
									}
									if !hashed {
										return fmt.Errorf("transaction output not hashed: %s", lisp.String(res))
									}
									dec(input.Address(), input.Quantity)
								}
								return nil
							}); err != nil {
								return err
							}
							for _, o := range t.Outputs {
								inc(o.Address, o.Quantity)
							}
							for _, c := range t.Content {
								if err := c.Verify(); err != nil {
									return err
								}
								if err := state.SetContent(c); err != nil {
									return err
								}
							}
							// thread the output back to input of next transaction
							h, err := state.Hash()
							if err != nil {
								return err
							}
							block.Output = thash.Hash(concat(output, h))
							block.FinalState = h
							bhash = block.Output
							return nil
						}); err != nil {
							return err
						}
					}
					return nil
				}); err != nil {
					return err
				}

				block.ID = thash.Hash(concat(block.Hash, block.Output, block.FinalState))

				x := big.NewInt(0).SetBytes(block.ID)
				if x.Cmp(block.Threshold) < 0 {
					if false {
						if err := state.db.ToGviz("trie.svg", "state"); err != nil {
							return err
						}
						return open.Run("trie.svg")
					}
					break
				}
			}
			allRounds = append(allRounds, rounds)
			sort.Ints(allRounds)
			fmt.Printf("***** %d rounds; %d median (%d)\n",
				rounds,
				allRounds[len(allRounds)/2],
				len(allRounds),
			)
			return nil
		}); err != nil {
			return err
		}
	}
	fmt.Printf("all rounds: %v\n", allRounds)
	for k := range times {
		fmt.Printf("median %s time (%d): %v\n", k, count(k), median(k))
	}

	return nil
}

func MaxHash(n int) *big.Int {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = 255
	}
	return big.NewInt(0).SetBytes(buf)
}

func Difficulty(numerator, denominator *big.Int) *big.Int {
	return big.NewInt(0).Div(numerator, denominator)
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

func (t Transaction) Hash() ([]byte, error) {
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
