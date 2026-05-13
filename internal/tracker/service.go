package tracker

import (
	"fmt"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/modelbuilder/torrent"
	"my-torrent/internal/storage"
)

type Service interface {
	Announce(torrent *torrent.Torrent) []peers.Peer
}

type service struct {
	sStorage storage.ServerStorage
	tStorage storage.TorrentsStorage
	client   Client
	port     int
}

func NewService(cl Client, ts storage.TorrentsStorage, ss storage.ServerStorage) Service {
	return &service{
		sStorage: ss,
		tStorage: ts,
		client:   cl,
		port:     -1,
	}
}

func (s *service) Announce(torrent *torrent.Torrent) []peers.Peer {
	args, err := s.getArguments(torrent)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	body, err := s.client.Announce(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	response, err := ParseAnnounceResponse(body)
	if err != nil {
		return nil
	}

	fmt.Printf("Saving peers from torrent %x\n", torrent.Info.Hash)

	err = s.tStorage.SavePeers(torrent.Info.Hash[:], response.Peers)
	if err != nil {
		fmt.Printf("Error while savin peers from torrent %x\n", torrent.Info.Hash)
		return nil
	}

	_ = response.Interval // todo: what to do with this?

	return response.Peers
}

func (s *service) getPort() int {
	if s.port != -1 {
		return s.port
	}

	port, err := s.sStorage.GetPort()
	if err != nil {
		fmt.Printf("Error while geting port from db, err = %v\n", err)
		return -1
	}

	s.port = port

	return s.port
}

func (s *service) getArguments(torrent *torrent.Torrent) (arguments, error) {
	event := "started" // stopped completed

	var uploaded int64 = 0
	var downloaded int64 = 0

	left := torrent.Info.Length
	if left == 0 {
		for _, f := range torrent.Info.Files {
			left = f.Length
		}
	}

	myId, err := s.sStorage.GetId()
	if err != nil {
		return arguments{}, err
	}

	port := s.getPort()

	args := arguments{
		peerId:       myId,
		uploaded:     uploaded,
		downloaded:   downloaded,
		left:         left,
		event:        event,
		infoHash:     torrent.Info.Hash,
		port:         port,
		announce:     torrent.Announce,
		announcelist: torrent.AnnounceList,
	}

	return args, nil
}
