package lisp

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/xoba/turd/thash"
	"github.com/xoba/turd/tnet"
)

// valid types: string, *big.Int, []byte, time.Time, []Exp, Func, or error
type Exp interface{}

type Func func(...Exp) Exp

var (
	env   Exp
	Nil   Exp = list()
	True  Exp = "t"
	False Exp = Nil
)

func ExpToBool(e Exp) bool {
	switch t := e.(type) {
	case string:
		return t == "t"
	default:
		return false
	}
}

func BoolToExp(v bool) Exp {
	if v {
		return True
	}
	return False
}

func String(e Exp) string {
	switch t := e.(type) {
	case string:
		return t
	case []byte:
		return marshal(t)
	case *big.Int:
		return t.String()
	case time.Time:
		return t.Format(TimeFormat)
	case []Exp:
		if len(t) == 2 && t[0] == "quote" {
			return fmt.Sprintf("'%s", String(t[1]))
		}
		var list []string
		for _, e := range t {
			list = append(list, String(e))
		}
		return fmt.Sprintf("(%s)", strings.Join(list, " "))
	case Func:
		return String(t())
	case error:
		return fmt.Sprintf("error: %v", t)
	default:
		panic(fmt.Errorf(`[can't stringify type %T (%v)]`, t, t))
	}
}

func manifest(a Exp) Exp {
	if f, ok := a.(Func); ok {
		a = f()
	}
	return a
}

func one(args ...Exp) Exp {
	return manifest(args[0])
}

func two(args ...Exp) (Exp, Exp) {
	x, y := args[0], args[1]
	return manifest(x), manifest(y)
}

func three(args ...Exp) (Exp, Exp, Exp) {
	x, y, z := args[0], args[1], args[2]
	return manifest(x), manifest(y), manifest(z)
}

func checkargs(args []Exp) error {
	for _, a := range args {
		if e, ok := a.(error); ok {
			return e
		}
	}
	return nil
}

func apply(f Func, args ...Exp) Exp {
	if err := checkargs(args); err != nil {
		return err
	}
	return f(args...)
}

// ----------------------------------------------------------------------
// AXIOMS
// ----------------------------------------------------------------------

func quote(args ...Exp) Exp {
	return args[0]
}

func atom(args ...Exp) Exp {
	switch t := one(args...).(type) {
	case string, *big.Int, time.Time:
		return True
	case []Exp:
		return BoolToExp(len(t) == 0)
	default:
		return fmt.Errorf("illegal atom call: %T %v", t, t)
	}
}

func eq(args ...Exp) Exp {
	x, y := two(args...)
	switch x := x.(type) {
	case []Exp:
		switch y := y.(type) {
		case []Exp: // both lists:
			return BoolToExp(len(x) == 0 && len(y) == 0)
		default:
			return False
		}
	default:
		switch y := y.(type) {
		case []Exp:
			return False
		default: // both not lists:
			return BoolToExp(String(x) == String(y))
		}
	}
}

func car(args ...Exp) Exp {
	switch t := one(args...).(type) {
	case []Exp:
		switch len(t) {
		case 0:
			return Nil
		default:
			return t[0]
		}
	default:
		return fmt.Errorf("car needs list, got %T %v", t, t)
	}
}

func cdr(args ...Exp) Exp {
	switch t := one(args...).(type) {
	case []Exp:
		switch len(t) {
		case 0:
			return Nil
		default:
			return t[1:]
		}
	default:
		return fmt.Errorf("cdr needs list, got %T %v", t, t)
	}
}

func cons(args ...Exp) Exp {
	x := args[0]
	y := manifest(args[1])
	switch t := y.(type) {
	case []Exp:
		var out []Exp
		out = append(out, x)
		out = append(out, t...)
		return out
	default:
		return fmt.Errorf("cons needs a list")
	}
}

func cond(args ...Exp) Exp {
	if err := checkargs(args); err != nil {
		return err
	}
	for _, a := range args {
		switch t := a.(type) {
		case []Exp:
			if ExpToBool(manifest(t[0])) {
				return manifest(t[1])
			}
		default:
			return fmt.Errorf("illegal cond arg type %T: %v", t, t)
		}
	}
	return fmt.Errorf("cond fallthrough: %s", String(args))
}

func display(args ...Exp) Exp {
	a := one(args...)
	fmt.Printf("(display %s)\n", String(a))
	return a
}

func list(args ...Exp) Exp {
	return args
}

func compute(args []Exp, f func(...*big.Int) *big.Int) Exp {
	ints, err := toInts(args)
	if err != nil {
		return err
	}
	return f(ints...)
}

func toInts(args []Exp) ([]*big.Int, error) {
	var out []*big.Int
	add := func(i *big.Int) {
		out = append(out, i)
	}
	for _, e := range args {
		switch t := e.(type) {
		case string:
			var i big.Int
			if _, ok := i.SetString(t, 10); !ok {
				return nil, fmt.Errorf("can't parse %q", t)
			}
			add(&i)
		case *big.Int:
			add(t)
		default:
			return nil, fmt.Errorf("not an int: %T", e)
		}
	}
	return out, nil
}

func exp(args ...Exp) Exp {
	return compute(args, func(args ...*big.Int) *big.Int {
		return big.NewInt(0).Exp(args[0], args[1], args[2])
	})
}

func mul(args ...Exp) Exp {
	return compute(args, func(args ...*big.Int) *big.Int {
		return big.NewInt(0).Mul(args[0], args[1])
	})
}

func add(args ...Exp) Exp {
	return compute(args, func(args ...*big.Int) *big.Int {
		return big.NewInt(0).Add(args[0], args[1])
	})
}

func sub(args ...Exp) Exp {
	return compute(args, func(args ...*big.Int) *big.Int {
		return big.NewInt(0).Sub(args[0], args[1])
	})
}

func marshal(buf []byte) string {
	return base64.RawStdEncoding.EncodeToString(buf)
}

func unmarshal(e Exp) ([]byte, error) {
	switch t := e.(type) {
	case string:
		return base64.RawStdEncoding.DecodeString(t)
	case []byte:
		return t, nil
	default:
		return nil, fmt.Errorf("can't unmarshal %T", t)
	}
}

// hashes content
func hash(args ...Exp) Exp {
	buf, err := unmarshal(one(args...))
	if err != nil {
		return err
	}
	return thash.Hash(buf)
}

// concats two blobs
func concat(args ...Exp) Exp {
	x, y := two(args...)
	xb, err := unmarshal(x)
	if err != nil {
		return err
	}
	yb, err := unmarshal(y)
	if err != nil {
		return err
	}
	var out []byte
	out = append(out, xb...)
	out = append(out, yb...)
	return out
}

// creates a new private key
func newkey(args ...Exp) Exp {
	key, err := tnet.NewKey()
	if err != nil {
		return err
	}
	buf, err := key.MarshalBinary()
	if err != nil {
		return err
	}
	return buf
}

// derives public from private key
// (pub private) -> public
func pub(args ...Exp) Exp {
	buf, err := unmarshal(one(args...))
	if err != nil {
		return err
	}
	var private tnet.PrivateKey
	if err := private.UnmarshalBinary(buf); err != nil {
		return err
	}
	public, err := private.Public().MarshalBinary()
	if err != nil {
		return err
	}
	return public
}

// sign blob with private key
// (sign private blob) -> signature
func sign(args ...Exp) Exp {
	x, y := two(args...)
	var private tnet.PrivateKey
	{
		buf, err := unmarshal(x)
		if err != nil {
			return err
		}
		if err := private.UnmarshalBinary(buf); err != nil {
			return err
		}
	}
	blob, err := unmarshal(y)
	if err != nil {
		return err
	}
	sig, err := private.Sign(blob)
	if err != nil {
		return err
	}
	return sig
}

// verify blob with public key and signature
// (verify public blob signature) -> t or () [true or false]
func verify(args ...Exp) Exp {
	x, y, z := three(args...)
	var public tnet.PublicKey
	{
		buf, err := unmarshal(x)
		if err != nil {
			return err
		}
		if err := public.UnmarshalBinary(buf); err != nil {
			return err
		}
	}
	blob, err := unmarshal(y)
	if err != nil {
		return err
	}
	sig, err := unmarshal(z)
	if err != nil {
		return err
	}
	if err := public.Verify(blob, sig); err != nil {
		return False
	}
	return True
}

const TimeFormat = "2006-01-02T15:04:05.000Z"

func after(args ...Exp) Exp {
	x, y := two(args...)
	parse := func(e Exp) (*time.Time, error) {
		switch t := e.(type) {
		case time.Time:
			return &t, nil
		case string:
			tx, err := time.Parse(TimeFormat, t)
			if err != nil {
				return nil, err
			}
			return &tx, nil
		default:
			return nil, fmt.Errorf("not a time %T: %v", t, t)
		}
	}
	tx, err := parse(x)
	if err != nil {
		return err
	}
	ty, err := parse(y)
	if err != nil {
		return err
	}
	if tx.After(*ty) {
		return True
	}
	return False
}

// in development:

func err(args ...Exp) Exp {
	x := one(args...)
	return errors.New(String(x))
}

func runes(args ...Exp) Exp {
	x := one(args...)
	op, ok := x.(string)
	if !ok {
		return fmt.Errorf("not a string")
	}
	var out []Exp
	for _, r := range op {
		out = append(out, string(r))
	}
	return out
}

// return true if it's car, cdr, cadr, cddr, ..., caaar, etc., else false
func iscxr(args ...Exp) Exp {
	x := args[0]
	op, ok := x.(string)
	if !ok {
		return fmt.Errorf("not a string")
	}
	runes := []rune(op)
	if len(runes) < 4 {
		return False
	}
	if runes[0] != 'c' {
		return False
	}
	if runes[len(runes)-1] != 'r' {
		return False
	}
	for _, r := range runes[1 : len(runes)-1] {
		switch r {
		case 'a', 'd':
		default:
			return False
		}
	}
	return True
}

func cxr(args ...Exp) Exp {
	x, y := args[0], args[1]
	op, ok := x.(string)
	if !ok {
		return fmt.Errorf("not a string")
	}
	runes := []rune(op)
	n := len(runes)
	if len(runes) < 3 || runes[0] != 'c' || runes[n-1] != 'r' {
		return fmt.Errorf("not a cxr: %q", op)
	}
	e := y
	for i := 0; i < n-2; i++ {
		switch runes[n-i-2] {
		case 'a':
			e = car(e)
		case 'd':
			e = cdr(e)
		default:
			return fmt.Errorf("not a cxr: %q", op)
		}
	}
	return e
}

func funcall(args ...Exp) Exp {
	e := args[0]
	a := args[1]
	lambda := cadr(e)
	args2 := cddr(e)
	var out []Exp
	add := func(e Exp) {
		out = append(out, e)
	}
	add(eval(lambda, a))
	for _, e := range args2.([]Exp) {
		add(e)
	}
	return out
}
