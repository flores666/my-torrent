package torrent

import (
	"time"
)

type Torrent struct {
	Announce     string       // URL of the main tracker (required field)
	AnnounceList [][]string   // list of backup trackers (tracker tiers for redundancy)
	Info         *TorrentInfo // core metadata of the torrent (files, pieces, name, etc.)
	HttpSeeds    []string
	Comment      string    // optional comment from torrent creator
	CreatedBy    string    // name of the program that created the torrent
	CreationDate time.Time // timestamp when the torrent was created
	Encoding     string    // string encoding used in the torrent (usually UTF-8)
	UrlList      []string  // web seed URLs (HTTP/FTP sources for downloading)
}

type TorrentInfo struct {
	Files       []TorrentFile // list of files (used in multi-file torrents)
	Length      int64         // total size of the file (single-file torrent only)
	Name        string        // suggested name for file or root directory
	PieceLength int64         // size of each piece in bytes (commonly 256 KiB)
	Pieces      []byte        // concatenated SHA-1 hashes of all pieces (20 bytes each)
	Private     bool          // private torrent flag (inherited or duplicated from root info)
	MD5Sum      []byte        // MD5 checksum (legacy, rarely used in modern torrents)
	Hash        string
}

type TorrentFile struct {
	Length int64    // size of the file in bytes (multi-file torrents only)
	Path   []string // hierarchical file path (directories + filename)
	MD5Sum []byte   // MD5 checksum of the file (legacy field, rarely used)
}
