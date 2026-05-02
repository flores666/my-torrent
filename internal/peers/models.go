package peers

type Response struct {
	Interval int
	Peers    []Peer
}

type Peer struct {
	Id   string
	Ip   string
	Port int
}
