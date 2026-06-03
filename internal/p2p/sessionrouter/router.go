package sessionrouter

import (
	"fmt"
	"my-torrent/internal/p2p/sessionrouter/sessionroutes"
	"my-torrent/lib/messages"
)

type Router interface {
	Register(sessionroutes.Route)
	Route(peerId, torrentInfoHash, message []byte) error
}

type router struct {
	handlers []sessionroutes.Route
}

func NewRouter() Router {
	return &router{}
}

func (r *router) Register(h sessionroutes.Route) {
	r.handlers = append(r.handlers, h)
}

func (r *router) Route(peerId, torrentInfoHash, message []byte) error {
	msgId := messages.GetId(message)
	handler := findHandler(msgId, r.handlers)
	if handler == nil {
		return fmt.Errorf("Not found message handler with id %d\n", msgId)
	}

	fmt.Printf("Handling message with id %d\n", msgId)

	return handler.Handle(peerId, torrentInfoHash, message)
}

func findHandler(id uint8, handlers []sessionroutes.Route) sessionroutes.Route {
	for _, h := range handlers {
		if h.GetMessageId() == id {
			return h
		}
	}

	return nil
}
