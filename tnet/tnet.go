// Package tnet is a network abstraction
package tnet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"net"
)

// all communication on network is authenticated
// by public key
type Network interface {
	Dial(to Node) (Conn, error)
	Listen() (Listener, error)
}

type Listener interface {
	Accept() (Conn, error)
	io.Closer
}

type Conn interface {
	Remote() Node // can be used in Network.Dial()
	Receive() ([]byte, error)
	Send([]byte) error
	io.Closer
}

type Node struct {
	Address   string     // how to reach the node
	PublicKey *PublicKey // node's public key
}

type conn struct {
	self  *PrivateKey
	other *PublicKey
	c     net.Conn
}

func (c conn) Remote() Node {
	panic("Remote unimplemented")
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

func receive(r io.Reader) ([]byte, error) {
	var n0 uint64
	if err := binary.Read(r, binary.BigEndian, &n0); err != nil {
		return nil, err
	}
	const max = 1000 * 1000
	if n0 > max {
		return nil, fmt.Errorf("can't handle buffers bigger than %d bytes", max)
	}
	buf := make([]byte, n0)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	if uint64(n) != n0 {
		return nil, fmt.Errorf("read %d/%d bytes", n, n0)
	}
	return buf, nil
}

func send(w io.Writer, buf []byte) error {
	var n0 uint64 = uint64(len(buf))
	if err := binary.Write(w, binary.BigEndian, n0); err != nil {
		return err
	}
	n, err := w.Write(buf)
	if err != nil {
		return err
	}
	if uint64(n) != n0 {
		return fmt.Errorf("wrote %d/%d bytes", n, n0)
	}
	return nil
}

func (c conn) Send(buf []byte) error {
	p := packet{
		Payload: buf,
	}
	r, s, err := ecdsa.Sign(rand.Reader, c.self.k, sha256d(p.Payload))
	if err != nil {
		return err
	}
	p.R = r
	p.S = s
	m, err := asn1.Marshal(p)
	if err != nil {
		return err
	}
	return send(c.c, m)
}

type packet struct {
	Payload []byte
	R, S    *big.Int
}

func sha256d(buf []byte) []byte {
	sha256 := func(x []byte) []byte {
		h := sha256.Sum256(x)
		return h[:]
	}
	return sha256(sha256(buf))
}

func (c conn) Close() error {
	return c.c.Close()
}

func (n network) Dial(to Node) (Conn, error) {
	c, err := net.Dial("tcp", to.Address)
	if err != nil {
		return nil, err
	}
	// send our public key
	buf1, err := n.key.Public().MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := send(c, buf1); err != nil {
		return nil, err
	}
	// receive other's public key
	buf2, err := receive(c)
	if err != nil {
		return nil, err
	}
	var other PublicKey
	if err := other.UnmarshalBinary(buf2); err != nil {
		return nil, err
	}
	cn := conn{
		self:  n.key,
		other: &other,
		c:     c,
	}
	return cn, nil
}

type listener struct {
	key *PrivateKey
	x   net.Listener
}

func (l listener) Accept() (Conn, error) {
	c, err := l.x.Accept()
	if err != nil {
		return nil, err
	}
	// receive other's public key
	buf2, err := receive(c)
	if err != nil {
		return nil, err
	}
	var other PublicKey
	if err := other.UnmarshalBinary(buf2); err != nil {
		return nil, err
	}
	// send our public key
	buf1, err := l.key.Public().MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := send(c, buf1); err != nil {
		return nil, err
	}

	cn := conn{
		self:  l.key,
		other: &other,
		c:     c,
	}
	return cn, nil
}

func (l listener) Close() error {
	return l.x.Close()
}

func (n network) Listen() (Listener, error) {
	x, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}
	return listener{x: x, key: n.key}, nil
}

type network struct {
	addr string
	key  *PrivateKey
}

func NewNetwork(pk *PrivateKey, port int) (Network, error) {
	if pk == nil {
		key, err := NewKey()
		if err != nil {
			return nil, err
		}
		pk = key
	}
	n := network{
		addr: fmt.Sprintf("localhost:%d", port),
		key:  pk,
	}
	return n, nil
}
