package main

import (
	"encoding/json"
	"fmt"
	"log"
	"my-torrent/internal/requests"
	"my-torrent/internal/torrent"
	"os"
)

func main() {
	torrent, err := torrent.Open("debian.torrent")
	if err != nil {
		log.Fatalln(err.Error())
	}

	peers, listener := requests.GetPeers(torrent)
	defer listener.Close()

	fmt.Println("Writing to file")

	bytes, _ := json.Marshal(peers)
	os.WriteFile("response.json", bytes, 0644)
}
