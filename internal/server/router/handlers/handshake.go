package handlers

import (
	"errors"
	"fmt"
	"my-torrent/internal/p2p"
	"my-torrent/internal/storage"
	"my-torrent/lib/constants"
	"my-torrent/lib/peerreader"
	"my-torrent/lib/peerreader/handshake"
	"net"
	"strconv"
	"strings"
	"time"
)

type handler struct {
	tStorage storage.TorrentsStorage
	sStorage storage.ServerStorage
	manager  p2p.SessionManager
	reader   *peerreader.PeerReaderComposite
}

func NewHandshake(ts storage.TorrentsStorage, ss storage.ServerStorage, sm p2p.SessionManager, r *peerreader.PeerReaderComposite) MessageHandler {
	return &handler{
		tStorage: ts,
		sStorage: ss,
		manager:  sm,
		reader:   r,
	}
}

func (h *handler) Handle(conn net.Conn) (_ bool, err error) {
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

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

	resp := handshake.BuildBytes(t.Info.Hash, myId)

	conn.SetDeadline(time.Now().Add(constants.TIMEOUT))

	_, err = conn.Write(resp)
	if err != nil {
		addr := conn.RemoteAddr()
		return true, fmt.Errorf("could not write handshake data to %s:%s", addr.Network(), addr.String())
	}

	// todo: _, err = conn.Write(bitfield)

	inHandshake, _ := peerreader.ReadAs[*peerreader.HandshakeMessage](h.reader, body)
	ip, port := getIpAndPort(conn.RemoteAddr())

	session := &p2p.Session{
		Conn:     conn,
		Ip:       ip,
		Port:     port,
		PeerId:   inHandshake.PeerId,
		InfoHash: inHandshake.InfoHash,
	}
	err = h.manager.Start(session)

	if err != nil {
		return true, err
	}

	return true, nil
}

func getIpAndPort(addr net.Addr) (string, int) {
	parts := strings.Split(addr.String(), ":")
	port, _ := strconv.Atoi(parts[1])
	return parts[0], port
}
