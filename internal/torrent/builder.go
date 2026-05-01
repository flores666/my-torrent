package torrent

import (
	"errors"
	"fmt"
	"my-torrent/internal/bencode"
	"my-torrent/internal/utils/url"
	"time"
)

func Build(root bencode.BDict) (*Torrent, error) {
	file := Torrent{}

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
				file.Info, err = buildInfo(val)
				if err != nil {
					return nil, err
				}

				file.Info.HashEncoded = url.Encode(v.Original)
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
				return nil, err
			}

			for _, item := range val {
				b := item.([]byte)

				file.UrlList = append(file.UrlList, string(b))
			}
		default:
			return nil, fmt.Errorf("could not determine field with name = %s\n", k)
		}
	}

	return &file, nil
}

func buildInfo(m map[string]any) (*TorrentInfo, error) {
	info := TorrentInfo{}

	for k, v := range m {
		switch k {
		case length:
			if val, ok := v.(int64); ok {
				info.Length = val
			}
		case name:
			if val, ok := v.([]byte); ok {
				info.Name = string(val)
			}
		case pieceLength:
			if val, ok := v.(int64); ok {
				info.PieceLength = val
			}
		case pieces:
			if val, ok := v.([]byte); ok {
				if len(val)%20 != 0 {
					return nil, errors.New("invalid pieces len")
				}

				info.Pieces = val
			}
		case files:
			if val, ok := v.([]map[string]any); ok {
				var err error
				info.Files, err = buildFiles(val)
				if err != nil {
					return nil, err
				}
			}
		case private:
			if val, ok := v.(int64); ok {
				info.Private = val == 1
			}
		default:
			return nil, fmt.Errorf("could not determine field with name = %s\n", k)
		}
	}

	return &info, nil
}

func buildFiles(fm []map[string]any) ([]TorrentFile, error) {
	files := make([]TorrentFile, 0, len(fm))

	for _, fileMap := range fm {
		file := TorrentFile{}

		for k, v := range fileMap {
			switch k {
			case length:
				if val, ok := v.(int64); ok {
					file.Length = val
				}
			case path:
				if val, ok := v.([][]byte); ok {
					file.Path = findFilePaths(val)
				}
			case md5sum:
				if val, ok := v.([]byte); ok {
					file.MD5Sum = val
				}
			default:
				return nil, fmt.Errorf("could not determine field with name = %s\n", k)
			}
		}

		files = append(files, file)
	}

	return files, nil
}

func findFilePaths(bytes [][]byte) []string {
	paths := make([]string, 0, len(bytes))

	for _, v := range bytes {
		paths = append(paths, string(v))
	}

	return paths
}
