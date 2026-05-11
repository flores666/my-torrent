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

func Validate(my, resp []byte) bool {
	if len(my) != msgLen || len(my) != len(resp) {
		return false
	}

	if resp[0] != pstrLen {
		return false
	}

	if string(resp[1:1+pstrLen]) != pstr {
		return false
	}

	if !slices.Equal(my[infoHashStart:infoHashStart+infoHashLen], resp[infoHashStart:infoHashStart+infoHashLen]) {
		return false
	}

	// todo: only for tests
	// if slices.Equal(my[peerIdStart:], resp[peerIdStart:]) {
	// 	return false
	// }

	return true
}
