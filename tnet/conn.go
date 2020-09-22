package tnet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"fmt"
	"math/big"
	"net"
)

type conn struct {
	self   *PrivateKey
	other  *PublicKey
	c      net.Conn
	remote string
}

type packet struct {
	Payload []byte
	R, S    *big.Int
}

func (c conn) Remote() Node {
	return Node{Address: c.remote, PublicKey: c.other}
}

func (c conn) Receive() ([]byte, error) {
	buf, err := receive(c.c)
	if err != nil {
		return nil, err
	}
	var p packet
	if _, err := asn1.Unmarshal(buf, &p); err != nil {
		return nil, err
	}
	if !ecdsa.Verify(c.other.k, sha256d(p.Payload), p.R, p.S) {
		return nil, fmt.Errorf("can't authenticate packet")
	}
	return p.Payload, nil
}

func (c conn) Send(buf []byte) (err error) {
	p := packet{
		Payload: buf,
	}
	p.R, p.S, err = ecdsa.Sign(rand.Reader, c.self.k, sha256d(p.Payload))
	if err != nil {
		return err
	}
	m, err := asn1.Marshal(p)
	if err != nil {
		return err
	}
	return send(c.c, m)
}

func (c conn) Close() error {
	return c.c.Close()
}
