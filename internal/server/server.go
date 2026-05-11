package server

import (
	"errors"
	"fmt"
	"net"
	"time"

	"my-torrent/internal/peerclient/peerreader"
	"my-torrent/lib/messagereader"
)

const deadline = 5 * time.Second

type Server struct {
	reader *peerreader.PeerReaderComposite
}

func NewServer(r *peerreader.PeerReaderComposite) *Server {
	return &Server{r}
}

func (s *Server) ListenAndServe() (int, error) {
	port, listener := getFreePort()
	if listener == nil {
		return -1, errors.New("no available ports")
	}

	go func() {
		addr := listener.Addr()

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Error while listening on %s:%s\n", addr.Network(), addr.String())
				return
			}

			go func(conn net.Conn) {
				defer conn.Close()
				conn.SetDeadline(time.Now().Add(deadline))

				body, err := readBody(conn)
				if err != nil {
					fmt.Printf("Error while reading data on %s:%s: %v\n", addr.Network(), addr.String(), err)
					return
				}

				conn.Write(body)
			}(conn)
		}
	}()

	return port, nil
}

func readBody(conn net.Conn) ([]byte, error) {
	var err error
	var body []byte

	if body, err = messagereader.ReadHandshake(conn); err == nil {
		return body, nil
	}

	if body, err = messagereader.ReadMessage(conn); err == nil {
		return body, nil
	}

	return nil, errors.New("unknown message format")
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
