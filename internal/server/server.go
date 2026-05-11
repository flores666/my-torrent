package server

import (
	"errors"
	"fmt"
	"net"
	"time"

	"my-torrent/internal/peerclient/peerreader"
	"my-torrent/internal/server/router"
)

const deadline = 5 * time.Second

type Server struct {
	reader *peerreader.PeerReaderComposite
	router router.MessageRouter
}

func NewServer(r *peerreader.PeerReaderComposite, mr router.MessageRouter) *Server {
	return &Server{
		reader: r,
		router: mr,
	}
}

func (s *Server) ListenAndServe() (int, error) {
	port, listener := getFreePort()
	if listener == nil {
		return -1, errors.New("no available ports")
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
