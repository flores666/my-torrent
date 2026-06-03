package messages

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

// ReadBody reads bytes from conn and returns bytes of message
func ReadBody(conn net.Conn) ([]byte, error) {
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

func GetLength(msg []byte) uint32 {
	return binary.BigEndian.Uint32(msg[:4])
}

func GetId(msg []byte) uint8 {
	return msg[4]
}
