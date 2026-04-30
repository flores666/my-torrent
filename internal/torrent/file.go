package torrent

import (
	"fmt"
	"my-torrent/internal/bencode"
	"time"
)

type File struct {
	Announce     string // Url of the tracker
	Info         *torrentInfo
	Comment      string
	CreatedBy    string
	CreationDate time.Time
	UrlList      []string
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

func CreateFile(root bencode.BDict) *File {
	file := File{}

	for k, v := range root {
		switch k {
		case announce:
			if val, err := bencode.GetTypedValue[[]byte](v); err == nil {
				file.Announce = string(val)
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
			} else if val, err := bencode.GetTypedValue[int64](v); err == nil {
				file.CreationDate = time.UnixMicro(int64(val))
			}
		case info:
			if val, err := bencode.GetTypedValue[map[string]any](v); err == nil {
				file.Info = buildInfo(val)
			}
		case urlList:
			val, err := bencode.GetTypedValue[[]interface{}](v)
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			urlList := make([]string, 0, len(val))

			for _, item := range val {
				b := item.([]byte)

				urlList = append(urlList, string(b))
			}

			file.UrlList = urlList
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
