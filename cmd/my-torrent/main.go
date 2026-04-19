package main

import (
	"my-torrent/internal/bencode"
)

func main() { //					   0123456789
	//res, err := bencode.Decode([]byte("d4:spami32e5:hello5:worldl4:spami42eee"))
	res, _ := bencode.Decode([]byte("4:spami32el2:hii10ee"))

	bencode.PrintParsed(res)
}
