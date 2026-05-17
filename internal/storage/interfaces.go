package storage

import (
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/modelbuilder/torrent"
)

type TorrentsStorage interface {
	Save(*torrent.Torrent) error
	Find([]byte) (*torrent.Torrent, error)
	Remove([]byte) error
	SavePeers(torrentHash []byte, peers []peers.Peer) error
	GetPeers(torrentHash []byte) ([]peers.Peer, error)
	UpdatePeerStatus(torrentHash []byte, ip string, port int, status string) error
	MarkPeerHandshakeReceived(infoHash []byte, ip string, port int) error
}

type ServerStorage interface {
	GetId() (string, error)
	SetId(string) error
	GetPort() (int, error)
	SetPort(int) error
}
