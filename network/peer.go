package network

import "strconv"

// Peer ...
type Peer struct {
	address string
	port    int
	name    string
	Send    chan Message
}

func (p Peer) String() string {
	if p.name != "" {
		return p.name
	} else {
		return p.address + ":" + strconv.Itoa(p.port)
	}
}
