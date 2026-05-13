package peerclient

import (
	"errors"
	"fmt"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/peerclient/peerreader"
	"my-torrent/internal/storage"
	"my-torrent/lib/handshake"
	"net"
	"time"
)

const deadline = 5 * time.Second

type PeerClient interface {
	Handshake(peer *peers.Peer, infoHash [20]byte) (*peerreader.HandshakeMessage, error)
}

type peerClient struct {
	reader   *peerreader.PeerReaderComposite
	tStorage storage.TorrentsStorage
	sStorage storage.ServerStorage
}

func NewPeerClient(r *peerreader.PeerReaderComposite, st storage.TorrentsStorage, ss storage.ServerStorage) PeerClient {
	return &peerClient{
		reader:   r,
		tStorage: st,
		sStorage: ss,
	}
}

func (c *peerClient) Handshake(peer *peers.Peer, infoHash [20]byte) (*peerreader.HandshakeMessage, error) {
	myId, err := c.sStorage.GetId()
	if err != nil {
		return nil, err
	}

	myHandshake := handshake.BuildBytes(infoHash, myId)

	conn, err := sendRequest("tcp", peer.Address(), myHandshake)
	if err != nil {
		if dbErr := c.tStorage.UpdatePeerStatus(infoHash[:], peer.Ip, peer.Port, PeerStatusFailed); dbErr != nil {
			fmt.Printf("error while updating peer %s status, err = %v\n", peer.Address(), dbErr)
		}

		return nil, err
	}

	fmt.Printf("Got response from %s\n", peer.Address())

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

	fmt.Printf("successful handshake, updating peer %s status\n", peer.Address())
	if err = c.tStorage.MarkPeerHandshakeReceived(response.InfoHash[:], string(response.PeerId[:]), peer.Ip, peer.Port, PeerStatusHandshakeReceived); err != nil {
		fmt.Printf("error while updating peer %s status\n", peer.Address())
		return nil, nil
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
