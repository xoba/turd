// Package tnet is a network abstraction
package tnet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
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
	Address   string           // how to reach the node
	PublicKey *ecdsa.PublicKey // true identity of the node
}

type conn struct {
	c net.Conn
}

func (c conn) Remote() Node {
	panic("Remote unimplemented")
}

func (c conn) Receive() ([]byte, error) {
	var n0 uint64
	if err := binary.Read(c.c, binary.BigEndian, &n0); err != nil {
		return nil, err
	}
	const max = 1000 * 1000
	if n0 > max {
		return nil, fmt.Errorf("can't handle buffers bigger than %d bytes", max)
	}
	buf := make([]byte, n0)
	n, err := c.c.Read(buf)
	if err != nil {
		return nil, err
	}
	if uint64(n) != n0 {
		return nil, fmt.Errorf("read %d/%d bytes", n, n0)
	}
	return buf, nil
}

func (c conn) Send(buf []byte) error {
	var n0 uint64 = uint64(len(buf))
	if err := binary.Write(c.c, binary.BigEndian, n0); err != nil {
		return err
	}
	n, err := c.c.Write(buf)
	if err != nil {
		return err
	}
	if uint64(n) != n0 {
		return fmt.Errorf("wrote %d/%d bytes", n, n0)
	}
	return nil
}

func (c conn) Close() error {
	return c.c.Close()
}

func (n network) Dial(to Node) (Conn, error) {
	c, err := net.Dial("tcp", to.Address)
	if err != nil {
		return nil, err
	}
	var cn conn
	cn.c = c
	return &cn, nil
}

type listener struct {
	x net.Listener
}

func (l listener) Accept() (Conn, error) {
	c, err := l.x.Accept()
	if err != nil {
		return nil, err
	}
	return conn{c: c}, nil
}

func (l listener) Close() error {
	return l.x.Close()
}

func (n network) Listen() (Listener, error) {
	x, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}
	return listener{x: x}, nil
}

type network struct {
	port int
	addr string
	key  *ecdsa.PrivateKey
}

func NewKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func NewNetwork(pk *ecdsa.PrivateKey, port int) (Network, error) {
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
