package main

import (
	"fmt"
	"log"
	"my-torrent/internal/modelbuilder/torrent"
	"my-torrent/internal/p2p"
	"my-torrent/internal/server"
	"my-torrent/internal/server/router"
	"my-torrent/internal/server/router/handlers"
	"my-torrent/internal/storage"
	"my-torrent/internal/tracker"
	"my-torrent/lib/constants"
	"my-torrent/lib/db"
	"my-torrent/lib/peerreader"
)

func main() {
	//#region dependencies
	db := db.MustLoadNewSqlliteDB("sqlite", "egt.db")
	defer db.Close()

	tStorage := storage.NewSqlLiteTorrentStorage(db)
	sStorage := storage.NewSqlLiteServerStorage(db)

	peerReader := newPeerReader()

	peerService := p2p.NewService(p2p.NewClient(peerReader), tStorage, sStorage)
	trackerService := tracker.NewService(tracker.NewClient(), tStorage, sStorage)
	_ = trackerService

	sessionsManager := p2p.NewSessionManager(tStorage)

	server := server.NewServer(peerReader, newMessagesRouter(tStorage, sStorage, sessionsManager, peerReader), sStorage)
	port, err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to create new server, err = %v\n", err)
	}

	fmt.Printf("Server started on port %d\n", port)

	// #endregion dependencies

	torrent, err := torrent.Open("debian.torrent")
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = tStorage.Save(torrent)
	if err != nil {
		log.Fatalln(err.Error())
	}

	pList, err := trackerService.Announce(torrent)
	if err != nil {
		log.Fatalln(err)
	}

	for i, connectedCount := 0, 0; i < len(pList) && connectedCount < constants.MAX_DOWNLOAD_SLOTS; i++ {
		peer := pList[i]

		fmt.Printf("Sending handshake request to %s\n", peer.Address())

		session, err := peerService.DialHandshake(peer, torrent.Info.Hash)
		if err != nil {
			log.Fatalf("Error while sending handshake request, error = %v", err)
		}

		err = sessionsManager.Download(session)
		if err != nil {
			connectedCount++
		}
	}
}

func newPeerReader() *peerreader.PeerReaderComposite {
	r := peerreader.NewPeerReaderComposite()
	r.Register(peerreader.NewHandshakeReader().Read)

	return r
}

func newMessagesRouter(tStorage storage.TorrentsStorage, sStorage storage.ServerStorage, sm p2p.SessionManager, pr *peerreader.PeerReaderComposite) router.MessageRouter {
	r := router.NewMessageRouter()
	r.Register(handlers.NewHandshake(tStorage, sStorage, sm, pr).Handle)

	return r
}
