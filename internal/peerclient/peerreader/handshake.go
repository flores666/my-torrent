package peerreader

type HandshakeReader struct {
}

func NewHandshakeReader() *HandshakeReader {
	return &HandshakeReader{}
}

func (h *HandshakeReader) Read(body []byte) (any, error) {
	return &HandshakeMessage{
		Message: body,
	}, nil
}
