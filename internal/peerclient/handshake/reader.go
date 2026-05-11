package handshake

import (
	"encoding/hex"
	"errors"
)

const (
	handshakeLenV1 = 68
	protocol       = "BitTorrent protocol"
)

type HandshakeMessage struct {
	PstrLen  int      // string length of pstr
	Pstr     []byte   // string identifier of the protocol
	Reserved [8]byte  // 8 reserved bytes.
	InfoHash [20]byte // 20-byte SHA1 hash of the info key in the metainfo file. This is the same info_hash that is transmitted in tracker requests
	PeerId   [20]byte // 20-byte string used as a unique ID for the client
}

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

// Read supports only BitTorrent v1 protocol
func (h *Reader) Read(body []byte) (any, error) {
	msg := HandshakeMessage{}

	if !validate(body) {
		return nil, errors.New("invalid handshake message body")
	}

	msg.PstrLen = int(body[0])
	msg.Pstr = body[1 : msg.PstrLen+1]
	msg.Reserved = [8]byte(body[1+msg.PstrLen : 1+msg.PstrLen+8])
	msg.InfoHash = [20]byte(body[1+msg.PstrLen+8 : 1+msg.PstrLen+8+20])
	msg.PeerId = [20]byte(body[1+msg.PstrLen+8+20 : 1+msg.PstrLen+8+20+20])

	return &msg, nil
}

func validate(msg []byte) bool {
	if len(msg) != handshakeLenV1 {
		return false
	}

	if msg[0] != 19 {
		return false
	}

	if string(msg[1:1+19]) != protocol {
		return false
	}

	_, err := hex.DecodeString(string(msg[1+19+8 : 1+19+8+20]))
	if err != nil {
		return false
	}

	_, err = hex.DecodeString(string(msg[1+19+8+20:]))
	if err != nil {
		return false
	}

	return true
}
