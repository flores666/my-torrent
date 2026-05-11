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

const (
	pstr     = "BitTorrent protocol"
	reserved = "00000000"
)

const deadline = 5 * time.Second

type PeerClient interface {
	Handshake(peer *peers.Peer, infoHash [20]byte) error
}

type peerClient struct {
	reader *peerreader.PeerReaderComposite
}

func NewPeerClient(r *peerreader.PeerReaderComposite) PeerClient {
	return &peerClient{r}
}

func (c *peerClient) Handshake(peer *peers.Peer, infoHash [20]byte) error {
	myHandshake := createHandshake(infoHash, constants.PEER_ID)

	conn, err := sendRequest("tcp", peer.Address(), myHandshake)
	if err != nil {
		return err
	}

	defer conn.Close()

	addr := conn.RemoteAddr()

	body, err := connection.ReadHandshake(conn)
	if err != nil {
		fmt.Printf("Error while reading data on %s:%s: %v\n", addr.Network(), addr.String(), err)
		return err
	}

	if !handshake.Validate(myHandshake, body) {
		err = fmt.Errorf("invalid handshake response from %s:%s", addr.Network(), addr.String())
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("%x\n", body)

	return nil
}

func createHandshake(infoHash [20]byte, peerId string) []byte {
	buff := make([]byte, 0, 68)
	buff = append(buff, 19)
	buff = append(buff, []byte("BitTorrent protocol")...)
	buff = append(buff, make([]byte, 8)...)
	buff = append(buff, infoHash[:]...)
	buff = append(buff, []byte(peerId)...)
	return buff
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
