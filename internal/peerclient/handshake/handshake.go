package handshake

import (
	"slices"
)

const (
	msgLen = 68

	pstrLen = 19
	pstr    = "BitTorrent protocol"

	reservedStart = 1 + pstrLen
	reservedLen   = 8

	infoHashStart = reservedStart + reservedLen
	infoHashLen   = 20

	peerIdStart = infoHashStart + infoHashLen
)

// Validate validates response handshake
func Validate(first, second []byte) bool {
	if len(first) != msgLen || len(first) != len(second) {
		return false
	}

	if second[0] != pstrLen {
		return false
	}

	if string(second[1:1+pstrLen]) != pstr {
		return false
	}

	if !slices.Equal(first[infoHashStart:infoHashStart+infoHashLen], second[infoHashStart:infoHashStart+infoHashLen]) {
		return false
	}

	// todo: only for dev, uncomment this
	// if slices.Equal(first[peerIdStart:], second[peerIdStart:]) {
	// 	return false
	// }

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
