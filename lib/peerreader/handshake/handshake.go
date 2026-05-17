package handshake

import (
	"fmt"
	"io"
	"net"
	"slices"
	"time"
)

const (
	MsgLen = 68

	PstrLen = 19
	Pstr    = "BitTorrent protocol"

	ReservedStart = 1 + PstrLen
	ReservedLen   = 8

	InfoHashStart = ReservedStart + ReservedLen
	InfoHashLen   = 20

	PeerIdStart = InfoHashStart + InfoHashLen
)

const readTimeout = 5 * time.Second

// Validate validates response handshake
func ValidateResponse(first, second []byte) bool {
	if len(first) != MsgLen || len(first) != len(second) {
		return false
	}

	if second[0] != PstrLen {
		return false
	}

	if string(second[1:1+PstrLen]) != Pstr {
		return false
	}

	if !slices.Equal(first[InfoHashStart:InfoHashStart+InfoHashLen], second[InfoHashStart:InfoHashStart+InfoHashLen]) {
		return false
	}

	if slices.Equal(first[PeerIdStart:], second[PeerIdStart:]) {
		return false
	}

	return true
}

func Validate(msg []byte) bool {
	if len(msg) != MsgLen {
		return false
	}

	if msg[0] != PstrLen {
		return false
	}

	if string(msg[1:1+PstrLen]) != Pstr {
		return false
	}

	return true
}

// BuildBytes creates handshake request body
func BuildBytes(infoHash [20]byte, peerId string) []byte {
	buff := make([]byte, 0, 68)
	buff = append(buff, 19)
	buff = append(buff, []byte("BitTorrent protocol")...)
	buff = append(buff, make([]byte, 8)...)
	buff = append(buff, infoHash[:]...)
	buff = append(buff, []byte(peerId)...)
	return buff
}

// ReadBody reads bytes from conn and returns bytes of handshake
func ReadBody(conn net.Conn) ([]byte, error) {
	_ = conn.SetReadDeadline(time.Now().Add(readTimeout))

	pstrlen := make([]byte, 1)

	_, err := io.ReadFull(conn, pstrlen)
	if err != nil {
		return nil, err
	}

	if pstrlen[0] != 19 {
		return nil, fmt.Errorf("invalid pstrlen: %d", pstrlen[0])
	}

	rest := make([]byte, 67) // fixed for BitTorrent

	_, err = io.ReadFull(conn, rest)

	if err != nil {
		return nil, err
	}

	return append(pstrlen, rest...), nil
}
