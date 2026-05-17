package p2p

import (
	"errors"
	"fmt"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/lib/constants"
	"my-torrent/lib/peerreader"
	"my-torrent/lib/peerreader/handshake"
	"net"
	"time"
)

type Client interface {
	DialHandshake(peer *peers.Peer, handshake []byte) (*Session, error)
}

type client struct {
	reader *peerreader.PeerReaderComposite
}

func NewClient(r *peerreader.PeerReaderComposite) Client {
	return &client{r}
}

func (c *client) DialHandshake(peer *peers.Peer, handhsake []byte) (*Session, error) {
	conn, err := sendRequest("tcp", peer.Address(), handhsake)
	if err != nil {
		return nil, err
	}

	data, err := handshake.ReadBody(conn)
	if err != nil {
		return nil, err
	}

	if !handshake.ValidateResponse(handhsake, data) {
		err = fmt.Errorf("invalid handshake response from %s", peer.Address())
		return nil, err
	}

	response, err := peerreader.ReadAs[*peerreader.HandshakeMessage](c.reader, data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Session{
		Conn:     conn,
		PeerId:   response.PeerId,
		Ip:       peer.Ip,
		Port:     peer.Port,
		InfoHash: response.InfoHash,
		Done:     make(chan struct{}),
	}, nil
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

	conn.SetDeadline(time.Now().Add(constants.TIMEOUT))

	_, err = conn.Write(message)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
