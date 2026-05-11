package peerclient

import (
	"errors"
	"fmt"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/peerclient/peerreader"
	"my-torrent/lib/constants"
	"my-torrent/lib/handshake"
	"net"
	"time"
)

const deadline = 5 * time.Second

type PeerClient interface {
	Handshake(peer *peers.Peer, infoHash [20]byte) (*peerreader.HandshakeMessage, error)
}

type peerClient struct {
	reader *peerreader.PeerReaderComposite
}

func NewPeerClient(r *peerreader.PeerReaderComposite) PeerClient {
	return &peerClient{r}
}

func (c *peerClient) Handshake(peer *peers.Peer, infoHash [20]byte) (*peerreader.HandshakeMessage, error) {
	myHandshake := handshake.BuildBytes(infoHash, constants.PEER_ID)

	conn, err := sendRequest("tcp", peer.Address(), myHandshake)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	body, err := handshake.ReadBody(conn)
	if err != nil {
		return nil, err
	}

	if !handshake.ValidateResponse(myHandshake, body) {
		addr := conn.RemoteAddr()
		err = fmt.Errorf("invalid handshake response from %s:%s", addr.Network(), addr.String())
		return nil, err
	}

	response, err := peerreader.ReadAs[*peerreader.HandshakeMessage](c.reader, body)
	if err != nil {
		return nil, err
	}

	return response, nil
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
