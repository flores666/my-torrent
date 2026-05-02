package main

import (
	"encoding/json"
	"log"
	"my-torrent/internal/peers"
	"my-torrent/internal/torrent"
	"os"
)

func main() {
	torrent, err := torrent.Open("debian.torrent")
	if err != nil {
		log.Fatal(err.Error())
	}

	bytes, _ := json.Marshal(torrent)
	os.WriteFile("debian.json", bytes, 0644)

	r := "*peers*"

	response, err := peers.Parse(r)
	if err != nil {
		log.Fatal(err.Error())
	}
	bytes2, _ := json.Marshal(response)
	os.WriteFile("response.json", bytes2, 0644)
}
