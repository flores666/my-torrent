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

	external, _ := peerreader.ReadAs[*peerreader.HandshakeMessage](h.reader, body)

	t, _ := h.tStorage.Find(external.InfoHash[:])
	if t == nil {
		return true, fmt.Errorf("could not find torrent with hash %x", external.InfoHash[:])
	}

	conn.SetDeadline(time.Now().Add(constants.TIMEOUT))

	myHandshake := h.buildHandshake(external.InfoHash)
	_, err = conn.Write(myHandshake)
	if err != nil {
		addr := conn.RemoteAddr()
		return true, fmt.Errorf("could not write handshake to %s:%s", addr.Network(), addr.String())
	}

	bitfield, err := h.tStorage.GetPieces(external.InfoHash)
	if err == nil && len(bitfield) > 0 {
		_, err = conn.Write(bitfield)
		if err != nil {
			addr := conn.RemoteAddr()
			return true, fmt.Errorf("could not write bitfield to %s:%s", addr.Network(), addr.String())
		}
	}

	ip, port := getIpAndPort(conn.RemoteAddr())

	session := &p2p.Session{
		Conn:     conn,
		Ip:       ip,
		Port:     port,
		PeerId:   external.PeerId,
		InfoHash: external.InfoHash,
	}

	err = h.manager.Upload(session)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (h *handler) buildHandshake(hash [20]byte) []byte {
	myId, _ := h.sStorage.GetId()

	return handshake.BuildBytes(hash, myId)
}

func getIpAndPort(addr net.Addr) (string, int) {
	host, portStr, _ := net.SplitHostPort(addr.String())

	port, _ := strconv.Atoi(portStr)
	return host, port
}
