package lisp

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/xoba/turd/thash"
	"github.com/xoba/turd/tnet"
)

// valid types: string, []Exp, Func, or error
type Exp interface{}

type Func func(...Exp) Exp

var (
	env   Exp
	Nil   Exp = list()
	True  Exp = "t"
	False Exp = Nil
)

func String(e Exp) string {
	switch t := e.(type) {
	case string:
		return t
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
		panic(fmt.Errorf("can't stringify type %T %v", t, t))
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

// ----------------------------------------------------------------------
// AXIOMS
// ----------------------------------------------------------------------

// TODO: maybe natively handle type []byte, rather than base64-encoded strings?

// return true if it's car, cdr, cadr, cddr, ..., caaar, etc., else false
func iscxr(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	x := args[0]
	op, ok := x.(string)
	if !ok {
		return fmt.Errorf("not a string")
	}
	runes := []rune(op)
	if len(runes) < 3 {
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
	if err := checklen(2, args); err != nil {
		return err
	}
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

func quote(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	return args[0]
}

func atom(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	switch t := one(args...).(type) {
	case string:
		return True
	case []Exp:
		return boolToExp(len(t) == 0)
	default:
		return fmt.Errorf("illegal atom call: %T %v", t, t)
	}
}

func eq(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x, y := two(args...)
	switch x := x.(type) {
	case string:
		switch y := y.(type) {
		case string: // both atoms
			return boolToExp(x == y)
		case []Exp:
			return False
		default:
			return fmt.Errorf("bad second argument to eq: %T", y)
		}
	case []Exp:
		switch y := y.(type) {
		case string:
			return False
		case []Exp: // both lists
			return boolToExp(len(x) == 0 && len(y) == 0)
		default:
			return fmt.Errorf("bad second argument to eq: %T", y)
		}
	default:
		return fmt.Errorf("bad first argument to eq: %T", x)
	}
}

func car(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
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
	if err := checklen(1, args); err != nil {
		return err
	}
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
	if err := checklen(2, args); err != nil {
		return err
	}
	x := args[0]
	y := manifest(args[1])
	switch y.(type) {
	case []Exp:
		var out []Exp
		out = append(out, x)
		out = append(out, y.([]Exp)...)
		return out
	default:
		return fmt.Errorf("cons needs a list")
	}
}

func cond(args ...Exp) Exp {
	if err := checkargs(args); err != nil {
		return err
	}
	get := func(i interface{}) Exp {
		switch v := i.(type) {
		case Func:
			return v()
		default:
			return v
		}
	}
	for _, a := range args {
		switch t := a.(type) {
		case []Exp:
			if err := checklen(2, t); err != nil {
				return err
			}
			p, e := t[0], t[1]
			if expToBool(get(p)) {
				return get(e)
			}
		default:
			return fmt.Errorf("illegal cond arg type %T: %v", t, t)
		}
	}
	return fmt.Errorf("cond fallthrough")
}

func display(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	a := one(args...)
	fmt.Printf("(display %s)\n", String(a))
	return a
}

func list(args ...Exp) Exp {
	return args
}

func mult(args ...Exp) Exp {
	return arith("mult", args, func(i, j *big.Int) *big.Int {
		var z big.Int
		z.Mul(i, j)
		return &z
	})
}

func plus(args ...Exp) Exp {
	return arith("plus", args, func(i, j *big.Int) *big.Int {
		var z big.Int
		z.Add(i, j)
		return &z
	})
}

func minus(args ...Exp) Exp {
	return arith("minus", args, func(i, j *big.Int) *big.Int {
		var z big.Int
		z.Sub(i, j)
		return &z
	})
}

func arith(name string, args []Exp, f func(*big.Int, *big.Int) *big.Int) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
	x, y := two(args...)
	xs, ok := x.(string)
	if !ok {
		return fmt.Errorf("x not string")
	}
	ys, ok := y.(string)
	if !ok {
		return fmt.Errorf("y not string")
	}
	//fmt.Printf("%s(%s,%s)\n", name, xs, ys)
	var xv, yv big.Int
	set := func(i *big.Int, s string) error {
		if _, ok := i.SetString(s, 10); !ok {
			return fmt.Errorf("can't parse %q", s)
		}
		return nil
	}
	if err := set(&xv, xs); err != nil {
		return err
	}
	if err := set(&yv, ys); err != nil {
		return err
	}
	return f(&xv, &yv).String()
}

func marshal(buf []byte) Exp {
	return base64.RawStdEncoding.EncodeToString(buf)
}

func unmarshal(e Exp) ([]byte, error) {
	s, ok := e.(string)
	if !ok {
		return nil, fmt.Errorf("unmarshal needs string, got %s", String(e))
	}
	return base64.RawStdEncoding.DecodeString(s)
}

// hashes content
func hash(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
	buf, err := unmarshal(one(args...))
	if err != nil {
		return err
	}
	return marshal(thash.Hash(buf))
}

// concats two blobs
func concat(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
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
	return marshal(out)
}

// creates a new private key
func newkey(args ...Exp) Exp {
	if err := checklen(0, args); err != nil {
		return err
	}
	key, err := tnet.NewKey()
	if err != nil {
		return err
	}
	buf, err := key.MarshalBinary()
	if err != nil {
		return err
	}
	return marshal(buf)
}

// derives public from private key
// (pub private) -> public
func pub(args ...Exp) Exp {
	if err := checklen(1, args); err != nil {
		return err
	}
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
	return marshal(public)
}

// sign blob with private key
// (sign private blob) -> signature
func sign(args ...Exp) Exp {
	if err := checklen(2, args); err != nil {
		return err
	}
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
	return marshal(sig)
}

// verify blob with public key and signature
// (verify public blob signature) -> t or () [true or false]
func verify(args ...Exp) Exp {
	if err := checklen(3, args); err != nil {
		return err
	}
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
	if err := checklen(2, args); err != nil {
		return err
	}
	x, y := two(args...)
	parse := func(e Exp) (*time.Time, error) {
		s, ok := e.(string)
		if !ok {
			return nil, fmt.Errorf("time not a string")
		}
		tx, err := time.Parse(TimeFormat, s)
		if err != nil {
			return nil, err
		}
		return &tx, nil
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
