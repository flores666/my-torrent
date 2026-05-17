package peerreader

import (
	"errors"
	"my-torrent/lib/peerreader/handshake"
)

type HandshakeMessage struct {
	PstrLen  int      // string length of pstr
	Pstr     []byte   // string identifier of the protocol
	Reserved [8]byte  // 8 reserved bytes.
	InfoHash [20]byte // 20-byte SHA1 hash of the info key in the metainfo file. This is the same info_hash that is transmitted in tracker requests
	PeerId   [20]byte // 20-byte string used as a unique ID for the client
}

type HandshakeReader struct {
}

func NewHandshakeReader() *HandshakeReader {
	return &HandshakeReader{}
}

// Read supports only BitTorrent v1 protocol
func (h *HandshakeReader) Read(body []byte) (any, error) {
	msg := HandshakeMessage{}

	if !handshake.Validate(body) {
		return nil, errors.New("invalid handshake message body")
	}

	msg.PstrLen = int(body[0])
	msg.Pstr = body[1 : msg.PstrLen+1]
	msg.Reserved = [8]byte(body[1+msg.PstrLen : 1+msg.PstrLen+8])
	msg.InfoHash = [20]byte(body[1+msg.PstrLen+8 : 1+msg.PstrLen+8+20])
	msg.PeerId = [20]byte(body[1+msg.PstrLen+8+20 : 1+msg.PstrLen+8+20+20])

	return &msg, nil
}
