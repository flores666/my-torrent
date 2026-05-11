package torrent

import (
	"errors"
	"fmt"
	"my-torrent/internal/bencode"
	"os"
)

func Open(filePath string) (*Torrent, error) {
	fmt.Println("Opening torrent file")

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fmt.Println("Reading torrent file")

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())
	count, err := file.Read(data)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Decoding torrent file, got %d bytes\n", count)

	res, err := bencode.Decode(data[:count])
	if err != nil {
		return nil, err
	}

	if len(res) <= 0 {
		return nil, errors.New("parsed bencode contains 0 elements")
	}

	fmt.Println("Building torrent file")

	if dict, ok := res[0].Value.(bencode.BDict); ok {
		return Build(dict)
	}

	return nil, errors.New("could not read .torrent file")
}
