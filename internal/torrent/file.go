package torrent

import (
	"fmt"
	"my-torrent/internal/bencode"
	"time"
)

type File struct {
	Announce     string // Url of the tracker
	Info         torrentInfo
	Comment      string
	CreatedBy    string
	CreationDate time.Time
}

type torrentInfo struct {
	Files       []torrentFile
	Length      int64
	Name        string // suggested filename where the file is to be saved (if one file)/suggested directory name where the files are to be saved (if multiple files)
	PieceLength int64  // number of bytes per piece. This is commonly 28 KiB = 256 KiB = 262,144 B.
	Pieces      []byte
}

type torrentFile struct {
	Length int64    // size of the file in bytes (only when one file is being shared though)
	Path   []string // A list of strings corresponding to subdirectory names, the last of which is the actual file name
}

const (
	announce     = "announce"
	info         = "info"
	comment      = "comment"
	createdBy    = "created by"
	creationDate = "creation date"
)

func CreateFile(root bencode.BDict) *File {
	file := File{}

	for k, v := range root {
		switch k {
		case announce:
			if val, err := bencode.GetTypedValue[string](v); err == nil {
				file.Announce = val
			}
		case comment:
			if val, err := bencode.GetTypedValue[string](v); err == nil {
				file.Comment = val
			}
		case createdBy:
			if val, err := bencode.GetTypedValue[string](v); err == nil {
				file.CreatedBy = val
			}
		case creationDate:
			if val, err := bencode.GetTypedValue[int64](v); err == nil {
				file.CreationDate = time.UnixMicro(val)
			} else if val, err := bencode.GetTypedValue[int](v); err == nil {
				file.CreationDate = time.UnixMicro(int64(val))
			}
		default:
			fmt.Printf("Warning: Could not determine field with name = %s\n", k)
		}
	}

	return &file
}
