package torrent

import (
	"fmt"
	"my-torrent/internal/bencode"
	"time"
)

type File struct {
	Announce     string       // URL of the main tracker (required field)
	AnnounceList [][]string   // list of backup trackers (tracker tiers for redundancy)
	Info         *torrentInfo // core metadata of the torrent (files, pieces, name, etc.)
	HttpSeeds    []string
	Comment      string    // optional comment from torrent creator
	CreatedBy    string    // name of the program that created the torrent
	CreationDate time.Time // timestamp when the torrent was created
	Encoding     string    // string encoding used in the torrent (usually UTF-8)
	UrlList      []string  // web seed URLs (HTTP/FTP sources for downloading)
}

type torrentInfo struct {
	Files       []torrentFile // list of files (used in multi-file torrents)
	Length      int64         // total size of the file (single-file torrent only)
	Name        string        // suggested name for file or root directory
	PieceLength int64         // size of each piece in bytes (commonly 256 KiB)
	Pieces      []byte        // concatenated SHA-1 hashes of all pieces (20 bytes each)
	PieceHashes [][]byte      // slice of individual 20-byte SHA-1 piece hashes
	Private     bool          // private torrent flag (inherited or duplicated from root info)
	MD5Sum      []byte        // MD5 checksum (legacy, rarely used in modern torrents)
}

type torrentFile struct {
	Length int64    // size of the file in bytes (multi-file torrents only)
	Path   []string // hierarchical file path (directories + filename)
	MD5Sum []byte   // MD5 checksum of the file (legacy field, rarely used)
}

func CreateFile(root bencode.BDict) *File {
	file := File{}

	for k, v := range root {
		switch k {
		case announce:
			if val, err := bencode.GetTypedValue[[]byte](v); err == nil {
				file.Announce = string(val)
			}
		case announceList:
			if val, err := bencode.GetTypedValue[[][][]byte](v); err == nil {
				result := make([][]string, len(val))

				for i, tier := range val {
					for _, tracker := range tier {
						result[i] = append(result[i], string(tracker))
					}
				}

				file.AnnounceList = result
			}
		case comment:
			if val, err := bencode.GetTypedValue[[]byte](v); err == nil {
				file.Comment = string(val)
			}
		case createdBy:
			if val, err := bencode.GetTypedValue[[]byte](v); err == nil {
				file.CreatedBy = string(val)
			}
		case creationDate:
			if val, err := bencode.GetTypedValue[int64](v); err == nil {
				file.CreationDate = time.UnixMicro(val)
			}
		case encoding:
			if val, err := bencode.GetTypedValue[[]byte](v); err == nil {
				file.Encoding = string(val)
			}
		case info:
			if val, err := bencode.GetTypedValue[map[string]any](v); err == nil {
				file.Info = buildInfo(val)
			}
		case httpseeds:
			if val, err := bencode.GetTypedValue[[]interface{}](v); err == nil {
				for _, x := range val {
					file.HttpSeeds = append(file.HttpSeeds, string(x.([]byte)))
				}
			}
		case urlList:
			val, err := bencode.GetTypedValue[[]interface{}](v)
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			for _, item := range val {
				b := item.([]byte)

				file.UrlList = append(file.UrlList, string(b))
			}
		default:
			fmt.Printf("Warning: Could not determine field with name = %s\n", k)
		}
	}

	return &file
}

func buildInfo(m map[string]any) *torrentInfo {
	info := torrentInfo{}

	for k, v := range m {
		switch k {
		case length:
			if val, err := bencode.GetTypedValue[int64](v); err == nil {
				info.Length = val
			}
		case name:
			if val, err := bencode.GetTypedValue[[]byte](v); err == nil {
				info.Name = string(val)
			}
		case pieceLength:
			if val, err := bencode.GetTypedValue[int64](v); err == nil {
				info.PieceLength = val
			}
		case pieces:
			if val, err := bencode.GetTypedValue[[]byte](v); err == nil {
				if len(val)%20 != 0 {
					fmt.Println("invalid pieces length")
					return nil
				}

				info.Pieces = val
			}
		case files:
			if val, err := bencode.GetTypedValue[[]map[string]any](v); err == nil {
				info.Files = buildFiles(val)
			}
		case private:
			if val, err := bencode.GetTypedValue[int64](v); err == nil {
				info.Private = val == 1
			}
		default:
			fmt.Printf("Warning: Could not determine field with name = %s\n", k)
		}
	}

	return &info
}

func buildFiles(fm []map[string]any) []torrentFile {
	files := make([]torrentFile, 0, len(fm))

	for _, fileMap := range fm {
		file := torrentFile{}

		for k, v := range fileMap {
			switch k {
			case length:
				if val, err := bencode.GetTypedValue[int64](v); err == nil {
					file.Length = val
				}
			case path:
				if val, err := bencode.GetTypedValue[[][]byte](v); err == nil {
					file.Path = findFilePaths(val)
				}
			case md5sum:
				if val, err := bencode.GetTypedValue[[]byte](v); err == nil {
					file.MD5Sum = val
				}
			default:
				fmt.Printf("Warning: Could not determine field with name = %s\n", k)
			}
		}

		files = append(files, file)
	}

	return files
}

func findFilePaths(bytes [][]byte) []string {
	paths := make([]string, 0, len(bytes))

	for _, v := range bytes {
		paths = append(paths, string(v))
	}

	return paths
}
