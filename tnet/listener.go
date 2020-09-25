package tnet

import (
	"fmt"
)

type listener struct {
	a Acceptor
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

	// send our key and nonce
	self, err := NewKeyAndNonce(key.Public())
	if err != nil {
		return nil, err
	}
	if err := self.send(c0); err != nil {
		return nil, err
	}

	// receive other's key and nonce
	other, err := receiveKeyAndNonce(c0)
	if err != nil {
		return nil, err
	}

	selfKey, err := self.GenerateSharedKey(key)
	if err != nil {
		return nil, err
	}
	otherKey, err := other.GenerateSharedKey(key)
	if err != nil {
		return nil, err
	}

	cn, err := newConn(c0, selfKey, otherKey)
	if err != nil {
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
	// receive other's address
	remote, err := cn.Receive()
	if err != nil {
		return nil, err
	}
	cn.remote = Node{Address: string(remote), PublicKey: other.Key}
	return cn, nil
}

func (ln listener) Close() error {
	return ln.a.Close()
}
