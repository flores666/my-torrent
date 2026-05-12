package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func MustLoadNewSqlliteDB(driverName, dbPath string) *sql.DB {
	migrate := false

	_, err := os.Stat(dbPath)
	if err != nil {
		migrate = true
	}

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

CREATE INDEX IF NOT EXISTS idx_torrent_peers_torrent_hash
ON torrent_peers(torrent_hash);

CREATE INDEX IF NOT EXISTS idx_torrent_peers_status
ON torrent_peers(status);

CREATE INDEX IF NOT EXISTS idx_torrent_peers_peer_id
ON torrent_peers(peer_id);`

	_, err := db.Exec(migration)
	if err != nil {
		return err
	}

	return err
}

// #endregion
