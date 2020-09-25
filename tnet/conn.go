package tnet

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"math/big"
)

type conn struct {
	self           cipher.AEAD // key we encrypt with
	other          cipher.AEAD // key we decrypt with
	remote         Node
	c              RawConn
	sent, received *big.Int
}

func newConn(c RawConn, self, other []byte) (*conn, error) {
	fmt.Printf("self = 0x%x\n", self)
	fmt.Printf("other = 0x%x\n", other)
	if len(self) != 32 || len(other) != 32 {
		return nil, fmt.Errorf("need 256 bit keys")
	}
	self0, err := aes.NewCipher(self)
	if err != nil {
		return nil, err
	}
	self1, err := cipher.NewGCM(self0)
	if err != nil {
		return nil, err
	}
	other0, err := aes.NewCipher(other)
	if err != nil {
		return nil, err
	}
	other1, err := cipher.NewGCM(other0)
	if err != nil {
		return nil, err
	}
	return &conn{
		c:        c,
		self:     self1,
		other:    other1,
		sent:     big.NewInt(0),
		received: big.NewInt(0),
	}, nil
}

func (c conn) Remote() Node {
	return c.remote
}

func (c *conn) Receive() ([]byte, error) {
	cipherText, err := receive(c.c)
	if err != nil {
		return nil, err
	}
	plaintext, err := c.other.Open(nil, c.nonce(c.received, c.other.NonceSize()), cipherText, nil)
	if err != nil {
		return nil, err
	}
	c.received = inc(c.received)
	return plaintext, nil
}

func (c *conn) nonce(i *big.Int, n int) []byte {
	return i.FillBytes(make([]byte, n))
}

func (c *conn) Send(buf []byte) (err error) {
	cipherText := c.self.Seal(nil, c.nonce(c.sent, c.self.NonceSize()), buf, nil)
	c.sent = inc(c.sent)
	return send(c.c, cipherText)
}

func (c conn) Close() error {
	return c.c.Close()
}

func inc(i *big.Int) *big.Int {
	var z big.Int
	one := big.NewInt(1)
	z.Add(i, one)
	return &z
}
