package peerreader

import "errors"

type PeerReaderComposite struct {
	handlers []func(body []byte) (any, error)
}

const knownResponseMessageTypesCount = 12

func NewPeerReaderComposite() *PeerReaderComposite {
	return &PeerReaderComposite{
		handlers: make([]func(body []byte) (any, error), 0, knownResponseMessageTypesCount),
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
