package main

import (
	"fmt"
	"my-torrent/internal/bencode"
)

func main() { //					   0123456789
	res, err := bencode.Decode([]byte("d4:spami32e5:hello5:worlde"))
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res)
}
