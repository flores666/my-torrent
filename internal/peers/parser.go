package peers

import (
	"errors"
	"fmt"
	"my-torrent/internal/bencode"
)

func Parse(response string) (*Response, error) {
	fmt.Println("Decoding response")

	res, err := bencode.Decode([]byte(response))
	if err != nil {
		return nil, err
	}

	if len(res) <= 0 {
		return nil, errors.New("parsed bencode contains 0 elements")
	}

	fmt.Println("Building response struct")

	if dict, ok := res[0].Value.(bencode.BDict); ok {
		return Build(dict)
	}

	return nil, errors.New("could not read response")
}
