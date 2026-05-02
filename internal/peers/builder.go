package peers

import (
	"fmt"
	"my-torrent/internal/bencode"
)

func Build(root bencode.BDict) (*Response, error) {
	response := Response{}
	for k, v := range root {
		switch k {
		case interval:
			if val, err := bencode.GetTypedValue[int64](v); err == nil {
				response.Interval = int(val)
			}
		case peers:
			if val, err := bencode.GetTypedValue[[]any](v); err == nil {
				response.Peers, err = buildPeers(val)
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

func buildPeers(pMaps []any) ([]Peer, error) {
	peers := make([]Peer, 0)

	for _, pm := range pMaps {
		peer := Peer{}
		if pMap, ok := pm.(map[string]any); ok {
			for k, v := range pMap {
				switch k {
				case id:
					if val, ok := v.([]byte); ok {
						peer.Id = string(val)
					}
				case ip:
					if val, ok := v.([]byte); ok {
						peer.Ip = string(val)
					}
				case port:
					if val, ok := v.(int64); ok {
						peer.Port = int(val)
					}
				default:
					return nil, fmt.Errorf("could not determine field with name = %s\n", k)
				}
			}
		}

		peers = append(peers, peer)
	}

	return peers, nil
}
