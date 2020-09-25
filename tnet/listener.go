package tnet

import (
	"bytes"
	"fmt"
)

type listener struct {
	addr string
	a    Acceptor
}

func (ln listener) Accept(keys ...*PrivateKey) (Conn, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("needs at least one key")
	}
	if len(keys) > 1 {
		return nil, fmt.Errorf("only one key supported")
	}
	key := keys[0]
	c0, err := ln.a.Accept()
	if err != nil {
		return nil, err
	}

	var cleanup func()
	{
		// if we exit with error, close the connection
		cleanup = func() {
			c0.Close()
		}
		defer func() {
			if cleanup == nil {
				return
			}
			cleanup()
		}()
	}

	// receive other's key and nonce
	other, err := receiveKeyAndNonce(c0)
	if err != nil {
		return nil, err
	}

	{
		buf, err := receive(c0)
		if err != nil {
			return nil, err
		}
		if len(buf) > 0 {
			if !bytes.Equal(buf, key.Public().Hash()) {
				return nil, fmt.Errorf("don't have key")
			}
		}
	}

	// send our key and nonce
	self, err := NewKeyAndNonce(key.Public())
	if err != nil {
		return nil, err
	}
	if err := self.send(c0); err != nil {
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

	// receive version
	version, err := cn.Receive()
	if err != nil {
		return nil, err
	}
	if string(version) != Version {
		return nil, fmt.Errorf("bad version %q", string(version))
	}

	// send address
	if err := cn.Send([]byte(ln.addr)); err != nil {
		return nil, err
	}

	// receive other's address
	remote, err := cn.Receive()
	if err != nil {
		return nil, err
	}
	cn.remote = Node{Address: string(remote), PublicKey: other.Key}
	cleanup = nil
	return cn, nil
}

func (ln listener) Close() error {
	return ln.a.Close()
}
