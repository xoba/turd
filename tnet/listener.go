package tnet

import "net"

type listener struct {
	key *PrivateKey
	ln   net.Listener
}

func (ln listener) Accept() (Conn, error) {
	c, err := ln.ln.Accept()
	if err != nil {
		return nil, err
	}
	// receive other's address
	remote, err := receive(c)
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
	buf1, err := ln.key.Public().MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := send(c, buf1); err != nil {
		return nil, err
	}
	cn := conn{
		remote: string(remote),
		self:   ln.key,
		other:  &other,
		c:      c,
	}
	return cn, nil
}

func (ln listener) Close() error {
	return ln.ln.Close()
}
