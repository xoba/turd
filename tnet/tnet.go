// Package tnet is a network abstraction
package tnet

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// all communication on network is authenticated by public key
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
	Receive() (*Hashed, error)
	Send([]byte) error
	io.Closer
}

type Node struct {
	Address   string     // how to reach the node
	PublicKey *PublicKey // node's public key
}

type Hashed struct {
	Payload []byte
	Hash    []byte
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

type network struct {
	addr string
	key  *PrivateKey
}

func (n network) Dial(to Node) (Conn, error) {
	c, err := net.Dial("tcp", to.Address)
	if err != nil {
		return nil, err
	}
	// send our address:
	if err := send(c, []byte(n.addr)); err != nil {
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
	if to.PublicKey != nil {
		if !other.Equal(to.PublicKey) {
			return nil, fmt.Errorf("wrong key")
		}
	}
	cn := conn{
		remote: to.Address,
		self:   n.key,
		other:  &other,
		c:      c,
	}
	return cn, nil
}

func (n network) Listen() (Listener, error) {
	x, err := net.Listen("tcp", n.addr)
	if err != nil {
		return nil, err
	}
	return listener{ln: x, key: n.key}, nil
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

func sha256d(buf []byte) []byte {
	sha256 := func(x []byte) []byte {
		h := sha256.Sum256(x)
		return h[:]
	}
	return sha256(sha256(buf))
}
