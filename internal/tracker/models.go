package tracker

import (
	"errors"
	"fmt"
	"my-torrent/internal/bencode"
	"my-torrent/internal/modelbuilder/peers"
)

const (
	intervalKey = "interval"
	peersKey    = "peers"
)

type arguments struct {
	announce     string
	announcelist [][]string
	peerId       string
	port         int
	uploaded     int64
	downloaded   int64
	left         int64
	event        string
	infoHash     [20]byte
}

type AnnounceResponse struct {
	Interval int
	Peers    []peers.Peer
}

func ParseAnnounceResponse(response []byte) (*AnnounceResponse, error) {
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
		return buildAnnounceResponse(dict)
	}

	return nil, errors.New("could not read response")
}

func buildAnnounceResponse(root bencode.BDict) (*AnnounceResponse, error) {
	response := AnnounceResponse{}
	for k, v := range root {
		switch k {
		case intervalKey:
			if val, err := bencode.GetTypedValue[int64](v); err == nil {
				response.Interval = int(val)
			}
		case peersKey:
			if val, err := bencode.GetTypedValue[[]any](v); err == nil {
				response.Peers, err = peers.BuildPeers(val)
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, fmt.Errorf("could not determine field with name = %s\n", k)
		}
	}

	return &response, nil
}
