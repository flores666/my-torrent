package main

import (
	"encoding/json"
	"fmt"
	"log"
	"my-torrent/internal/bencode"
	"my-torrent/internal/torrent"
	"os"
)

func main() {
	fmt.Println("Opening torrent file")

	file, err := os.Open("debian.torrent")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Println("Reading torrent file")

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, stat.Size())
	count, err := file.Read(data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoding torrent file, got %d bytes\n", count)

	res, err := bencode.Decode(data[:count])
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("debian.json", bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if dict, ok := res[0].(bencode.BDict); ok {
		tf := torrent.CreateFile(dict)
		byte, _ := json.Marshal(tf)

		fmt.Printf("TorrentFile: %s\n", string(byte))
	}
}
