package storage

import "my-torrent/internal/modelbuilder/torrent"

type Storage interface {
	Save(*torrent.Torrent) error
	Find([]byte) (*torrent.Torrent, error)
	Remove([]byte) error
}
