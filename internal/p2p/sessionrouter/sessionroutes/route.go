package sessionroutes

type Route interface {
	GetMessageId() uint8
	Handle(peerId, torrentInfoHash, msg []byte) error
}
