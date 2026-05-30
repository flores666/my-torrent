package tracker

import (
	"my-torrent/internal/modelbuilder/peers"
)

type arguments struct {
	announce     string
	announcelist [][]string
	peerId       string
	port         int
	uploaded     int64
	downloaded   int64
	left         int64
	event        string
	infoHash     [20]byte
}

type AnnounceResponse struct {
	Interval int
	Peers    []*peers.Peer
}
