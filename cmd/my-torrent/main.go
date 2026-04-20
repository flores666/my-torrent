package main

import (
	"fmt"
	"my-torrent/internal/bencode"
)

func main() { //4:spami32el2:hii10eed3:bar4:spam3:fooi42ee
	res, err := bencode.Decode([]byte("d3:bar4:spam3:fooi42ee"))
	if err != nil {
		fmt.Printf("Error - %s\n", err)
	}

	bencode.PrintParsed(res)
}
