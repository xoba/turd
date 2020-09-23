package tnet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"fmt"
	"math/big"
)

type conn struct {
	self           *PrivateKey
	other          *PublicKey
	c              RawConn
	remote         string
	sent, received *big.Int
}

func newConn(c RawConn, self *PrivateKey, other *PublicKey, remote string) *conn {
	return &conn{
		c:        c,
		self:     self,
		other:    other,
		remote:   remote,
		sent:     big.NewInt(0),
		received: big.NewInt(0),
	}
}

func inc(i *big.Int) *big.Int {
	var z big.Int
	one := big.NewInt(1)
	z.Add(i, one)
	return &z
}

type packet struct {
	Seq     *big.Int
	Payload []byte
	R, S    *big.Int
}

// hash of everything except signature R&S
func (p packet) hash() ([]byte, error) {
	w := new(bytes.Buffer)
	if err := send(w, p.Seq.Bytes()); err != nil {
		return nil, err
	}
	if err := send(w, p.Payload); err != nil {
		return nil, err
	}
	return sha256d(w.Bytes()), nil
}

func (c conn) Remote() Node {
	return Node{Address: c.remote, PublicKey: c.other}
}

func (c *conn) Receive() ([]byte, error) {
	buf, err := receive(c.c)
	if err != nil {
		return nil, err
	}
	var p packet
	if _, err := asn1.Unmarshal(buf, &p); err != nil {
		return nil, err
	}
	h, err := p.hash()
	if err != nil {
		return nil, err
	}
	if !ecdsa.Verify(c.other.k, h, p.R, p.S) {
		return nil, fmt.Errorf("can't authenticate packet")
	}
	if p.Seq.Cmp(c.received) != 0 {
		return nil, fmt.Errorf("sequence mismatch: got %s vs expected %s", p.Seq, c.received)
	}
	c.received = inc(c.received)
	return p.Payload, nil
}

func (c *conn) Send(buf []byte) (err error) {
	p := packet{
		Seq:     c.sent,
		Payload: buf,
	}
	c.sent = inc(c.sent)
	h, err := p.hash()
	if err != nil {
		return err
	}
	p.R, p.S, err = ecdsa.Sign(rand.Reader, c.self.k, h)
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
