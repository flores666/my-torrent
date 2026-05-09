package connection

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	maxMessageSize = 1 << 20 // 1 MB safety cap
	readTimeout    = 5 * time.Second
)

func ReadHandshake(conn net.Conn) ([]byte, error) {
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

func ReadMessage(conn net.Conn) ([]byte, error) {
	_ = conn.SetReadDeadline(time.Now().Add(readTimeout))

	prefix := make([]byte, 4)
	if _, err := io.ReadFull(conn, prefix[:]); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(prefix[:])

	// keep-alive
	if length == 0 {
		return prefix[:], nil
	}

	if length > maxMessageSize {
		return nil, fmt.Errorf("message too large: %d", length)
	}

	body := make([]byte, length)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, err
	}

	return append(prefix[:], body...), nil
}
