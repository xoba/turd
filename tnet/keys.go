package tnet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
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
