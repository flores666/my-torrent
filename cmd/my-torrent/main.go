package main

import (
	"fmt"
	"my-torrent/internal/bencode"
)

func main() { //4:spami32el2:hii10eed3:bar4:spam3:fooi42ee
	res, err := bencode.Decode([]byte("d4:infod4:name12:example_file6:lengthi123456789e12:piece lengthi262144e6:pieces20:abcdefghijklmnopqrst5:filesld6:lengthi1024e4:pathl3:src8:main.goeed6:lengthi2048e4:pathl3:src7:util.goeed6:lengthi4096e4:pathl4:test11:input.txteeee8:announce35:http://tracker.example.com/announce13:announce-listll36:http://tracker1.example.com/announceel36:http://tracker2.example.com/announceee7:comment22:This is a test torrent10:created by11:chatgpt-bot13:creation datei1713600000e5:nodesll9:127.0.0.1i6881eel8:10.0.0.5i51413eee5:extrad5:flagsl4:DHT3:PEX8:ut_metadatae7:privatei1e8:metadatad6:author8:john_doe7:version3:1.04:tagsl3:dev7:example4:testeeee"))
	if err != nil {
		fmt.Printf("Error - %s\n", err)
	}

	bencode.PrintParsed(res)
}
