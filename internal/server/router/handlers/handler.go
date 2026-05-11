package handlers

import "net"

type MessageHandler interface {
	Handle(conn net.Conn) (bool, error)
}
