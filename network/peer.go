package network

import (
	"strconv"
)

// Peer ...
type Peer struct {
	address string
	port    int
	name    string
	Send    chan Message
}

type PeersMap interface {
	Get(address string) Peer
	Find(address string) (bool, Peer)
}

func (p Peer) String() string {
	if p.name != "" {
		return p.name
	} else {
		return p.address + ":" + strconv.Itoa(p.port)
	}
}

func (p Peer) FullAddress() string {
	return p.address + ":" + strconv.Itoa(p.port)
}

func (p Peer) Name() string {
	return p.name
}

func (p *Peer) SetName(name string) {
	p.name = name
}

func CreatePeer(addr string, port int) Peer {
	return Peer{
		address: addr,
		port:    port,
		Send:    make(chan Message),
	}
}
