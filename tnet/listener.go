package tnet

import (
	"bytes"
	"fmt"
)

type listener struct {
	addr string
	a    Acceptor
}

type cleaner struct {
	clean    bool
	deferred func()
}

func newCleaner(deferred func()) *cleaner {
	return &cleaner{
		deferred: deferred,
	}
}

func (c *cleaner) MarkClean() {
	c.clean = true
}

func (c *cleaner) Cleanup() {
	if c.clean {
		return
	}
	c.deferred()
}

func (ln listener) Accept(keys ...*PrivateKey) (Conn, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("needs at least one key")
	}
	if len(keys) > 1 {
		return nil, fmt.Errorf("only one key supported")
	}
	key := keys[0]
	insecureConn, err := ln.a.Accept()
	if err != nil {
		return nil, err
	}

	cleaner := newCleaner(func() {
		insecureConn.Close()
	})
	defer cleaner.Cleanup()

	// receive other's key and nonce
	other, err := receiveKeyAndNonce(insecureConn)
	if err != nil {
		return nil, err
	}

	// receive other's request for our key
	{
		buf, err := receive(insecureConn)
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
	if err := self.send(insecureConn); err != nil {
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

	secure, err := newConn(insecureConn, selfKey, otherKey)
	if err != nil {
		return nil, err
	}

	if err := secure.negotiate(ln.addr, other.Key); err != nil {
		return nil, err
	}

	cleaner.MarkClean()
	return secure, nil
}

func (ln listener) Close() error {
	return ln.a.Close()
}
