package handlers

import (
	"errors"
	"fmt"
	"my-torrent/internal/modelbuilder/torrent"
	"my-torrent/lib/handshake"
	"net"
)

type handler struct {
}

func NewHandshake() MessageHandler {
	return &handler{}
}

func (h *handler) Handle(conn net.Conn) (bool, error) {
	defer conn.Close()

	body, err := handshake.ReadBody(conn)
	if err != nil {
		return false, err
	}

	if !handshake.Validate(body) {
		return false, errors.New("Error while validating handshake")
	}

	hash := body[handshake.InfoHashStart : handshake.InfoHashStart+handshake.InfoHashLen]

	//todo real torrent from db
	t := findTorrent(hash)
	if t == nil {
		return true, fmt.Errorf("could not find torrent with hash %x", hash)
	}

	//todo real id from db
	resp := handshake.BuildBytes(t.Info.Hash, "12345678912345678912")
	_, err = conn.Write(resp)
	if err != nil {
		addr := conn.RemoteAddr()
		return true, fmt.Errorf("could not write handshake data to %s:%s", addr.Network(), addr.String())
	}

	return true, nil
}

func findTorrent(hash []byte) *torrent.Torrent {
	return &torrent.Torrent{
		Info: &torrent.TorrentInfo{Hash: [20]byte(hash)},
	}
}
