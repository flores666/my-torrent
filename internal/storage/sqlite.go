package storage

import (
	"database/sql"
	"my-torrent/internal/modelbuilder/torrent"
)

type sqlliteStorage struct {
	db *sql.DB
}

func NewSqlLiteStorage(s *sql.DB) Storage {
	return &sqlliteStorage{
		db: s,
	}
}

func (s *sqlliteStorage) Save(*torrent.Torrent) error {
	return nil
}

func (s *sqlliteStorage) Find([]byte) (*torrent.Torrent, error) {
	return nil, nil
}

func (s *sqlliteStorage) Remove([]byte) error {
	return nil
}
