package p2p

import (
	"errors"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/lib/peerreader/handshake"
	"net"
	"time"
)

const deadline = 5 * time.Second

type Client interface {
	Handshake(peer *peers.Peer, handshake []byte) ([]byte, error)
}

type client struct {
}

func NewClient() Client {
	return &client{}
}

func (c *client) Handshake(peer *peers.Peer, handhsake []byte) ([]byte, error) {
	conn, err := sendRequest("tcp", peer.Address(), handhsake)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	body, err := handshake.ReadBody(conn)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func sendRequest(network, addr string, message []byte) (net.Conn, error) {
	if addr == "" {
		return nil, errors.New("empty address")
	}

	if len(message) == 0 {
		return nil, errors.New("empty message")
	}

	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(deadline))

	_, err = conn.Write(message)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
