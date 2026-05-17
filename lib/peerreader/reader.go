package peerreader

type PeerReader interface {
	Read(body []byte) (any, error)
}
