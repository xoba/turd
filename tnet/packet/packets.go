package packet

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Connection interface {
	Receiver
	Sender
	io.Closer
}

const MaxBuffer = 1000 * 1000

func NewConn(rwc io.ReadWriteCloser) Connection {
	return packetconn{rwc}
}

type Receiver interface {
	Receive() ([]byte, error)
}
type Sender interface {
	Send([]byte) error
}

type packetconn struct {
	io.ReadWriteCloser
}

func (pc packetconn) Receive() ([]byte, error) {
	return receive(pc)
}

func (pc packetconn) Send(buf []byte) error {
	return send(pc, buf)
}

func send(w io.Writer, buf []byte) error {
	var bufSize uint64 = uint64(len(buf))
	if bufSize > MaxBuffer {
		return fmt.Errorf("can't handle buffers bigger than %d bytes", MaxBuffer)
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

func receive(r io.Reader) ([]byte, error) {
	var bufSize uint64
	if err := binary.Read(r, binary.BigEndian, &bufSize); err != nil {
		return nil, err
	}
	if bufSize > MaxBuffer {
		return nil, fmt.Errorf("can't handle buffers bigger than %d bytes", MaxBuffer)
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
