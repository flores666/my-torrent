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

	Done chan struct{}

	closeOnce sync.Once
	wm        sync.Mutex
}

func (s *Session) Close() error {
	var err error

	s.closeOnce.Do(func() {
		close(s.Done)
		err = s.Conn.Close()
	})

	return err
}

// Write writes data to the connection with mutex lock
func (s *Session) Write(data []byte) (int, error) {
	s.wm.Lock()
	defer s.wm.Unlock()
	return s.Conn.Write(data)
}
