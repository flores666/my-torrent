package main

import (
	"fmt"
	"log"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/modelbuilder/torrent"
	"my-torrent/internal/peerclient"
	"my-torrent/internal/peerclient/peerreader"
	"my-torrent/internal/server"
	"my-torrent/internal/server/router"
	"my-torrent/internal/server/router/handlers"
	"my-torrent/internal/storage"
	"my-torrent/lib/constants"
	"my-torrent/lib/db"
)

func main() {
	db := db.MustLoadNewSqlliteDB("sqlite", "egt.db")
	defer db.Close()

	store := storage.NewSqlLiteStorage(db)
	peerReader := newPeerReader()

	server := server.NewServer(peerReader, newMessagesRouter())
	peerClient := peerclient.NewPeerClient(peerReader, store)

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

	response, err := peerClient.Handshake(peer, torrent.Info.Hash)
	if err != nil {
		log.Fatalf("Error while sending handshake request, error = %v", err)
	}
	_ = response
}

func newPeerReader() *peerreader.PeerReaderComposite {
	r := peerreader.NewPeerReaderComposite()
	r.Register(peerreader.NewHandshakeReader().Read)

	return r
}

func newMessagesRouter() router.MessageRouter {
	r := router.NewMessageRouter()
	r.Register(handlers.NewHandshake().Handle)

	return r
}
