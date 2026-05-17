package handlers

import (
	"errors"
	"fmt"
	"my-torrent/internal/storage"
	"my-torrent/lib/peerreader/handshake"
	"net"
)

type handler struct {
	tStorage storage.TorrentsStorage
	sStorage storage.ServerStorage
}

func NewHandshake(ts storage.TorrentsStorage, ss storage.ServerStorage) MessageHandler {
	return &handler{
		tStorage: ts,
		sStorage: ss,
	}
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

	t, _ := h.tStorage.Find(hash)
	if t == nil {
		return true, fmt.Errorf("could not find torrent with hash %x", hash)
	}

	myId, err := h.sStorage.GetId()
	if err != nil || myId == "" {
		return true, fmt.Errorf("own peer id is not configured")
	}

	resp := handshake.BuildBytes(t.Info.Hash, "12345678912345678912") // todo
	_, err = conn.Write(resp)
	if err != nil {
		addr := conn.RemoteAddr()
		return true, fmt.Errorf("could not write handshake data to %s:%s", addr.Network(), addr.String())
	}

	return true, nil
}
