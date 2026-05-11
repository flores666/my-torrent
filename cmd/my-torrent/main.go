package main

import (
	"fmt"
	"log"
	"my-torrent/internal/connection"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/modelbuilder/torrent"
	"my-torrent/internal/peerclient"
	"my-torrent/internal/peerclient/handshake"
	"my-torrent/internal/peerclient/peerreader"
	"my-torrent/lib/constants"
)

func main() {
	peerReader := newPeerReader()
	server := connection.NewServer(peerReader)
	peerClient := peerclient.NewPeerClient(peerReader)

	port, err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to create new server, err = %v\n", err)
	}

	fmt.Println("Server started")

	torrent, err := torrent.Open("debian.torrent")
	if err != nil {
		log.Fatalln(err.Error())
	}

	//_ = trackerclient.GetPeers(torrent, port)

	peer := peers.PickPeer([]peers.Peer{
		{
			Id:   constants.PEER_ID,
			Ip:   "localhost",
			Port: port,
		},
	}, port)

	fmt.Println("Sending handshake request")

	err = peerClient.Handshake(peer, torrent.Info.Hash, constants.PEER_ID)
	if err != nil {
		log.Fatal("Error while sending handshake request")
	}
}

func newPeerReader() *peerreader.PeerReaderComposite {
	reader := peerreader.NewPeerReaderComposite()
	reader.Register(handshake.NewReader().Read)

	return reader
}
