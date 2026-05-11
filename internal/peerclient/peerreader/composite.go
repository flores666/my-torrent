package peerreader

import (
	"errors"
	"my-torrent/lib/constants"
)

type PeerReaderComposite struct {
	handlers []func(body []byte) (any, error)
}

func NewPeerReaderComposite() *PeerReaderComposite {
	return &PeerReaderComposite{
		handlers: make([]func(body []byte) (any, error), 0, constants.P2P_MESSAGE_TYPES_COUNT),
	}
}

func (r *PeerReaderComposite) Register(f func(body []byte) (any, error)) {
	r.handlers = append(r.handlers, f)
}

func ReadAs[T any](reader *PeerReaderComposite, body []byte) (T, error) {
	var zero T

	for _, handle := range reader.handlers {
		if message, err := handle(body); err == nil {
			if v, ok := message.(T); ok {
				return v, nil
			}
		}
	}

	return zero, errors.New("unhandled message content")
}
