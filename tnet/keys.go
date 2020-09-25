package tnet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"

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
	return "public:" + base64.StdEncoding.EncodeToString(p.Hash())
}

func (p PublicKey) Hash() []byte {
	buf, _ := p.MarshalBinary()
	return hash(buf)
}

func GenerateSharedKey(nonce []byte, self *PrivateKey, other *PublicKey) ([]byte, error) {
	a, _ := other.k.Curve.ScalarMult(other.k.X, other.k.Y, self.k.D.Bytes())
	w := new(bytes.Buffer)
	w.Write(nonce)
	w.Write(a.Bytes())
	return hash(w.Bytes()), nil
}

type KeyAndNonce struct {
	Key   *PublicKey
	Nonce []byte
}

func NewKeyAndNonce(pk *PublicKey) (*KeyAndNonce, error) {
	nonce, err := createNonce(100)
	if err != nil {
		return nil, err
	}
	return &KeyAndNonce{Key: pk, Nonce: nonce}, nil
}

func (kn *KeyAndNonce) send(w io.Writer) error {
	buf, err := kn.Key.MarshalBinary()
	if err != nil {
		return err
	}
	if err := send(w, buf); err != nil {
		return err
	}
	return send(w, kn.Nonce)
}

func receiveKeyAndNonce(r io.Reader) (*KeyAndNonce, error) {
	buf, err := receive(r)
	if err != nil {
		return nil, err
	}
	var pk PublicKey
	if err := pk.UnmarshalBinary(buf); err != nil {
		return nil, err
	}
	kn := KeyAndNonce{
		Key: &pk,
	}
	kn.Nonce, err = receive(r)
	if err != nil {
		return nil, err
	}
	return &kn, nil
}

func createNonce(n int) ([]byte, error) {
	nonce := make([]byte, n)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return nonce, nil
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
		kn1, err := NewKeyAndNonce(key1.Public())
		if err != nil {
			return err
		}
		s1, err := GenerateSharedKey(kn1.Nonce, key1, key2.Public())
		if err != nil {
			return err
		}
		fmt.Printf("shared1 = %x\n", s1)
		s2, err := GenerateSharedKey(kn1.Nonce, key2, key1.Public())
		if err != nil {
			return err
		}
		fmt.Printf("shared2 = %x\n", s2)
	}
	return nil
}
