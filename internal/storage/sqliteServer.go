package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

type sqlliteServerStorage struct {
	db *sql.DB
}

func NewSqlLiteServerStorage(s *sql.DB) ServerStorage {
	return &sqlliteServerStorage{
		db: s,
	}
}

func (s *sqlliteServerStorage) SetPort(port int) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port")
	}

	const query = `
		UPDATE client_identity
		SET
			port = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = 1;
	`

	_, err := s.db.Exec(query, port)
	if err != nil {
		return fmt.Errorf("update port: %w", err)
	}

	return nil
}

func (s *sqlliteServerStorage) GetId() (string, error) {
	const query = `
		SELECT peer_id
		FROM client_identity
		WHERE id = 1;
	`

	var peerID string

	err := s.db.QueryRow(query).Scan(&peerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return peerID, nil
}

func (s *sqlliteServerStorage) SetId(peerID string) error {
	if len(peerID) != 20 {
		return fmt.Errorf("peer id must be 20 characters")
	}

	const query = `
		INSERT INTO client_identity (id, peer_id)
		VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET
			peer_id = excluded.peer_id,
			updated_at = CURRENT_TIMESTAMP;
	`

	_, err := s.db.Exec(query, peerID)
	if err != nil {
		return err
	}

	return nil
}

func (s *sqlliteServerStorage) GetPort() (int, error) {
	const query = `
		SELECT port
		FROM client_identity
		WHERE id = 1;
	`

	var port int

	err := s.db.QueryRow(query).Scan(&port)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return port, nil
}
