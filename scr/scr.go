// Package scr is for transaction scripting language
package scr

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/xoba/turd/cnfg"
	"github.com/xoba/turd/thash"
)

type Test struct {
	A string    `asn1:"printable"`
	B string    `asn1:"utf8"`
	X time.Time `asn1:"utc"`
	Y time.Time `asn1:"generalized"`
}

type Transaction struct {
	Inputs  []Input
	Outputs []Output
}

func (t Transaction) Signature(key *ecdsa.PrivateKey) ([]byte, error) {
	panic("")
}

type Output struct {
	Tokens *big.Int
	Script Steps
}

type Input struct {
	Transaction []byte
	Index       int
	Script      Steps
}

func Run(cnfg.Config) error {

	script := Script{
		RawTransaction: []byte("this is a test"),
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	sig, err := CreateSignature(script.RawTransaction, key)
	if err != nil {
		return err
	}
	pub, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}

	first := func() (out Steps) {
		add := func(s Step) {
			out = append(out, s)
		}
		add(Op("push"))
		add(Data(sig)) // sig
		add(Op("push"))
		add(Data(pub)) // pubkey
		return
	}()

	second := func() (out Steps) {
		add := func(s Step) {
			out = append(out, s)
		}
		add(Op("dup"))
		add(Op("hash"))
		add(Op("push"))
		add(Data(thash.Hash(pub)))
		add(Op("equal"))
		add(Op("verify"))
		add(Op("checksig"))
		return
	}()

	script.Steps = Add(first, second)
	fmt.Println(script)

	var s stack
	if err := script.Run(&s); err != nil {
		return err
	}
	fmt.Printf("final stack = %v\n", s)
	if b := s.pop(); !b.Boolean() {
		return fmt.Errorf("bad result: %x", b)
	}
	return nil
}

func testAsn1() error {
	var list []interface{}
	add := func(x interface{}) {
		list = append(list, x)
	}
	adde := func(x interface{}, e error) {
		if e != nil {
			panic(e)
		}
		list = append(list, x)
	}
	now := time.Now()
	add(Test{
		A: "abc",
		B: "abc",
		X: now,
		Y: now,
	})
	add(time.Second)
	add("ab*c")
	add(now)                                // utctime
	add(now.Add(50 * 365 * 24 * time.Hour)) // force generalized time
	add(now.UTC())
	add(523)
	a := big.NewInt(9223372036854775807)
	add(a)
	adde(parseOID("1.3.6.1.4.1.11129")) // google
	adde(parseOID("1.3.6.1.4.1.18300")) // xoba
	add(asn1.Enumerated(5))
	add(5)
	add([]byte("abc"))
	add(asn1.BitString{
		Bytes:     []byte("abc"),
		BitLength: 3,
	})
	buf, err := asn1.MarshalWithParams(list, "")
	if err != nil {
		return err
	}
	fmt.Println(base64.StdEncoding.EncodeToString(buf))
	for {
		var i interface{}
		rest, err := asn1.Unmarshal(buf, &i)
		if err != nil {
			return err
		}
		fmt.Printf("%d left; got: %T %v\n", len(rest), i, i)
		if len(rest) == 0 {
			break
		}
		buf = rest
	}
	return nil
}

func parseOID(id string) (asn1.ObjectIdentifier, error) {
	var out asn1.ObjectIdentifier
	for _, x := range strings.Split(id, ".") {
		v, err := strconv.ParseUint(x, 10, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, int(v))
	}
	return out, nil
}

type Script struct {
	PC             int // program counter
	RawTransaction []byte
	Steps          Steps
}

// checks that the signature is valid for given data
func CheckSignature(data, sig, pubkey []byte) error {
	pub, err := x509.ParsePKIXPublicKey(pubkey)
	if err != nil {
		return err
	}
	digest := thash.Hash(data)
	if !ecdsa.VerifyASN1(pub.(*ecdsa.PublicKey), digest, sig) {
		return fmt.Errorf("signature invalid")
	}
	return nil
}

func CreateSignature(data []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	digest := thash.Hash(data)
	sig, err := ecdsa.SignASN1(rand.Reader, key, digest)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func testSigs() error {
	data := []byte("this is a test")
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	sig, err := CreateSignature(data, key)
	if err != nil {
		return err
	}
	pub, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}
	if err := CheckSignature(data, sig, pub); err != nil {
		return err
	}
	return nil
}

func (s *Script) Next() *Step {
	if s.PC >= len(s.Steps) {
		return nil
	}
	x := s.Steps[s.PC]
	s.PC++
	return &x
}

type Steps []Step

func Add(a, b Steps) (out Steps) {
	out = append(out, a...)
	out = append(out, b...)
	return
}

type Step struct {
	Opcode string `json:",omitempty",asn1:"utf8,omitempty,optional"`
	Data   blob   `json:",omitempty",asn1:",omitempty,optional"`
}

var (
	OP_PUSH     = Op("push")     // the following data in script is pushed to stack
	OP_DUP      = Op("dup")      // duplicate top of stack
	OP_HASH     = Op("hash160")  // pop stack, hash and push
	OP_EQUAL    = Op("equal")    // pop two items of stack, push boolean value of equality
	OP_VERIFY   = Op("verify")   // pop, and if not true, exit with error
	OP_CHECKSIG = Op("checksig") // pop sig and pubkey and check signature is valid for the transaction, otherwise exit with error
)

func Op(name string) Step {
	return Step{Opcode: name}
}
func Data(d []byte) Step {
	return Step{Data: d}
}

func (s Step) String() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}
func (s Script) String() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func New() *Script {
	panic("")
}

func Parse() (*Script, error) {
	panic("")
}

func (script *Script) Run(s *stack) error {
	wrap := func(e error) error {
		return fmt.Errorf("at pc=%d: %w", script.PC-1, e)
	}
	for {
		x := script.Next()
		if x == nil {
			break
		}
		switch x.Opcode {
		case "dup":
			d := s.pop()
			s.push(d)
			s.push(d)
		case "hash":
			d := s.pop()
			s.push(thash.Hash(d))
		case "push":
			d := script.Next()
			s.push(blob(d.Data))
		case "equal":
			a, b := s.pop(), s.pop()
			if bytes.Equal(a, b) {
				s.push([]byte{1})
			} else {
				s.push([]byte{0})
			}
		case "verify":
			if !s.pop().Boolean() {
				return fmt.Errorf("unverified")
			}
		case "checksig":
			// TODO: the order seems wrong here:
			pub, sig := s.pop(), s.pop()
			if err := CheckSignature(script.RawTransaction, sig, pub); err != nil {
				s.push([]byte{0})
			} else {
				s.push([]byte{1})
			}
		default:
			return wrap(fmt.Errorf("illegal op: %q", x.Opcode))
		}
	}
	return nil
}

type blob []byte

func (b blob) Boolean() bool {
	var i big.Int
	i.SetBytes(b)
	return i.Cmp(big.NewInt(0)) != 0
}

type stack struct {
	Blobs []blob
}

func (s stack) String() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (s *stack) push(b blob) {
	s.Blobs = append(s.Blobs, b)
}

func (s *stack) len() int {
	return len(s.Blobs)
}

func (s *stack) pop() blob {
	n := s.len()
	if n == 0 {
		return nil
	}
	x := s.Blobs[n-1]
	s.Blobs = s.Blobs[:n-1]
	return x
}
