package peers

import (
	"math/rand/v2"
)

func PickPeer(pr []Peer, port int) *Peer {
	if len(pr) == 0 {
		return nil
	}

	if len(pr) == 1 {
		return &pr[0]
	}

	//todo implement good algorithm
	return &pr[rand.IntN(len(pr)-1)]
}
