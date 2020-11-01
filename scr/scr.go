// Package scr is for transaction scripting language
package scr

import (
	"encoding/asn1"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/xoba/turd/cnfg"
)

type Test struct {
	A string    `asn1:"printable"`
	B string    `asn1:"utf8"`
	X time.Time `asn1:"utc"`
	Y time.Time `asn1:"generalized"`
}

func Run(cnfg.Config) error {
	return testAsn1()
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
}

func New() *Script {
	panic("")
}

func Parse() (*Script, error) {
	panic("")
}

func (s Script) Run() error {
	panic("")
}
