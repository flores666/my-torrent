package sessionroutes

import (
	"my-torrent/internal/storage"
	"my-torrent/lib/constants"
	"my-torrent/lib/messages"
)

type BitfieldRoute struct {
	storage storage.TorrentsStorage
}

func NewBitfieldMessageHandler(ts storage.TorrentsStorage) Route {
	return &BitfieldRoute{ts}
}

func (h *BitfieldRoute) Handle(peerId, torrentInfoHash, msg []byte) error {
	length := messages.GetLength(msg)
	bitfield := msg[5:length]

	return h.storage.UpdatePeerBitfield(peerId, torrentInfoHash, bitfield)
}

func (h *BitfieldRoute) GetMessageId() uint8 {
	return constants.MID_BITFIELD
}
