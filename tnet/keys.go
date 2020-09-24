package tnet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/xoba/turd/cnfg"
)

func NewKey() (*PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{k: key}, nil
}

type PrivateKey struct {
	k *ecdsa.PrivateKey
}

type PublicKey struct {
	k *ecdsa.PublicKey
}

func (p PrivateKey) Public() *PublicKey {
	return &PublicKey{k: &p.k.PublicKey}
}

func (p PublicKey) Equal(o *PublicKey) bool {
	return p.k.Equal(o.k)
}

func (p PublicKey) MarshalBinary() (data []byte, err error) {
	return x509.MarshalPKIXPublicKey(p.k)
}

func (p *PublicKey) UnmarshalBinary(data []byte) error {
	key, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return err
	}
	pk, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("wrong key type: %T", pk)
	}
	p.k = pk
	return nil
}

func (p PublicKey) String() string {
	buf, _ := p.MarshalBinary()
	return "public:" + base64.StdEncoding.EncodeToString(sha256d(buf))
}

func GenerateSharedKey(nonce []byte, self *PrivateKey, other *PublicKey) ([]byte, error) {
	a, _ := other.k.Curve.ScalarMult(other.k.X, other.k.Y, self.k.D.Bytes())
	w := new(bytes.Buffer)
	w.Write(nonce)
	w.Write(a.Bytes())
	return sha256d(w.Bytes()), nil
}

func SharedKey(cnfg.Config) error {
	key1, err := NewKey()
	if err != nil {
		return err
	}
	key2, err := NewKey()
	if err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		nonce := make([]byte, 100)
		rand.Read(nonce)
		s1, err := GenerateSharedKey(nonce, key1, key2.Public())
		if err != nil {
			return err
		}
		fmt.Printf("shared1 = %x\n", s1)
		s2, err := GenerateSharedKey(nonce, key2, key1.Public())
		if err != nil {
			return err
		}
		fmt.Printf("shared2 = %x\n", s2)
	}
	return nil
}
