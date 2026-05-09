package peers

import "fmt"

type Response struct {
	Interval int
	Peers    []Peer
}

type Peer struct {
	Id   string
	Ip   string
	Port int
}

func (p *Peer) Address() string {
	return fmt.Sprintf("%s:%d", p.Ip, p.Port)
}
