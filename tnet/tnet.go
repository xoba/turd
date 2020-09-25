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
	c0, err := net.Dial("tcp", to.Address)
	if err != nil {
		return nil, err
	}

	// send our key and nonce
	self, err := NewKeyAndNonce(key.Public())
	if err != nil {
		return nil, err
	}
	if err := self.send(c0); err != nil {
		return nil, err
	}

	// send public key we're looking for:
	if pk := to.PublicKey; pk == nil {
		if err := send(c0, nil); err != nil {
			return nil, err
		}
	} else {
		if err := send(c0, pk.Hash()); err != nil {
			return nil, err
		}
	}

	// receive other's key and nonce
	other, err := receiveKeyAndNonce(c0)
	if err != nil {
		return nil, err
	}

	selfKey, err := GenerateSharedKey(self.Nonce, key, other.Key)
	if err != nil {
		return nil, err
	}
	otherKey, err := GenerateSharedKey(other.Nonce, key, other.Key)
	if err != nil {
		return nil, err
	}

	cn, err := newConn(c0, selfKey, otherKey)
	if err != nil {
		return nil, err
	}

	// send version
	if err := cn.Send([]byte(Version)); err != nil {
		return nil, err
	}
	// send our address
	if err := cn.Send([]byte(n.addr)); err != nil {
		return nil, err
	}
	cn.remote = Node{Address: "", PublicKey: other.Key}
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

// TODO: replace with scrypt
func hash(buf []byte) []byte {
	sha256 := func(x []byte) []byte {
		h := sha256.Sum256(x)
		return h[:]
	}
	return sha256(sha256(buf))
}
