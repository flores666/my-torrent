package p2p

import (
	"net"
	"sync"
)

type Session struct {
	Conn net.Conn

	Ip   string
	Port int

	PeerId   [20]byte
	InfoHash [20]byte

	AmChoking      bool
	PeerChoking    bool
	AmInterested   bool
	PeerInterested bool

	Bitfield []byte

	Done      chan struct{}
	closeOnce sync.Once
}

func (s *Session) Close() error {
	var err error

	s.closeOnce.Do(func() {
		close(s.Done)
		err = s.Conn.Close()
	})

	return err
}
