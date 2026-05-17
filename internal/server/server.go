package server

import (
	"errors"
	"fmt"
	"net"

	"my-torrent/internal/server/router"
	"my-torrent/internal/storage"
	"my-torrent/lib/peerreader"
)

type Server struct {
	reader  *peerreader.PeerReaderComposite
	router  router.MessageRouter
	storage storage.ServerStorage
}

func NewServer(r *peerreader.PeerReaderComposite, mr router.MessageRouter, ss storage.ServerStorage) *Server {
	return &Server{
		reader:  r,
		router:  mr,
		storage: ss,
	}
}

func (s *Server) ListenAndServe() (int, error) {
	port, listener := getFreePort()
	if listener == nil {
		return -1, errors.New("no available ports")
	}

	err := s.storage.SetPort(port)
	if err != nil {
		fmt.Println(err)
		return -1, errors.New("failed to save port")
	}

	go func() {
		addr := listener.Addr()

		for {
			//todo max connections
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Error while listening on %s:%s\n", addr.Network(), addr.String())
				return
			}

			go s.router.Route(conn)
		}
	}()

	return port, nil
}

func getFreePort() (int, net.Listener) {
	var listener net.Listener
	var err error

	const start int = 6881
	const end int = 6889

	for port := start; port <= end; port++ {
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))

		if err == nil {
			fmt.Println("Using port:", port)
			return port, listener
		}
	}

	return 0, nil
}
