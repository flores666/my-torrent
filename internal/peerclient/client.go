package peerclient

import (
	"errors"
	"fmt"
	"my-torrent/internal/connection"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/peerclient/handshake"
	"my-torrent/internal/peerclient/peerreader"
	"my-torrent/lib/constants"
	"net"
	"time"
)

const deadline = 5 * time.Second

type PeerClient interface {
	Handshake(peer *peers.Peer, infoHash [20]byte) (*handshake.Message, error)
}

type peerClient struct {
	reader *peerreader.PeerReaderComposite
}

func NewPeerClient(r *peerreader.PeerReaderComposite) PeerClient {
	return &peerClient{r}
}

func (c *peerClient) Handshake(peer *peers.Peer, infoHash [20]byte) (*handshake.Message, error) {
	myHandshake := handshake.BuildBytes(infoHash, constants.PEER_ID)

	conn, err := sendRequest("tcp", peer.Address(), myHandshake)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	addr := conn.RemoteAddr()

	body, err := connection.ReadHandshake(conn)
	if err != nil {
		return nil, err
	}

	if !handshake.Validate(myHandshake, body) {
		err = fmt.Errorf("invalid handshake response from %s:%s", addr.Network(), addr.String())
		return nil, err
	}

	response, err := peerreader.ReadAs[*handshake.Message](c.reader, body)
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
