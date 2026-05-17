package router

import (
	"fmt"
	"my-torrent/lib/constants"
	"net"
)

type MessageRouter interface {
	Register(func(conn net.Conn) (bool, error))
	Route(conn net.Conn)
}

type router struct {
	handlers []func(conn net.Conn) (bool, error)
}

func NewMessageRouter() MessageRouter {
	return &router{
		handlers: make([]func(conn net.Conn) (bool, error), 0, constants.P2P_MESSAGE_TYPES_COUNT),
	}
}

func (r *router) Route(conn net.Conn) {
	var found bool
	var err error

	for _, handle := range r.handlers {
		if found, err = handle(conn); found {
			if err != nil {
				// todo log in handler
				fmt.Printf("Error while handling request: %v", err)
			}

			return
		}
	}

	addr := conn.RemoteAddr()
	fmt.Printf("Could not find route for message from %s:%s\n", addr.Network(), addr.String())
	conn.Close()
}

func (r *router) Register(h func(conn net.Conn) (bool, error)) {
	r.handlers = append(r.handlers, h)
}
