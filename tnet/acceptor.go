package tnet

import (
	"io"
	"net"
)

type Acceptor interface {
	Accept() (io.ReadWriteCloser, error)
	io.Closer
}

type acceptor struct {
	ln net.Listener
}

func (a acceptor) Accept() (io.ReadWriteCloser, error) {
	return a.ln.Accept()
}

func (a acceptor) Close() error {
	return a.ln.Close()
}
