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
	c, err := ln.a.Accept()
	if err != nil {
		return nil, err
	}
	keyReceiver := func(receive func() ([]byte, error)) func() (*PublicKey, error) {
		return func() (*PublicKey, error) {
			buf, err := receive()
			if err != nil {
				return nil, err
			}
			var other PublicKey
			if err := other.UnmarshalBinary(buf); err != nil {
				return nil, err
			}
			return &other, nil
		}
	}
	receiveKey := keyReceiver(func() ([]byte, error) {
		return receive(c)
	})
	// receive other's public key
	other, err := receiveKey()
	if err != nil {
		return nil, err
	}
	cn := newConn(c, keys[0], other, "")
	receiveKeySigned := keyReceiver(func() ([]byte, error) {
		return cn.Receive()
	})
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
	cn.remote = string(remote)
	str, err := cn.Receive()
	switch x := string(str); x {
	case "none":
	case "expect":
		other, err := receiveKeySigned()
		if err != nil {
			return nil, err
		}
		if !other.Equal(key.Public()) {
			return nil, fmt.Errorf("we don't have key %s", other)
		}
	default:
		return nil, fmt.Errorf("illegal: %q", x)
	}
	// send our public key
	buf1, err := key.Public().MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := send(c, buf1); err != nil {
		return nil, err
	}
	return cn, nil
}

func (ln listener) Close() error {
	return ln.a.Close()
}
