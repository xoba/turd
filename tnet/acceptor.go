package tnet

import (
	"io"
	"net"
)

type Acceptor interface {
	Accept() (RawConn, error)
	Close() error
}

type RawConn io.ReadWriteCloser

type acceptor struct {
	ln net.Listener
}

func (a acceptor) Accept() (RawConn, error) {
	return a.ln.Accept()
}

func (a acceptor) Close() error {
	return a.ln.Close()
}
