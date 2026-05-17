package p2p

import (
	"fmt"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/storage"
	"my-torrent/lib/peerreader/handshake"
	"my-torrent/lib/peerstatuses"
)

type Service interface {
	Handshake(peer *peers.Peer, infoHash [20]byte) (*Session, error)
}

type service struct {
	client          Client
	torrentsStorage storage.TorrentsStorage
	serverStorage   storage.ServerStorage
}

func NewService(c Client, ts storage.TorrentsStorage, ss storage.ServerStorage) Service {
	return &service{
		client:          c,
		serverStorage:   ss,
		torrentsStorage: ts,
	}
}

func (s *service) Handshake(peer *peers.Peer, infoHash [20]byte) (*Session, error) {
	myId, err := s.serverStorage.GetId()
	if err != nil {
		return nil, err
	}

	myHandshake := handshake.BuildBytes(infoHash, myId)
	session, err := s.client.Handshake(peer, myHandshake)

	if err != nil {
		fmt.Println(err)
		s.updatePeerStatus(infoHash[:], peer.Ip, peer.Port, peerstatuses.Failed)

		return nil, err
	}

	fmt.Printf("successful handshake, updating peer %s status\n", peer.Address())

	if dbErr := s.updatePeerStatus(infoHash[:], peer.Ip, peer.Port, peerstatuses.Ready); dbErr != nil {
		fmt.Println(dbErr)
	}

	return session, nil
}

func (s *service) updatePeerStatus(torrentHash []byte, ip string, port int, status string) error {
	var err error = nil

	if err = s.torrentsStorage.UpdatePeerStatus(torrentHash, ip, port, status); err != nil {
		fmt.Println(err)

		if dbErr := s.torrentsStorage.UpdatePeerStatus(torrentHash, ip, port, peerstatuses.Failed); dbErr != nil {
			fmt.Println(dbErr)
		}
	}

	return err
}
