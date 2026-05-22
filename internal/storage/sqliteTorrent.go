package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/modelbuilder/torrent"
	"strings"
	"time"
)

type sqlliteTorrentStorage struct {
	db *sql.DB
}

func NewSqlLiteTorrentStorage(s *sql.DB) TorrentsStorage {
	return &sqlliteTorrentStorage{
		db: s,
	}
}

func (s *sqlliteTorrentStorage) Save(t *torrent.Torrent) error {
	if t == nil {
		return errors.New("torrent is nil")
	}

	if t.Info == nil {
		return errors.New("torrent info is nil")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
	INSERT OR REPLACE INTO torrents (
		hash,
		announce,
		comment,
		created_by,
		creation_date,
		encoding,
		name,
		length,
		piece_length,
		pieces,
		private,
		md5_sum
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		t.Info.Hash[:],
		t.Announce,
		t.Comment,
		t.CreatedBy,
		t.CreationDate.Format(time.RFC3339),
		t.Encoding,
		t.Info.Name,
		t.Info.Length,
		t.Info.PieceLength,
		t.Info.Pieces,
		boolToInt(t.Info.Private),
		t.Info.MD5Sum,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	DELETE FROM torrent_files
	WHERE torrent_hash = ?
	`, t.Info.Hash[:])
	if err != nil {
		return err
	}

	for i, file := range t.Info.Files {
		path := strings.Join(file.Path, "/")

		_, err = tx.Exec(`
		INSERT INTO torrent_files (
			torrent_hash,
			file_index,
			length,
			path,
			md5_sum
		)
		VALUES (?, ?, ?, ?, ?)
		`,
			t.Info.Hash[:],
			i,
			file.Length,
			path,
			file.MD5Sum,
		)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *sqlliteTorrentStorage) Find(hashInfo []byte) (*torrent.Torrent, error) {
	row := s.db.QueryRow(`
	SELECT
		announce,
		comment,
		created_by,
		creation_date,
		encoding,
		name,
		length,
		piece_length,
		pieces,
		private,
		md5_sum
	FROM torrents
	WHERE hash = ?
	`, hashInfo)

	var (
		torrentObj torrent.Torrent
		info       torrent.TorrentInfo

		creationDate string
		private      int
		hash         [20]byte
	)

	copy(hash[:], hashInfo)

	err := row.Scan(
		&torrentObj.Announce,
		&torrentObj.Comment,
		&torrentObj.CreatedBy,
		&creationDate,
		&torrentObj.Encoding,
		&info.Name,
		&info.Length,
		&info.PieceLength,
		&info.Pieces,
		&private,
		&info.MD5Sum,
	)
	if err != nil {
		return nil, err
	}

	info.Private = private == 1
	info.Hash = hash

	if creationDate != "" {
		torrentObj.CreationDate, _ = time.Parse(time.RFC3339, creationDate)
	}

	rows, err := s.db.Query(`
	SELECT
		length,
		path,
		md5_sum
	FROM torrent_files
	WHERE torrent_hash = ?
	ORDER BY file_index
	`, hashInfo)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			file torrent.TorrentFile
			path string
		)

		err = rows.Scan(
			&file.Length,
			&path,
			&file.MD5Sum,
		)
		if err != nil {
			return nil, err
		}

		file.Path = strings.Split(path, "/")

		info.Files = append(info.Files, file)
	}

	torrentObj.Info = &info

	return &torrentObj, nil
}

func (s *sqlliteTorrentStorage) Remove(hashInfo []byte) error {
	_, err := s.db.Exec(`
	DELETE FROM torrents
	WHERE hash = ?
	`, hashInfo)

	return err
}

func (s *sqlliteTorrentStorage) SavePeers(torrentHash []byte, peers []peers.Peer) error {
	if len(torrentHash) != 20 {
		return fmt.Errorf("torrent hash must be 20 bytes")
	}

	if len(peers) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`
		INSERT INTO torrent_peers (
			torrent_hash,
			peer_id,
			ip,
			port,
			status,
			failed_attempts,
			last_seen
		)
		VALUES (?, ?, ?, ?, 'discovered', 0, CURRENT_TIMESTAMP)
		ON CONFLICT(torrent_hash, ip, port) DO UPDATE SET
			peer_id = excluded.peer_id,
			last_seen = CURRENT_TIMESTAMP;
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, peer := range peers {
		if peer.Ip == "" {
			continue
		}

		if peer.Port <= 0 || peer.Port > 65535 {
			continue
		}

		_, err = stmt.Exec(
			torrentHash,
			nullableString(peer.Id),
			peer.Ip,
			peer.Port,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *sqlliteTorrentStorage) GetPeers(torrentHash []byte) ([]peers.Peer, error) {
	if len(torrentHash) != 20 {
		return nil, fmt.Errorf("torrent hash must be 20 bytes")
	}

	rows, err := s.db.Query(`
		SELECT
			peer_id,
			ip,
			port
		FROM torrent_peers
		WHERE torrent_hash = ?
		ORDER BY last_seen DESC;
	`, torrentHash)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pList := make([]peers.Peer, 0)

	for rows.Next() {
		var (
			peer   peers.Peer
			peerID sql.NullString
		)

		err = rows.Scan(
			&peerID,
			&peer.Ip,
			&peer.Port,
		)
		if err != nil {
			return nil, err
		}

		if peerID.Valid {
			peer.Id = peerID.String
		}

		pList = append(pList, peer)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return pList, nil
}

func (s *sqlliteTorrentStorage) UpdatePeerStatus(
	torrentHash []byte,
	ip string,
	port int,
	status string,
) error {
	if len(torrentHash) != 20 {
		return fmt.Errorf("torrent hash must be 20 bytes")
	}

	if ip == "" {
		return fmt.Errorf("peer ip is empty")
	}

	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid peer port: %d", port)
	}

	result, err := s.db.Exec(`
		UPDATE torrent_peers
		SET
			status = ?,
			last_seen = CURRENT_TIMESTAMP,
			last_connected_at = CASE
				WHEN ? IN ('connected', 'handshake_sent', 'handshake_received', 'ready')
				THEN CURRENT_TIMESTAMP
				ELSE last_connected_at
			END,
			failed_attempts = CASE
				WHEN ? = 'failed'
				THEN failed_attempts + 1
				ELSE failed_attempts
			END
		WHERE torrent_hash = ?
		  AND ip = ?
		  AND port = ?;
	`,
		status,
		status,
		status,
		torrentHash,
		ip,
		port,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *sqlliteTorrentStorage) InitPieces(infoHash []byte, piecesCount int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO torrent_pieces (
			torrent_hash,
			piece_index,
			completed
		)
		VALUES (?, ?, 0)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range piecesCount {
		if _, err := stmt.Exec(infoHash, i); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *sqlliteTorrentStorage) GetPieces(torrentHash [20]byte) ([]byte, error) {
	rows, err := s.db.Query(`
			SELECT completed
			FROM torrent_pieces
			WHERE torrent_hash = ?
			ORDER BY piece_index
		`, torrentHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pieces []byte

	for rows.Next() {
		var completed byte

		if err := rows.Scan(&completed); err != nil {
			return nil, err
		}

		pieces = append(pieces, completed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pieces, nil
}

// helpers

func boolToInt(v bool) int {
	if v {
		return 1
	}

	return 0
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}
