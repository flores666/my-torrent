package p2p

import (
	"net"
)

type Session struct {
	Conn net.Conn

	Ip   string
	Port int

	PeerId [20]byte

	AmChoking      bool
	PeerChoking    bool
	AmInterested   bool
	PeerInterested bool

	Bitfield []byte
}
