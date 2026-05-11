package peers

import (
	"errors"
	"fmt"
	"my-torrent/internal/bencode"
)

func Parse(response []byte) (*Response, error) {
	fmt.Println("Decoding peers response")

	res, err := bencode.Decode(response)
	if err != nil {
		return nil, err
	}

	if len(res) <= 0 {
		return nil, errors.New("parsed bencode contains 0 elements")
	}

	fmt.Println("Building peers response struct")

	if dict, ok := res[0].Value.(bencode.BDict); ok {
		return Build(dict)
	}

	return nil, errors.New("could not read response")
}
