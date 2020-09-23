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
	Dial(*PrivateKey, Node) (Conn, error)
	Listen() (Listener, error)
}

type Listener interface {
	// Accept a connection asking for any of the given keys
	Accept(...*PrivateKey) (Conn, error)
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

func NewTCPLocalhost(port int) (Network, error) {
	n := network{
		addr: fmt.Sprintf("localhost:%d", port),
	}
	return n, nil
}

type network struct {
	addr string
}

const Version = "1.0"

func (n network) Dial(key *PrivateKey, to Node) (Conn, error) {
	if key == nil {
		return nil, fmt.Errorf("needs key")
	}
	c, err := net.Dial("tcp", to.Address)
	if err != nil {
		return nil, err
	}
	keySender := func(send func([]byte) error) func(*PublicKey) error {
		return func(key *PublicKey) error {
			buf, err := key.MarshalBinary()
			if err != nil {
				return err
			}
			return send(buf)
		}
	}
	sendKey := keySender(func(buf []byte) error {
		return send(c, buf)
	})
	// send our own public key
	if err := sendKey(key.Public()); err != nil {
		return nil, err
	}
	cn := newConn(c, key, to.PublicKey, to.Address)
	// send version
	if err := cn.Send([]byte(Version)); err != nil {
		return nil, err
	}
	// send our address
	if err := cn.Send([]byte(n.addr)); err != nil {
		return nil, err
	}
	// send which key we expect
	if pk := to.PublicKey; pk == nil {
		if err := cn.Send([]byte("none")); err != nil {
			return nil, err
		}
	} else {
		if err := cn.Send([]byte("expect")); err != nil {
			return nil, err
		}
		sendKeySigned := keySender(func(buf []byte) error {
			return cn.Send(buf)
		})
		if err := sendKeySigned(to.PublicKey); err != nil {
			return nil, err
		}
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
	cn.other = &other
	return cn, nil
}

func (n network) Listen() (Listener, error) {
	ln, err := net.Listen("tcp", n.addr)
	if err != nil {
		return nil, err
	}
	return listener{a: acceptor{ln: ln}}, nil
}

func send(w io.Writer, buf []byte) error {
	var bufSize uint64 = uint64(len(buf))
	if bufSize > maxBufferSize {
		return fmt.Errorf("can't handle buffers bigger than %d bytes", maxBufferSize)
	}
	if err := binary.Write(w, binary.BigEndian, bufSize); err != nil {
		return err
	}
	n, err := w.Write(buf)
	if err != nil {
		return err
	}
	if uint64(n) != bufSize {
		return fmt.Errorf("wrote %d/%d bytes", n, bufSize)
	}
	return nil
}

const maxBufferSize = 1000 * 1000

func receive(r io.Reader) ([]byte, error) {
	var bufSize uint64
	if err := binary.Read(r, binary.BigEndian, &bufSize); err != nil {
		return nil, err
	}
	if bufSize > maxBufferSize {
		return nil, fmt.Errorf("can't handle buffers bigger than %d bytes", maxBufferSize)
	}
	buf := make([]byte, bufSize)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	if uint64(n) != bufSize {
		return nil, fmt.Errorf("read %d/%d bytes", n, bufSize)
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
