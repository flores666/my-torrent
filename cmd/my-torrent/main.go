package main

import (
	"fmt"
	"my-torrent/internal/bencode"
)

func main() {
	res, err := bencode.Decode([]byte("d8:announce26:http://tracker.example.com4:infod6:lengthi12345e4:name11:example.txt12:piece lengthi16384e6:pieces20:abcdefghijklmnopqrste4:listld3:foo3:bared3:numi42eeee"))
	//res, err := bencode.Decode([]byte("ld3:foo3:bared3:numi42eee"))
	if err != nil {
		fmt.Printf("Error - %s\n", err)
	}

	bencode.PrintParsed(res)
}
