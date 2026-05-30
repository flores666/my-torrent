package peers

import (
	"fmt"
)

func BuildPeers(pMaps []any) ([]*Peer, error) {
	peers := make([]*Peer, 0)

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

		peers = append(peers, &peer)
	}

	return peers, nil
}
