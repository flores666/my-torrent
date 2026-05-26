package db

import (
	"crypto/rand"
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func MustLoadNewSqlliteDB(driverName, dbPath string) *sql.DB {
	migrate := false

	_, err := os.Stat(dbPath)
	migrate = err != nil

	if driverName == "" {
		panic("driver name must be not empty")
	}

	if dbPath == "" {
		panic("dbPath must be not empty")
	}

	db, err := sql.Open(driverName, dbPath)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		panic(err)
	}

	if migrate {
		if err := applyMigration(db); err != nil {
			panic(err)
		}
	}

	return db
}

// #region migration
func applyMigration(db *sql.DB) error {
	migration := `PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS torrents (
    hash            BLOB PRIMARY KEY CHECK (length(hash) = 20),

    announce        TEXT NOT NULL,
    comment         TEXT,
    created_by      TEXT,
    creation_date   TEXT,
    encoding        TEXT,

    name            TEXT NOT NULL,
    length          INTEGER,
    piece_length    INTEGER NOT NULL,

    pieces          BLOB NOT NULL,
    private         INTEGER NOT NULL DEFAULT 0,
    md5_sum         BLOB,

    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS torrent_pieces (
    torrent_hash BLOB NOT NULL,
    piece_index  INTEGER NOT NULL,
    completed    INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (torrent_hash, piece_index),

    FOREIGN KEY (torrent_hash)
        REFERENCES torrents(info_hash)
        ON DELETE CASCADE,

    CHECK (completed IN (0, 1))
);

CREATE TABLE IF NOT EXISTS torrent_files (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    torrent_hash    BLOB NOT NULL,

    file_index      INTEGER NOT NULL,
    length          INTEGER NOT NULL,
    path            TEXT NOT NULL,
    md5_sum         BLOB,

    UNIQUE (torrent_hash, file_index),

    FOREIGN KEY (torrent_hash)
        REFERENCES torrents(hash)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS torrent_announce_list (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    torrent_hash    BLOB NOT NULL,

    tier            INTEGER NOT NULL,
    url_index       INTEGER NOT NULL,
    url             TEXT NOT NULL,

    UNIQUE (torrent_hash, tier, url_index),

    FOREIGN KEY (torrent_hash)
        REFERENCES torrents(hash)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS torrent_url_list (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    torrent_hash    BLOB NOT NULL,
    url             TEXT NOT NULL,

    UNIQUE (torrent_hash, url),

    FOREIGN KEY (torrent_hash)
        REFERENCES torrents(hash)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS torrent_http_seeds (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    torrent_hash    BLOB NOT NULL,
    url             TEXT NOT NULL,

    UNIQUE (torrent_hash, url),

    FOREIGN KEY (torrent_hash)
        REFERENCES torrents(hash)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS torrent_peers (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,

    torrent_hash        BLOB NOT NULL,

    peer_id             TEXT,
    bitfield 			BLOB,
    
    ip                  TEXT NOT NULL,
    port                INTEGER NOT NULL CHECK (port > 0 AND port <= 65535),

    status              TEXT NOT NULL DEFAULT 'discovered',
    failed_attempts     INTEGER NOT NULL DEFAULT 0,

    last_seen           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_connected_at   TEXT,

    UNIQUE (torrent_hash, ip, port),

    FOREIGN KEY (torrent_hash)
        REFERENCES torrents(hash)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS client_identity (
    id              INTEGER PRIMARY KEY CHECK (id = 1),

    peer_id         TEXT NOT NULL CHECK (length(peer_id) = 20),
    port            INTEGER NOT NULL CHECK (port > 0 AND port <= 65535),

    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO client_identity (
    id,
    peer_id,
    port
)
VALUES (
    1,
    ?,
    1
);

CREATE INDEX IF NOT EXISTS idx_torrent_peers_torrent_hash
ON torrent_peers(torrent_hash);

CREATE INDEX IF NOT EXISTS idx_torrent_peers_status
ON torrent_peers(status);

CREATE INDEX IF NOT EXISTS idx_torrent_peers_peer_id
ON torrent_peers(peer_id);`

	peerId, _ := generatePeerID()

	_, err := db.Exec(migration, peerId)
	if err != nil {
		return err
	}

	return err
}

// #endregion

func generatePeerID() (string, error) {
	const prefix = "-EG1305-"
	const totalLen = 20
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomLen := totalLen - len(prefix)

	buf := make([]byte, randomLen)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	result := make([]byte, randomLen)
	for i, b := range buf {
		result[i] = alphabet[int(b)%len(alphabet)]
	}

	return prefix + string(result), nil
}
